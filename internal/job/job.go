// Package job runs the per-file pipeline shared by every compress and convert
// command: existence check, dry-run preview, spinner, external tool steps,
// size report, in-place replacement and a final summary with exit status.
package job

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"charm.land/log/v2"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

const tailLines = 12

var errNoOutput = errors.New("no output produced")

// Step is one unit of work for a file: either an external command or a Go
// function with a human-readable description for dry runs.
type Step struct {
	Cmd  *exec.Cmd
	Desc string
	Fn   func() error
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
	// Noun names the inputs in the summary, e.g. "image files".
	Noun string
	// Compare reports the size delta and refuses to replace an original with
	// a larger result.
	Compare bool
	InPlace bool
	DryRun  bool
	// Code is the error code returned when any file fails.
	Code derrors.Code
}

// Run plans and executes one job per file. It never stops at the first
// failure; every file gets its turn and an error summarizing the failures is
// returned at the end so the CLI exits non-zero.
func Run(files []string, opts Options, plan func(file string) (*Job, error)) error {
	total := len(files)

	var done, failed, skipped int

	var saved int64

	for idx, file := range files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			failed++

			continue
		}

		job, err := plan(file)
		if err != nil {
			ui.Message.Errorf("%s: %v", file, err)

			failed++

			continue
		}

		if job.Skip != "" {
			ui.Message.Warnf("%s, skipping: %s", job.Skip, file)

			skipped++

			continue
		}

		origSize := util.FileSize(file)
		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, total, file, util.HumanSize(origSize))

		if job.Note != "" {
			ui.Message.Mutedf("  %s", job.Note)
		}

		if opts.DryRun {
			preview(job, opts.InPlace)

			done++

			continue
		}

		err = execute(job, opts)
		if err != nil {
			ui.Message.Errorf("%s failed: %s", strings.TrimSuffix(opts.Verb, "ing")+"ion", file)
			ui.Message.Mutedf("%s", err)

			failed++

			continue
		}

		newSize := util.FileSize(job.Output)

		switch {
		case opts.Compare:
			ui.Message.Successf("%s: %s → %s %s", file, util.HumanSize(origSize),
				util.HumanSize(newSize), util.CompressionRatio(origSize, newSize))

			if newSize >= origSize {
				ui.Message.Warnf("Result is not smaller than the original, keeping %s", file)

				if opts.InPlace {
					_ = os.Remove(job.Output)
				}

				skipped++

				continue
			}

			saved += origSize - newSize
		case isDir(job.Output):
			ui.Message.Successf("%s → %s/", file, job.Output)
		default:
			ui.Message.Successf("%s → %s (%s)", file, job.Output, util.HumanSize(newSize))
		}

		if opts.InPlace {
			replaceInPlace(job)
		}

		done++
	}

	if total > 1 {
		summarize(opts, done, failed, skipped, saved)
	}

	if failed > 0 {
		return derrors.Newf(opts.Code, "%d of %d files failed", failed, total)
	}

	return nil
}

func preview(job *Job, inPlace bool) {
	ui.Message.Infof("[dry-run] Would write %s%s", job.Output, util.InPlaceHint(inPlace))

	for _, step := range job.Steps {
		if step.Cmd != nil {
			ui.Message.Mutedf("  $ %s", util.Describe(step.Cmd))
		} else {
			ui.Message.Mutedf("  %s", step.Desc)
		}
	}
}

func execute(job *Job, opts Options) error {
	err := util.EnsureDir(job.Output)
	if err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	existed := util.FileExists(job.Output)

	spinner := ui.NewSpinner(nil, 0)
	spinner.SetSuffix(fmt.Sprintf("%s %s...", opts.Verb, job.Input))
	spinner.Start()

	defer spinner.Stop()

	for _, step := range job.Steps {
		err := runStep(step)
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

func runStep(step Step) error {
	if step.Cmd == nil {
		log.Debug("step", "desc", step.Desc)

		return step.Fn()
	}

	log.Debug("exec", "cmd", util.Describe(step.Cmd))

	out, err := step.Cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w\n%s", step.Cmd.Args[0], err, tail(out))
	}

	return nil
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
// and take its base name.
func replaceInPlace(job *Job) {
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

func summarize(opts Options, done, failed, skipped int, saved int64) {
	if opts.DryRun {
		ui.Message.Infof("[dry-run] %d %s previewed, nothing written", done, opts.Noun)

		return
	}

	msg := fmt.Sprintf("%s %d %s", opts.Done, done, opts.Noun)

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
