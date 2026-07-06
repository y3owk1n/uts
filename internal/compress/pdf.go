package compress

import (
	"fmt"
	"os/exec"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type PDFOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

func PDF(opts PDFOptions) error {
	dpi, settings, err := util.PDFDPI(opts.Quality)
	if err != nil {
		return err
	}

	if settings != "" {
		ui.Message.Infof("PDF compression at %s quality (preset=%s, %d DPI)", opts.Quality, settings, dpi)
	} else {
		ui.Message.Infof("PDF compression at %s quality (%d DPI)", opts.Quality, dpi)
	}

	total := len(opts.Files)
	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		out := util.OutputPath(file, "small")
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", i+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would compress %s -> %s (settings=%s)", file, out, settings)
			continue
		}

		util.EnsureDir(out)
		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Compressing %s...", file))
		sp.Start()

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

		output, err := exec.Command("gs", args...).CombinedOutput()
		sp.Stop()

		if err == nil && util.FileExists(out) {
			newSize := util.FileSize(out)
			ratio := util.CompressionRatio(origSize, newSize)
			ui.Message.Successf("%s: %s → %s %s", file, util.HumanSize(origSize), util.HumanSize(newSize), ratio)
		} else {
			ui.Message.Errorf("Compression failed: %s", file)
			ui.Message.Mutedf("gs: %s", string(output))
		}
	}

	if total > 1 {
		ui.Message.Successf("Compressed %d PDF files", total)
	}
	return nil
}
