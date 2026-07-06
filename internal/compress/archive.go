//nolint:mnd
package compress

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

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
		name := strings.TrimSuffix(base, filepath.Ext(base))

		return strings.TrimSuffix(name, ".tar")
	}

	parent := filepath.Dir(files[0])
	if parent == "." {
		return "archive"
	}

	return filepath.Base(parent)
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
		if !util.HasTool("zstd") {
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
		if !util.HasTool("brotli") {
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
		if !util.HasTool("zip") {
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
