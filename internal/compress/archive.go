//nolint:goconst,mnd
package compress

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

var errNoCompressionTools = derrors.New(derrors.CodeToolNotFound, "no compression tools available")

// ArchiveOptions represents options for archive creation.
type ArchiveOptions struct {
	Files     []string
	Algorithm string
	OutputDir string
	DryRun    bool
}

// Archive creates an archive of the given files.
func Archive(opts ArchiveOptions) error {
	ui.Message.Infof("Creating archive with %s algorithm", opts.Algorithm)

	name := deriveArchiveName(opts.Files)

	outDir := opts.OutputDir
	if outDir == "" {
		outDir = "."
	}

	if opts.DryRun {
		ui.Message.Infof("[dry-run] Would create archive from: %v", opts.Files)
		ui.Message.Infof("[dry-run] Output: %s", outDir)

		return nil
	}

	_ = os.MkdirAll(outDir, 0o755)

	if opts.Algorithm == "auto" {
		return archiveAuto(outDir, name, opts.Files)
	}

	archiveWith(opts.Algorithm, outDir, name, opts.Files)

	return nil
}

func deriveArchiveName(files []string) string {
	if len(files) == 1 {
		info, err := os.Stat(files[0])
		if err == nil && info.IsDir() {
			return filepath.Base(files[0])
		}

		base := filepath.Base(files[0])

		return stringsTrimSuffix(base, filepath.Ext(base))
	}

	parent := filepath.Dir(files[0])
	if parent == "." {
		return "archive"
	}

	return filepath.Base(parent)
}

func archiveAuto(outDir, name string, files []string) error {
	var bestAlgo, bestFile string

	bestSize := int64(999999999999)

	spinner := ui.NewSpinner(nil, 0)
	spinner.Start()

	algorithms := []string{"zstd", "xz", "brotli", "gzip"}
	for _, algo := range algorithms {
		spinner.SetSuffix(fmt.Sprintf("Trying %s...", algo))

		candidate := archiveWith(algo, outDir, name, files)
		if candidate != "" {
			size := util.FileSize(candidate)
			if size < bestSize {
				if bestFile != "" {
					_ = os.Remove(bestFile)
				}

				bestSize = size
				bestAlgo = algo
				bestFile = candidate
			} else {
				_ = os.Remove(candidate)
			}
		}
	}

	spinner.Stop()

	if bestFile != "" {
		ui.Message.Successf(
			"Best algorithm: %s → %s (%s)",
			bestAlgo,
			filepath.Base(bestFile),
			util.HumanSize(bestSize),
		)

		return nil
	}

	ui.Message.Errorf("No compression tools available — install: brew install zstd xz brotli gzip")

	return errNoCompressionTools
}

func archiveWith(algo, outDir, name string, files []string) string {
	var output string
	switch algo {
	case "gzip", "gz":
		output = filepath.Join(outDir, name+".tar.gz")
	case "zstd", "zst":
		output = filepath.Join(outDir, name+".tar.zst")
	case "xz":
		output = filepath.Join(outDir, name+".tar.xz")
	case "brotli", "br":
		output = filepath.Join(outDir, name+".tar.br")
	case "zip":
		output = filepath.Join(outDir, name+".zip")
	default:
		ui.Message.Errorf("Unknown algorithm: %s", algo)

		return ""
	}

	var cmd *exec.Cmd
	switch algo {
	case "gzip", "gz":
		cmd = exec.CommandContext(
			context.Background(),
			"tar",
			append([]string{"-czf", output}, files...)...,
		)
	case "zstd", "zst":
		if !hasTool("zstd") {
			return ""
		}

		cmd = exec.CommandContext(
			context.Background(),
			"tar",
			append([]string{"--zstd", "-cf", output}, files...)...,
		)
	case "xz":
		cmd = exec.CommandContext(
			context.Background(),
			"tar",
			append([]string{"-cJf", output}, files...)...,
		)
	case "brotli", "br":
		if !hasTool("brotli") {
			return ""
		}

		tarCmd := exec.CommandContext(
			context.Background(),
			"tar",
			append([]string{"-cf", "-"}, files...)...,
		)

		file, err := os.Create(output)
		if err != nil {
			return ""
		}
		defer file.Close() //nolint:errcheck

		brotliCmd := exec.CommandContext(context.Background(), "brotli", "-c")

		var pipeErr error

		brotliCmd.Stdin, pipeErr = tarCmd.StdoutPipe()
		if pipeErr != nil {
			return ""
		}

		brotliCmd.Stdout = file

		_ = tarCmd.Start()
		_ = brotliCmd.Run()
		_ = tarCmd.Wait()

		ui.Message.Successf(
			"Created %s (%s)",
			filepath.Base(output),
			util.HumanSize(util.FileSize(output)),
		)

		return output
	case "zip":
		if !hasTool("zip") {
			return ""
		}

		cmd = exec.CommandContext(
			context.Background(),
			"zip",
			append([]string{"-r", output}, files...)...,
		)
	}

	if cmd != nil {
		output, err := cmd.CombinedOutput()
		if err != nil {
			ui.Message.Mutedf("archive error: %s", string(output))

			return ""
		}
	}

	if util.FileExists(output) {
		ui.Message.Successf(
			"Created %s (%s)",
			filepath.Base(output),
			util.HumanSize(util.FileSize(output)),
		)
	}

	return output
}

func stringsTrimSuffix(str, suffix string) string {
	if len(str) >= len(suffix) && str[len(str)-len(suffix):] == suffix {
		return str[:len(str)-len(suffix)]
	}

	return str
}
