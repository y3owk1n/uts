//nolint:mnd
package archive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
)

// ExtractOptions represents options for archive extraction.
type ExtractOptions struct {
	Files     []string
	OutputDir string
	DryRun    bool
}

// ListOptions represents options for listing archive contents.
type ListOptions struct {
	Files []string
}

// Extract extracts archive files.
func Extract(opts ExtractOptions) error {
	for _, file := range opts.Files {
		if !fileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))

		outDir := opts.OutputDir
		if outDir == "" {
			outDir = "."
		}

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would extract %s -> %s/", file, outDir)

			continue
		}

		ui.Message.Stepf("Extracting %s → %s/", file, outDir)
		_ = os.MkdirAll(outDir, 0o755)

		spinner := ui.NewSpinner(nil, 0)
		spinner.SetSuffix(fmt.Sprintf("Extracting %s...", file))
		spinner.Start()

		var err error
		switch ext {
		case "zip":
			if !hasTool("unzip") {
				spinner.Stop()
				ui.Message.Errorf("unzip not found — install: brew install unzip")

				continue
			}

			err = exec.CommandContext(context.Background(), "unzip", "-qo", file, "-d", outDir).
				Run()
		case "tar":
			err = exec.CommandContext(context.Background(), "tar", "xf", file, "-C", outDir).Run()
		case "gz", "tgz":
			err = exec.CommandContext(context.Background(), "tar", "xzf", file, "-C", outDir).Run()
		case "zst", "zstd":
			if !hasTool("zstd") {
				spinner.Stop()
				ui.Message.Errorf("zstd not found — install: brew install zstd")

				continue
			}

			decomp := exec.CommandContext(
				context.Background(),
				"zstd",
				"-d",
				file,
				"--force",
				"-o",
				file+".unpacked",
			)

			err = decomp.Run()
			if err != nil {
				spinner.Stop()
				ui.Message.Errorf("zstd decompression failed: %s", file)

				continue
			}

			err = exec.CommandContext(context.Background(), "tar", "xf", file+".unpacked", "-C", outDir).
				Run()
			_ = os.Remove(file + ".unpacked")
		case "xz", "txz":
			err = exec.CommandContext(context.Background(), "xz", "-dk", file).Run()
			if err == nil {
				err = exec.CommandContext(context.Background(), "tar", "xf", strings.TrimSuffix(file, ".xz"), "-C", outDir).
					Run()
				_ = os.Remove(strings.TrimSuffix(file, ".xz"))
			}
		case "bz2", "tbz2":
			err = exec.CommandContext(context.Background(), "bunzip2", "-k", file).Run()
			if err == nil {
				err = exec.CommandContext(context.Background(), "tar", "xf", strings.TrimSuffix(file, ".bz2"), "-C", outDir).
					Run()
				_ = os.Remove(strings.TrimSuffix(file, ".bz2"))
			}
		default:
			spinner.Stop()
			ui.Message.Errorf("Unsupported archive: .%s", ext)

			continue
		}

		spinner.Stop()

		if err == nil {
			ui.Message.Successf("Extracted: %s -> %s/", file, outDir)
		} else {
			ui.Message.Errorf("Extraction failed: %s", file)
		}
	}

	return nil
}

// List lists the contents of archive files.
func List(opts ListOptions) error {
	for _, file := range opts.Files {
		if !fileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		ui.Message.Stepf("Contents of %s:", file)

		spinner := ui.NewSpinner(nil, 0)
		spinner.SetSuffix(fmt.Sprintf("Listing %s...", file))
		spinner.Start()

		var (
			output []byte
			err    error
		)
		switch ext {
		case "zip":
			if !hasTool("unzip") {
				spinner.Stop()
				ui.Message.Errorf("unzip not found — install: brew install unzip")

				continue
			}

			output, err = exec.CommandContext(context.Background(), "unzip", "-l", file).Output()
		case "tar":
			output, err = exec.CommandContext(context.Background(), "tar", "tf", file).Output()
		case "gz", "tgz":
			output, err = exec.CommandContext(context.Background(), "tar", "tzf", file).Output()
		case "zst", "zstd":
			if hasTool("zstd") {
				output, err = exec.CommandContext(context.Background(), "zstd", "-dc", file).
					Output()
			} else {
				spinner.Stop()
				ui.Message.Errorf("zstd not found — install: brew install zstd")

				continue
			}
		case "xz", "txz":
			output, err = exec.CommandContext(context.Background(), "xz", "-dc", file).Output()
		case "bz2", "tbz2":
			output, err = exec.CommandContext(context.Background(), "bzip2", "-dc", file).Output()
		default:
			spinner.Stop()
			ui.Message.Errorf("Unsupported archive: .%s", ext)

			continue
		}

		spinner.Stop()

		if err == nil {
			ui.Message.Mutedf("%s", string(output))
		} else {
			ui.Message.Errorf("Failed to list: %s", file)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func hasTool(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}
