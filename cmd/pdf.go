//nolint:goconst
package cmd

import (
	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/format"
)

// pdfCmd represents the PDF command.
var pdfCmd = &cobra.Command{
	Use:     "pdf",
	Aliases: []string{"p"},
	Short:   "Compress and convert PDF documents",
	Long: `Compress and convert PDF documents using Ghostscript, poppler and ImageMagick.

Conversions: PDF ↔ jpg/png images.`,
	Example: `  uts pdf compress thesis.pdf -q low
  uts pdf convert report.pdf --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// pdfCompressCmd represents the pdf compress command.
var pdfCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress PDFs using Ghostscript",
	Long: `Compress PDF documents using Ghostscript.

Quality: high (400 DPI, /printer), medium (300 DPI, /ebook), low (150 DPI, /screen), or raw DPI.`,
	Example: `  uts pdf compress thesis.pdf -q low
  uts pdf compress report.pdf -q medium -o ./web/
  uts pdf compress slides.pdf -q 300 --dry-run
  uts pdf compress ./documents -r`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := inputs(args, format.PDFExts)
		if err != nil {
			return err
		}

		log.Debug("compressing PDF files", "files", files)

		return compress.PDF(compress.PDFOptions{
			Files:     files,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

func init() {
	pdfCmd.AddCommand(pdfCompressCmd)
	pdfCmd.AddCommand(newPDFConvertCmd("convert", []string{"x"}))
}
