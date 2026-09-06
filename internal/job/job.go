// Package job runs the per-file pipeline shared by every compress and convert
// command: existence check, dry-run preview, spinner, external tool steps,
// size report, in-place replacement and a final summary with exit status.
package job

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"charm.land/log/v2"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/message"
	"github.com/y3owk1n/uts/internal/util"
)

const (
	tailLines      = 12
	percent        = 100
	secondsPerMin  = 60
	secondsPerHour = 3600
	// minProgressForETA avoids wild estimates in the first moments.
	minProgressForETA = 0.02
)

var errNoOutput = errors.New("no output produced")

// Step is one unit of work for a file: either an external command or a Go
// function with a human-readable description for dry runs.
type Step struct {
	Cmd  *exec.Cmd
	Desc string
	Fn   func() error
	// Progress, when set, makes the runner read the command's stdout line
	// by line and report completion as it goes.
	Progress *Progress
}

// Progress describes how to read completion out of a command's stdout.
type Progress struct {
	// Total is the amount of work in the same unit Parse returns, e.g.
	// seconds of media. Zero disables the display but still drains stdout.
	Total float64
	// Parse extracts the work done so far from one stdout line.
	Parse func(line string) (float64, bool)
}

// Exec builds a command step.
func Exec(name string, args ...string) Step {
	return Step{Cmd: exec.CommandContext(context.Background(), name, args...)}
}

// Do builds a function step.
func Do(desc string, fn func() error) Step {
	return Step{Desc: desc, Fn: fn}
}

// Job describes how one input becomes one output.
type Job struct {
	Input  string
	Output string
	Steps  []Step
	// Note is shown under the step line, e.g. the tool chosen for this file.
	Note string
	// Skip, when set, is the reason the file is left alone.
	Skip string
}

// Options configures a Run.
type Options struct {
	// Verb is the present participle shown while working, e.g. "Compressing".
	Verb string
	// Done is the past tense used in the summary, e.g. "Compressed".
	Done string
	// Noun names one input in the summary, e.g. "image file"; it is
	// pluralized as needed.
	Noun string
	// Compare reports the size delta and refuses to replace an original with
	// a larger result.
	Compare bool
	InPlace bool
	DryRun  bool
	// Jobs is how many files run at once. Values below 2 run sequentially.
	Jobs int
	// SkipExisting leaves a file alone when its output already exists, so an
	// interrupted batch can be re-run without redoing finished files.
	SkipExisting bool
	// Backup keeps the original as <name>.bak when InPlace replaces it.
	Backup bool
	// Code is the error code returned when any file fails.
	Code derrors.Code
}

type outcome int

const (
	outcomeDone outcome = iota
	outcomeFailed
	outcomeSkipped
)

type result struct {
	outcome outcome
	saved   int64
}

type runner struct {
	opts  Options
	plan  func(string) (*Job, error)
	total int
}

// Run plans and executes one job per file. It never stops at the first
// failure; every file gets its turn and an error summarizing the failures is
// returned at the end so the CLI exits non-zero.
func Run(files []string, opts Options, plan func(file string) (*Job, error)) error {
	runner := &runner{opts: opts, plan: plan, total: len(files)}

	var results []result
	if opts.Jobs > 1 && !opts.DryRun && len(files) > 1 {
		results = runner.parallel(files)
	} else {
		results = runner.sequential(files)
	}

	var done, failed, skipped int

	var saved int64

	for _, res := range results {
		switch res.outcome {
		case outcomeDone:
			done++
			saved += res.saved
		case outcomeFailed:
			failed++
		case outcomeSkipped:
			skipped++
		}
	}

	if runner.total > 1 {
		summarize(opts, done, failed, skipped, saved)
	}

	if failed > 0 {
		return derrors.Newf(opts.Code, "%d of %d files failed", failed, runner.total)
	}

	return nil
}

func (r *runner) sequential(files []string) []result {
	results := make([]result, 0, len(files))

	for idx, file := range files {
		results = append(results, r.process(idx, file, ui.Message, ui.NewSpinner(nil, 0)))
	}

	return results
}

