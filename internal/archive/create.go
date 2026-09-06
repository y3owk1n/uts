//nolint:goconst,mnd
package archive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// CreateOptions represents options for archive creation.
type CreateOptions struct {
	Files     []string
	Algorithm string
	OutputDir string
	DryRun    bool
}

// Algorithms lists the accepted --algorithm values.
var Algorithms = []string{"zip", "gzip", "zstd", "xz", "brotli"}

// Create archives the given files and directories into a single archive
// named after the input.
func Create(opts CreateOptions) error {
	algo := normalizeAlgo(opts.Algorithm)

	ext, err := archiveExt(algo)
	if err != nil {
		return err
	}

	var inputs []string

	for _, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		inputs = append(inputs, file)
	}

	if len(inputs) == 0 {
		return derrors.New(derrors.CodeFileNotFound, "no input files found")
	}

	outDir := opts.OutputDir
	if outDir == "" {
		outDir = "."
	}

	output := filepath.Join(outDir, deriveName(inputs)+ext)

	ui.Message.Infof("Creating %s archive from %d input(s)", algo, len(inputs))

	producer, consumer, err := createCmds(algo, output, inputs)
	if err != nil {
		return err
	}

	desc := util.Describe(consumer)
	if producer != nil {
		desc = util.Describe(producer) + " | " + desc + " > " + output
	}

	if opts.DryRun {
		ui.Message.Infof("[dry-run] Would create %s", output)
		ui.Message.Mutedf("  $ %s", desc)

		return nil
	}

	err = os.MkdirAll(outDir, 0o755)
	if err != nil {
		return derrors.Wrap(err, derrors.CodeArchiveFailed, "create output directory")
	}

	spinner := ui.NewSpinner(nil, 0)
	spinner.SetSuffix(fmt.Sprintf("Archiving to %s...", output))
	spinner.Start()

	if producer != nil {
		err = pipeToFile(producer, consumer, output)
	} else {
		_, err = run(consumer)
	}

	spinner.Stop()

	if err != nil {
		_ = os.Remove(output)

		ui.Message.Errorf("Archive creation failed")
		ui.Message.Mutedf("%v", err)

		return derrors.Wrap(err, derrors.CodeArchiveFailed, "archive creation failed")
	}

	ui.Message.Successf("Created %s (%s)", output, util.HumanSize(util.FileSize(output)))

	return nil
}

func normalizeAlgo(algo string) string {
	switch strings.ToLower(algo) {
	case "gz", "gzip":
		return "gzip"
	case "zst", "zstd":
		return "zstd"
	case "br", "brotli":
		return "brotli"
	default:
		return strings.ToLower(algo)
	}
}

func archiveExt(algo string) (string, error) {
	switch algo {
	case "gzip":
		return ".tar.gz", nil
	case "zstd":
		return ".tar.zst", nil
	case "xz":
		return ".tar.xz", nil
	case "brotli":
		return ".tar.br", nil
	case "zip":
		return ".zip", nil
	}

	return "", derrors.Newf(derrors.CodeInvalidInput,
		"unknown algorithm %q (use one of: %s)", algo, strings.Join(Algorithms, ", "))
}

// createCmds returns the command(s) that build the archive. When producer is
// non-nil its stdout must be piped into consumer and consumer's stdout is the
// archive.
func createCmds(algo, output string, inputs []string) (*exec.Cmd, *exec.Cmd, error) {
	ctx := context.Background()

	tarArgs := func(flags string) []string {
		return append([]string{flags, output}, inputs...)
	}

	switch algo {
	case "gzip":
		return nil, exec.CommandContext(ctx, "tar", tarArgs("-czf")...), nil
	case "xz":
		return nil, exec.CommandContext(ctx, "tar", tarArgs("-cJf")...), nil
	case "zip":
		err := needTool("zip", "zip")
		if err != nil {
			return nil, nil, err
		}

		return nil, exec.CommandContext(
			ctx,
			"zip",
			append([]string{"-qr", output}, inputs...)...), nil
	case "zstd", "brotli":
		err := needTool(algo, algo)
		if err != nil {
			return nil, nil, err
		}

		producer := exec.CommandContext(ctx, "tar", append([]string{"-cf", "-"}, inputs...)...)
		consumer := exec.CommandContext(ctx, algo, "-c")

		return producer, consumer, nil
	}

	return nil, nil, derrors.Newf(derrors.CodeInvalidInput, "unknown algorithm %q", algo)
}

func pipeToFile(producer, consumer *exec.Cmd, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}

	out, err := pipe(producer, consumer)
	if err != nil {
		_ = file.Close()

		return err
	}

	_, err = file.Write(out)
	if err != nil {
		_ = file.Close()

		return err
	}

	return file.Close()
}

// deriveName picks the archive base name: a lone input's name, otherwise the
// name of the directory the inputs live in.
func deriveName(files []string) string {
	if len(files) == 1 {
		clean := filepath.Clean(files[0])

		info, err := os.Stat(clean)
		if err == nil && info.IsDir() {
			if base := filepath.Base(clean); base != "." && base != string(filepath.Separator) {
				return base
			}

			return "archive"
		}

		base := filepath.Base(clean)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		return strings.TrimSuffix(name, ".tar")
	}

	parent := filepath.Base(filepath.Dir(filepath.Clean(files[0])))
	if parent == "." || parent == string(filepath.Separator) {
		return "archive"
	}

	return parent
}
