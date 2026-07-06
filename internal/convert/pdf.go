//nolint:goconst,mnd
package convert

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// PDFOptions represents options for PDF conversion.
type PDFOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	DryRun    bool
}

// PDF converts PDFs to images or images to PDFs.
func PDF(opts PDFOptions) error {
	if len(opts.Files) == 0 {
		ui.Message.Errorf("No files provided")

		return nil
	}

	firstExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(opts.Files[0]), "."))

	if firstExt == "pdf" {
		return pdfToImages(opts)
	}

	switch firstExt {
	case "jpg", "jpeg", "png", "webp", "gif", "bmp", "tiff", "tif":
		return imagesToPDF(opts)
	default:
		ui.Message.Errorf("Unsupported input: .%s (provide PDF or images)", firstExt)

		return nil
	}
}

func pdfToImages(opts PDFOptions) error {
	target := strings.ToLower(opts.Target)
	switch target {
	case "jpg", "jpeg", "png":
	default:
		ui.Message.Errorf("Unsupported target: .%s (use jpg or png)", target)

		return nil
	}

	dpi, _, err := util.PDFDPI(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof("Converting PDF → images at %d DPI", dpi)

	for idx, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		origSize := util.FileSize(file)
		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, len(opts.Files), file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would extract %s -> images (dpi=%d)", file, dpi)

			continue
		}

		outDir := opts.OutputDir
		if outDir == "" {
			base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			outDir = filepath.Join(filepath.Dir(file), base)
		}

		_ = os.MkdirAll(outDir, 0o755)

		spinner := ui.NewSpinner(nil, 0)
		spinner.SetSuffix(fmt.Sprintf("Extracting %s...", file))
		spinner.Start()

		var convertErr error
		switch {
		case hasTool("pdftoppm"):
			imgExt := target
			if target == "jpg" {
				imgExt = "jpeg"
			}

			convertErr = exec.CommandContext(context.Background(), "pdftoppm", "-"+imgExt, "-r", strconv.Itoa(dpi), file, filepath.Join(outDir, "page")).
				Run()
		case hasMagick():
			magick := "magick"
			if !hasTool("magick") {
				magick = "convert"
			}

			convertErr = exec.CommandContext(context.Background(), magick, "-density", strconv.Itoa(dpi), file, filepath.Join(outDir, "page.%03d."+target)).
				Run()
		default:
			spinner.Stop()
			ui.Message.Errorf(
				"PDF conversion tools not found — install: brew install poppler imagemagick",
			)

			return nil
		}

		spinner.Stop()

		if convertErr == nil {
			ui.Message.Successf("%s: pages extracted to %s/", file, outDir)
		} else {
			ui.Message.Errorf("Extraction failed: %s", file)
		}
	}

	return nil
}

func imagesToPDF(opts PDFOptions) error {
	target := strings.ToLower(opts.Target)
	if target != "pdf" {
		ui.Message.Errorf("Cannot combine images into .%s", target)

		return nil
	}

	ui.Message.Infof("Combining %d images into PDF", len(opts.Files))

	var validFiles []string
	for _, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		switch ext {
		case "jpg", "jpeg", "png", "webp", "gif", "bmp", "tiff", "tif":
			validFiles = append(validFiles, file)
		default:
			ui.Message.Warnf("Skipping non-image: %s", file)
		}
	}

	if len(validFiles) == 0 {
		ui.Message.Errorf("No valid image files provided")

		return nil
	}

	out := opts.OutputDir
	if out == "" {
		out = filepath.Dir(validFiles[0])
	}

	baseName := strings.TrimSuffix(filepath.Base(validFiles[0]), filepath.Ext(validFiles[0]))
	outPath := filepath.Join(out, baseName+".pdf")

	if opts.DryRun {
		ui.Message.Infof("[dry-run] Would combine %d images -> %s", len(validFiles), outPath)

		return nil
	}

	_ = os.MkdirAll(filepath.Dir(outPath), 0o755)

	ui.Message.Stepf("Combining %d images → PDF", len(validFiles))

	spinner := ui.NewSpinner(nil, 0)
	spinner.SetSuffix("Combining images...")
	spinner.Start()

	var convertErr error
	if hasMagick() {
		magick := "magick"
		if !hasTool("magick") {
			magick = "convert"
		}

		args := make([]string, 0, len(validFiles)+1)
		args = append(args, validFiles...)
		args = append(args, outPath)
		convertErr = exec.CommandContext(context.Background(), magick, args...).Run()
	} else {
		spinner.Stop()
		ui.Message.Errorf("ImageMagick not found — install: brew install imagemagick")

		return nil
	}

	spinner.Stop()

	if convertErr == nil && util.FileExists(outPath) {
		ui.Message.Successf(
			"%d images → %s (%s)",
			len(validFiles),
			filepath.Base(outPath),
			util.HumanSize(util.FileSize(outPath)),
		)
	} else {
		ui.Message.Errorf("PDF creation failed")
	}

	return nil
}