// parallel runs up to opts.Jobs files at once. Each file's messages are
// buffered and flushed when it finishes so lines never interleave, and one
// spinner tracks overall progress.
func (r *runner) parallel(files []string) []result {
	results := make([]result, len(files))
	sem := make(chan struct{}, r.opts.Jobs)

	var (
		workers   sync.WaitGroup
		flush     sync.Mutex
		completed int
	)

	spinner := ui.NewSpinner(nil, 0)
	spinner.SetSuffix(
		fmt.Sprintf("%s %d files, %d at a time...", r.opts.Verb, len(files), r.opts.Jobs),
	)
	spinner.Start()

	for idx, file := range files {
		workers.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			var out, errOut bytes.Buffer

			res := r.process(idx, file, ui.Message.WithWriters(&out, &errOut), nil)

			flush.Lock()
			defer flush.Unlock()

			results[idx] = res
			completed++

			spinner.Stop()

			_, _ = os.Stdout.Write(out.Bytes())
			_, _ = os.Stderr.Write(errOut.Bytes())

			if completed < len(files) {
				spinner.SetSuffix(fmt.Sprintf("%s %d/%d files, %d at a time...",
					r.opts.Verb, completed, len(files), r.opts.Jobs))
				spinner.Start()
			}
		})
	}

	workers.Wait()
	spinner.Stop()

	return results
}

// process handles one file and reports through printer. spinner is nil when the
// caller shows its own progress.
func (r *runner) process(
	idx int,
	file string,
	printer *message.Printer,
	spinner *ui.Spinner,
) result {
	opts := r.opts

	if !util.FileExists(file) {
		printer.Warnf("File not found: %s", file)

		return result{outcome: outcomeFailed}
	}

	job, err := r.plan(file)
	if err != nil {
		printer.Errorf("%s: %v", file, err)

		return result{outcome: outcomeFailed}
	}

	if job.Skip != "" {
		printer.Warnf("%s, skipping: %s", job.Skip, file)

		return result{outcome: outcomeSkipped}
	}

	if opts.SkipExisting && !opts.DryRun && outputExists(job) {
		printer.Warnf("Output exists, skipping: %s", job.Output)

		return result{outcome: outcomeSkipped}
	}

	origSize := util.FileSize(file)
	printer.Stepf("[%d/%d] %s (%s)", idx+1, r.total, file, util.HumanSize(origSize))

	if job.Note != "" {
		printer.Mutedf("  %s", job.Note)
	}

	if opts.DryRun {
		preview(printer, job, opts.InPlace)

		return result{outcome: outcomeDone}
	}

	err = execute(job, opts, spinner)
	if err != nil {
		printer.Errorf("%s failed: %s", strings.TrimSuffix(opts.Verb, "ing")+"ion", file)
		printer.Mutedf("%s", err)

		return result{outcome: outcomeFailed}
	}

	newSize := util.FileSize(job.Output)

	var saved int64

	switch {
	case opts.Compare:
		printer.Successf("%s: %s → %s %s", file, util.HumanSize(origSize),
			util.HumanSize(newSize), util.CompressionRatio(origSize, newSize))

		if newSize >= origSize {
			printer.Warnf("Result is not smaller than the original, keeping %s", file)

			if opts.InPlace {
				_ = os.Remove(job.Output)
			}

			return result{outcome: outcomeSkipped}
		}

		saved = origSize - newSize
	case isDir(job.Output):
		printer.Successf("%s → %s/", file, job.Output)
	default:
		printer.Successf("%s → %s (%s)", file, job.Output, util.HumanSize(newSize))
	}

	if opts.InPlace {
		replaceInPlace(job, opts.Backup)
	}

	return result{outcome: outcomeDone, saved: saved}
}

// outputExists reports whether a job's file output is already on disk.
// Directory outputs (extracted archives, PDF pages) are never "existing"
// because the directory is shared between inputs.
func outputExists(job *Job) bool {
	return job.Output != job.Input && util.FileExists(job.Output) && !isDir(job.Output)
}

func preview(printer *message.Printer, job *Job, inPlace bool) {
	printer.Infof("[dry-run] Would write %s%s", job.Output, util.InPlaceHint(inPlace))

	for _, step := range job.Steps {
		if step.Cmd != nil {
			printer.Mutedf("  $ %s", util.Describe(step.Cmd))
		} else {
			printer.Mutedf("  %s", step.Desc)
		}
	}
}

func execute(job *Job, opts Options, spinner *ui.Spinner) error {
	err := util.EnsureDir(job.Output)
	if err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	existed := util.FileExists(job.Output)

	if spinner != nil {
		spinner.SetSuffix(fmt.Sprintf("%s %s...", opts.Verb, job.Input))
		spinner.Start()

		defer spinner.Stop()
	}

	started := time.Now()
	report := func(done, total float64) {
		if total <= 0 || spinner == nil {
			return
		}

		spinner.SetSuffix(fmt.Sprintf("%s %s... %s", opts.Verb, job.Input,
			progressText(done, total, time.Since(started))))
	}

	for _, step := range job.Steps {
		err := runStep(step, report)
		if err != nil {
			if !existed {
				removeFile(job.Output)
			}

			return err
		}
	}

	if !util.FileExists(job.Output) {
		return fmt.Errorf("%w at %s", errNoOutput, job.Output)
	}

	return nil
}

