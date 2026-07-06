package archive

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/style"
)

type ExtractOptions struct {
	Files     []string
	OutputDir string
	DryRun    bool
}

type ListOptions struct {
	Files []string
}

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
		os.MkdirAll(outDir, 0755)

		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Extracting %s...", file))
		sp.Start()

		var err error
		switch ext {
		case "zip":
			if !hasTool("unzip") {
				sp.Stop()
				ui.Message.Errorf("unzip not found — install: brew install unzip")
				continue
			}
			err = exec.Command("unzip", "-qo", file, "-d", outDir).Run()
		case "tar":
			err = exec.Command("tar", "xf", file, "-C", outDir).Run()
		case "gz", "tgz":
			err = exec.Command("tar", "xzf", file, "-C", outDir).Run()
		case "zst", "zstd":
			if !hasTool("zstd") {
				sp.Stop()
				ui.Message.Errorf("zstd not found — install: brew install zstd")
				continue
			}
			decomp := exec.Command("zstd", "-d", file, "--force", "-o", file+".unpacked")
			if err = decomp.Run(); err != nil {
				sp.Stop()
				ui.Message.Errorf("zstd decompression failed: %s", file)
				continue
			}
			err = exec.Command("tar", "xf", file+".unpacked", "-C", outDir).Run()
			os.Remove(file + ".unpacked")
		case "xz", "txz":
			err = exec.Command("xz", "-dk", file).Run()
			if err == nil {
				err = exec.Command("tar", "xf", strings.TrimSuffix(file, ".xz"), "-C", outDir).Run()
				os.Remove(strings.TrimSuffix(file, ".xz"))
			}
		case "bz2", "tbz2":
			err = exec.Command("bunzip2", "-k", file).Run()
			if err == nil {
				err = exec.Command("tar", "xf", strings.TrimSuffix(file, ".bz2"), "-C", outDir).Run()
				os.Remove(strings.TrimSuffix(file, ".bz2"))
			}
		default:
			sp.Stop()
			ui.Message.Errorf("Unsupported archive: .%s", ext)
			continue
		}
		sp.Stop()

		if err == nil {
			ui.Message.Successf("Extracted: %s -> %s/", file, outDir)
		} else {
			ui.Message.Errorf("Extraction failed: %s", file)
		}
	}
	return nil
}

func List(opts ListOptions) error {
	palette := style.Default()
	for _, file := range opts.Files {
		if !fileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		ui.Message.Stepf("Contents of %s:", file)

		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Listing %s...", file))
		sp.Start()

		var output []byte
		var err error
		switch ext {
		case "zip":
			if !hasTool("unzip") {
				sp.Stop()
				ui.Message.Errorf("unzip not found — install: brew install unzip")
				continue
			}
			output, err = exec.Command("unzip", "-l", file).Output()
		case "tar":
			output, err = exec.Command("tar", "tf", file).Output()
		case "gz", "tgz":
			output, err = exec.Command("tar", "tzf", file).Output()
		case "zst", "zstd":
			if hasTool("zstd") {
				output, err = exec.Command("zstd", "-dc", file).Output()
			} else {
				sp.Stop()
				ui.Message.Errorf("zstd not found — install: brew install zstd")
				continue
			}
		case "xz", "txz":
			output, err = exec.Command("xz", "-dc", file).Output()
		case "bz2", "tbz2":
			output, err = exec.Command("bzip2", "-dc", file).Output()
		default:
			sp.Stop()
			ui.Message.Errorf("Unsupported archive: .%s", ext)
			continue
		}
		sp.Stop()

		if err == nil {
			fmt.Print(lipgloss.NewStyle().
				Foreground(palette.Text).
				Render(string(output)) + "\n")
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
