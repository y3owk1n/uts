//nolint:goconst
package compress

import (
	"fmt"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// PDFOptions represents options for PDF compression.
type PDFOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// PDF compresses PDF files using Ghostscript.
func PDF(opts PDFOptions) error {
	dpi, settings, err := util.PDFDPI(opts.Quality)
	if err != nil {
		return err
	}

	if !util.HasTool("gs") {
		return derrors.New(
			derrors.CodeToolNotFound,
			"ghostscript not found — install: brew install ghostscript",
		)
	}

	if settings != "" {
		ui.Message.Infof(
			"PDF compression at %s quality (preset=%s, %d DPI)",
			opts.Quality,
			settings,
			dpi,
		)
	} else {
		ui.Message.Infof("PDF compression at %s quality (%d DPI)", opts.Quality, dpi)
	}

	return job.Run(opts.Files, job.Options{
		Verb:    "Compressing",
		Done:    "Compressed",
		Noun:    "PDF files",
		Compare: true,
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeCompressionFailed,
	}, func(file string) (*job.Job, error) {
		if format.Ext(file) != "pdf" {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"not a PDF: .%s",
				format.Ext(file),
			)
		}

		out := util.CalcOutputPath(file, "small", opts.OutputDir)

		args := []string{
			"-sDEVICE=pdfwrite",
			"-dCompatibilityLevel=1.4",
			"-dNOPAUSE",
			"-dQUIET",
			"-dBATCH",
		}
		if settings != "" {
			args = append(args, "-dPDFSETTINGS="+settings)
		}

		args = append(args,
			fmt.Sprintf("-dColorImageResolution=%d", dpi),
			fmt.Sprintf("-dGrayImageResolution=%d", dpi),
			fmt.Sprintf("-dMonoImageResolution=%d", dpi),
			"-sOutputFile="+out, file,
		)

		return &job.Job{
			Input:  file,
			Output: out,
			Steps:  []job.Step{job.Exec("gs", args...)},
			Note:   "via ghostscript",
		}, nil
	})
}