func runStep(step Step, report func(done, total float64)) error {
	if step.Cmd == nil {
		log.Debug("step", "desc", step.Desc)

		return step.Fn()
	}

	log.Debug("exec", "cmd", util.Describe(step.Cmd))

	if step.Progress == nil {
		out, err := step.Cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %w\n%s", step.Cmd.Args[0], err, tail(out))
		}

		return nil
	}

	return runWithProgress(step, report)
}

// runWithProgress streams the command's stdout through the step's parser so
// the spinner can show how far along a long encode is. stderr is kept for
// the error message.
func runWithProgress(step Step, report func(done, total float64)) error {
	var stderr bytes.Buffer

	step.Cmd.Stderr = &stderr

	stdout, err := step.Cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%s: %w", step.Cmd.Args[0], err)
	}

	err = step.Cmd.Start()
	if err != nil {
		return fmt.Errorf("%s: %w", step.Cmd.Args[0], err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if done, ok := step.Progress.Parse(scanner.Text()); ok {
			report(done, step.Progress.Total)
		}
	}

	err = step.Cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s: %w\n%s", step.Cmd.Args[0], err, tail(stderr.Bytes()))
	}

	return nil
}

// progressText renders "42% · 1:12 left" from work done, total work and the
// time spent so far.
func progressText(done, total float64, elapsed time.Duration) string {
	frac := min(max(done/total, 0), 1)
	text := fmt.Sprintf("%3.0f%%", frac*percent)

	if frac > minProgressForETA && frac < 1 {
		remaining := time.Duration(float64(elapsed) * (1 - frac) / frac).Round(time.Second)
		text += " · " + clock(remaining) + " left"
	}

	return text
}

func clock(dur time.Duration) string {
	secs := int(dur.Seconds())

	if secs >= secondsPerHour {
		return fmt.Sprintf(
			"%d:%02d:%02d",
			secs/secondsPerHour,
			secs%secondsPerHour/secondsPerMin,
			secs%secondsPerMin,
		)
	}

	return fmt.Sprintf("%d:%02d", secs/secondsPerMin, secs%secondsPerMin)
}

// tail keeps the last lines of tool output and drops carriage-return progress
// updates so a failed ffmpeg run prints a few useful lines, not a screenful.
func tail(out []byte) string {
	var lines []string

	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if idx := strings.LastIndex(line, "\r"); idx >= 0 {
			line = line[idx+1:]
		}

		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) > tailLines {
		lines = lines[len(lines)-tailLines:]
	}

	return strings.Join(lines, "\n")
}

func isDir(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func removeFile(path string) {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return
	}

	_ = os.Remove(path)
}

// replaceInPlace makes the output take the original's place. Same-format
// results overwrite the original; format-changing results delete the original
// and take its base name. With backup the original survives as <name>.bak.
func replaceInPlace(job *Job, backup bool) {
	if backup {
		_ = os.Rename(job.Input, job.Input+".bak")
	}

	if format.Same(format.Ext(job.Output), format.Ext(job.Input)) {
		util.MaybeInPlace(job.Output, job.Input)

		return
	}

	util.RemoveInPlace(job.Input)

	target := util.ConvertOutputPath(job.Input, format.Ext(job.Output))
	if target != job.Output {
		_ = os.Rename(job.Output, target)
	}
}

func plural(count int, noun string) string {
	if count == 1 {
		return "1 " + noun
	}

	return fmt.Sprintf("%d %ss", count, noun)
}

func summarize(opts Options, done, failed, skipped int, saved int64) {
	if opts.DryRun {
		ui.Message.Infof("[dry-run] %s previewed, nothing written", plural(done, opts.Noun))

		return
	}

	msg := opts.Done + " " + plural(done, opts.Noun)

	var extra []string
	if skipped > 0 {
		extra = append(extra, fmt.Sprintf("%d skipped", skipped))
	}

	if failed > 0 {
		extra = append(extra, fmt.Sprintf("%d failed", failed))
	}

	if len(extra) > 0 {
		msg += " (" + strings.Join(extra, ", ") + ")"
	}

	if opts.Compare && saved > 0 {
		msg += ", saved " + util.HumanSize(saved)
	}

	if failed > 0 {
		ui.Message.Warnf("%s", msg)
	} else {
		ui.Message.Successf("%s", msg)
	}
}
