//nolint:mnd,goconst
package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
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

// PDF converts PDFs to images or combines images into a PDF, chosen by the
// extension of the first input.
func PDF(opts PDFOptions) error {
	if len(opts.Files) == 0 {
		return derrors.New(derrors.CodeInvalidInput, "no files provided")
	}

	target := format.Normalize(opts.Target)
	if !slices.Contains(format.PDFTargets, target) {
		return unsupportedTarget(target, format.PDFTargets)
	}

	first := format.Ext(opts.Files[0])

	switch {
	case first == "pdf" && target != "pdf":
		return pdfToImages(opts, target)
	case format.Classify(first) == format.Image && target == "pdf":
		return imagesToPDF(opts)
	case first == "pdf":
		return derrors.New(
			derrors.CodeInvalidInput,
			"input is already a PDF (use --to jpg or --to png)",
		)
	default:
		return derrors.Newf(derrors.CodeUnsupportedFormat,
			"unsupported input .%s (provide a PDF, or images with --to pdf)", first)
	}
}

func pdfToImages(opts PDFOptions, target string) error {
	dpi, _, err := util.PDFDPI(opts.Quality)
	if err != nil {
		return err
	}

	if !util.HasTool("pdftoppm") && util.MagickBin() == "" {
		return derrors.New(derrors.CodeToolNotFound,
			"PDF conversion tools not found — install: brew install poppler imagemagick")
	}

	ui.Message.Infof("Converting PDF → .%s pages at %d DPI", target, dpi)

	return job.Run(opts.Files, job.Options{
		Verb:   "Extracting",
		Done:   "Extracted",
		Noun:   "PDF file",
		DryRun: opts.DryRun,
		Code:   derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		if format.Ext(file) != "pdf" {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"not a PDF: .%s",
				format.Ext(file),
			)
		}

		base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		outDir := filepath.Join(filepath.Dir(file), base)
		if opts.OutputDir != "" {
			outDir = filepath.Join(opts.OutputDir, base)
		}

		mkdir := job.Do("mkdir -p "+outDir, func() error { return os.MkdirAll(outDir, 0o755) })

		if util.HasTool("pdftoppm") {
			flag := "-" + target
			if target == "jpg" {
				flag = "-jpeg"
			}

			return &job.Job{
				Input:  file,
				Output: outDir,
				Steps: []job.Step{mkdir, job.Exec("pdftoppm", flag, "-r", strconv.Itoa(dpi),
					file, filepath.Join(outDir, "page"))},
				Note: "via pdftoppm",
			}, nil
		}

		return &job.Job{
			Input:  file,
			Output: outDir,
			Steps: []job.Step{mkdir, job.Exec(util.MagickBin(), "-density", strconv.Itoa(dpi),
				file, filepath.Join(outDir, "page-%03d."+target))},
			Note: "via ImageMagick",
		}, nil
	})
}

func imagesToPDF(opts PDFOptions) error {
	bin := util.MagickBin()
	if bin == "" {
		return derrors.New(
			derrors.CodeToolNotFound,
			"ImageMagick not found — install: brew install imagemagick",
		)
	}

	var images []string

	for _, file := range opts.Files {
		switch {
		case !util.FileExists(file):
			ui.Message.Warnf("File not found: %s", file)
		case format.Classify(format.Ext(file)) != format.Image:
			ui.Message.Warnf("Skipping non-image: %s", file)
		default:
			images = append(images, file)
		}
	}

	if len(images) == 0 {
		return derrors.New(derrors.CodeInvalidInput, "no valid image files provided")
	}

	first := images[0]

	out := util.CalcConvertOutputPath(first, "pdf", opts.OutputDir)

	ui.Message.Infof("Combining %d images into %s", len(images), out)

	return job.Run([]string{first}, job.Options{
		Verb:   "Combining",
		Done:   "Combined",
		Noun:   "PDF",
		DryRun: opts.DryRun,
		Code:   derrors.CodeConversionFailed,
	}, func(string) (*job.Job, error) {
		args := make([]string, 0, len(images)+1)
		args = append(args, images...)
		args = append(args, out)

		return &job.Job{
			Input:  first,
			Output: out,
			Steps:  []job.Step{job.Exec(bin, args...)},
			Note:   fmt.Sprintf("%d images via ImageMagick", len(images)),
		}, nil
	})
}
