// Package archive creates, extracts and lists archives by driving tar, zip and
// the standalone compressors. Compressed tarballs are streamed through a pipe
// so no temporary file is ever written next to the archive.
//
//nolint:goconst,mnd
package archive

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/log/v2"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// kind is what the file name says about an archive.
type kind struct {
	zip bool
	// filter is the standalone decompressor needed before tar can read the
	// stream: "" (plain tar or a format tar detects itself), "zstd" or "brotli".
	filter string
}

// detect classifies an archive by its (possibly compound) extension. Plain
// .gz/.xz/.bz2 files are treated as compressed tarballs, which is what uts
// itself produces.
func detect(file string) (kind, error) {
	name := strings.ToLower(file)

	switch {
	case strings.HasSuffix(name, ".zip"):
		return kind{zip: true}, nil
	case strings.HasSuffix(name, ".tar"),
		strings.HasSuffix(name, ".gz"), strings.HasSuffix(name, ".tgz"),
		strings.HasSuffix(name, ".xz"), strings.HasSuffix(name, ".txz"),
		strings.HasSuffix(name, ".bz2"), strings.HasSuffix(name, ".tbz2"), strings.HasSuffix(name, ".tbz"):
		return kind{}, nil
	case strings.HasSuffix(name, ".zst"),
		strings.HasSuffix(name, ".zstd"),
		strings.HasSuffix(name, ".tzst"):
		return kind{filter: "zstd"}, nil
	case strings.HasSuffix(name, ".br"):
		return kind{filter: "brotli"}, nil
	}

	return kind{}, derrors.Newf(
		derrors.CodeUnsupportedFormat,
		"unsupported archive: %s",
		filepath.Base(file),
	)
}

func needTool(name, pkg string) error {
	if util.HasTool(name) {
		return nil
	}

	return derrors.Newf(
		derrors.CodeToolNotFound,
		"%s not found — install: brew install %s",
		name,
		pkg,
	)
}

// tarStep returns the step that runs "tar <tarArgs>" over the archive, piping
// through the standalone decompressor when the format needs one.
func tarStep(archiveKind kind, file string, tarArgs ...string) (job.Step, error) {
	if archiveKind.filter == "" {
		return job.Exec("tar", append(tarArgs, "-f", file)...), nil
	}

	err := needTool(archiveKind.filter, archiveKind.filter)
	if err != nil {
		return job.Step{}, err
	}

	producer := exec.CommandContext(context.Background(), archiveKind.filter, "-dc", file)
	consumer := exec.CommandContext(context.Background(), "tar", append(tarArgs, "-f", "-")...)

	return job.Do(util.Describe(producer)+" | "+util.Describe(consumer), func() error {
		_, err := pipe(producer, consumer)

		return err
	}), nil
}

// pipe runs producer | consumer and returns the consumer's stdout.
func pipe(producer, consumer *exec.Cmd) ([]byte, error) {
	log.Debug("exec", "cmd", util.Describe(producer)+" | "+util.Describe(consumer))

	var (
		producerErr bytes.Buffer
		consumerOut bytes.Buffer
		consumerErr bytes.Buffer
	)

	stdout, err := producer.StdoutPipe()
	if err != nil {
		return nil, err
	}

	producer.Stderr = &producerErr
	consumer.Stdin = stdout
	consumer.Stdout = &consumerOut
	consumer.Stderr = &consumerErr

	err = producer.Start()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", producer.Args[0], err)
	}

	consumerRunErr := consumer.Run()

	producerWaitErr := producer.Wait()
	if producerWaitErr != nil {
		return nil, fmt.Errorf(
			"%s: %w\n%s",
			producer.Args[0],
			producerWaitErr,
			strings.TrimSpace(producerErr.String()),
		)
	}

	if consumerRunErr != nil {
		return nil, fmt.Errorf(
			"%s: %w\n%s",
			consumer.Args[0],
			consumerRunErr,
			strings.TrimSpace(consumerErr.String()),
		)
	}

	return consumerOut.Bytes(), nil
}

// ExtractOptions represents options for archive extraction.
type ExtractOptions struct {
	Files     []string
	OutputDir string
	DryRun    bool
}

// Extract extracts archive files into the output directory (default: the
// current directory).
func Extract(opts ExtractOptions) error {
	outDir := opts.OutputDir
	if outDir == "" {
		outDir = "."
	}

	return job.Run(opts.Files, job.Options{
		Verb:   "Extracting",
		Done:   "Extracted",
		Noun:   "archive",
		DryRun: opts.DryRun,
		Code:   derrors.CodeArchiveFailed,
	}, func(file string) (*job.Job, error) {
		archiveKind, err := detect(file)
		if err != nil {
			return nil, err
		}

		mkdir := job.Do("mkdir -p "+outDir, func() error { return os.MkdirAll(outDir, 0o755) })

		var step job.Step

		if archiveKind.zip {
			err = needTool("unzip", "unzip")
			if err != nil {
				return nil, err
			}

			step = job.Exec("unzip", "-qo", file, "-d", outDir)
		} else {
			step, err = tarStep(archiveKind, file, "-x", "-C", outDir)
			if err != nil {
				return nil, err
			}
		}

		return &job.Job{Input: file, Output: outDir, Steps: []job.Step{mkdir, step}}, nil
	})
}

// ListOptions represents options for listing archive contents.
type ListOptions struct {
	Files []string
}

// List prints the contents of archive files without extracting them.
func List(opts ListOptions) error {
	failed := 0

	for _, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			failed++

			continue
		}

		ui.Message.Stepf("Contents of %s:", file)

		out, err := list(file)
		if err != nil {
			ui.Message.Errorf("Failed to list %s", file)
			ui.Message.Mutedf("%v", err)

			failed++

			continue
		}

		ui.Message.Mutedf("%s", strings.TrimRight(string(out), "\n"))
	}

	if failed > 0 {
		return derrors.Newf(
			derrors.CodeArchiveFailed,
			"%d of %d archives failed",
			failed,
			len(opts.Files),
		)
	}

	return nil
}

func list(file string) ([]byte, error) {
	archiveKind, err := detect(file)
	if err != nil {
		return nil, err
	}

	if archiveKind.zip {
		err = needTool("unzip", "unzip")
		if err != nil {
			return nil, err
		}

		return run(exec.CommandContext(context.Background(), "unzip", "-l", file))
	}

	if archiveKind.filter == "" {
		return run(exec.CommandContext(context.Background(), "tar", "-tf", file))
	}

	err = needTool(archiveKind.filter, archiveKind.filter)
	if err != nil {
		return nil, err
	}

	return pipe(
		exec.CommandContext(context.Background(), archiveKind.filter, "-dc", file),
		exec.CommandContext(context.Background(), "tar", "-tf", "-"),
	)
}

func run(cmd *exec.Cmd) ([]byte, error) {
	log.Debug("exec", "cmd", util.Describe(cmd))

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s: %w\n%s", cmd.Args[0], err, strings.TrimSpace(stderr.String()))
	}

	return out, nil
}
