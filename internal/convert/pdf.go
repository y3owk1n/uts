package convert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type PDFOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	DryRun    bool
}

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

	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		origSize := util.FileSize(file)
		ui.Message.Stepf("[%d/%d] %s (%s)", i+1, len(opts.Files), file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would extract %s -> images (dpi=%d)", file, dpi)
			continue
		}

		outDir := opts.OutputDir
		if outDir == "" {
			base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
			outDir = filepath.Join(filepath.Dir(file), base)
		}
		os.MkdirAll(outDir, 0755)

		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Extracting %s...", file))
		sp.Start()

		var convertErr error
		if hasTool("pdftoppm") {
			imgExt := target
			if target == "jpg" {
				imgExt = "jpeg"
			}
			convertErr = exec.Command("pdftoppm", "-"+imgExt, "-r", fmt.Sprintf("%d", dpi), file, filepath.Join(outDir, "page")).Run()
		} else if hasMagick() {
			magick := "magick"
			if !hasTool("magick") {
				magick = "convert"
			}
			convertErr = exec.Command(magick, "-density", fmt.Sprintf("%d", dpi), file, filepath.Join(outDir, "page.%03d."+target)).Run()
		} else {
			sp.Stop()
			ui.Message.Errorf("No PDF conversion tool found. Install: brew install poppler imagemagick")
			return nil
		}
		sp.Stop()

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

	os.MkdirAll(filepath.Dir(outPath), 0755)
	sp := ui.NewSpinner(nil, 0)
	sp.SetSuffix("Combining images...")
	sp.Start()

	var convertErr error
	if hasMagick() {
		magick := "magick"
		if !hasTool("magick") {
			magick = "convert"
		}
		args := append(validFiles, outPath)
		convertErr = exec.Command(magick, args...).Run()
	} else {
		sp.Stop()
		ui.Message.Errorf("ImageMagick not found. Install: brew install imagemagick")
		return nil
	}
	sp.Stop()

	if convertErr == nil && util.FileExists(outPath) {
		ui.Message.Successf("%d images → %s (%s)", len(validFiles), filepath.Base(outPath), util.HumanSize(util.FileSize(outPath)))
	} else {
		ui.Message.Errorf("PDF creation failed")
	}

	return nil
}
