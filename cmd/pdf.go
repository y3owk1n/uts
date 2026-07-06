package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/convert"
)

var pdfCmd = &cobra.Command{
	Use:     "pdf",
	Aliases: []string{"p"},
	Short:   "Compress and convert PDF documents",
	Long: `Compress and convert PDF documents using Ghostscript and ImageMagick.

Conversions: PDF ↔ jpg/png images.`,
	Example: `  uts pdf compress thesis.pdf -q low
  uts pdf convert report.pdf --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pdfCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress PDFs using Ghostscript",
	Long: `Compress PDF documents using Ghostscript.

Quality: high (300 DPI, /printer), medium (150 DPI, /ebook), low (72 DPI, /screen), or raw DPI.`,
	Example: `  uts pdf compress thesis.pdf -q low
  uts pdf compress report.pdf -q medium -o ./web/
  uts pdf compress slides.pdf -q 300 --dry-run
  uts pdf compress '*.pdf' -r`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("compressing PDF files", "files", args)
		return compress.PDF(compress.PDFOptions{
			Files:     args,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

var pdfConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between PDF and images",
	Long: `Convert PDF documents to images or combine images into a PDF.

PDF → images: creates <basename>/ directory with page-*.ext files (jpg/png)
images → PDF: combines images into a single PDF`,
	Example: `  uts pdf convert report.pdf --to jpg
  uts pdf convert slides.pdf --to png -q high
  uts pdf convert '*.jpg' '*.png' --to pdf
  uts pdf convert images/*.png --to pdf`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to jpg, --to png, --to pdf")
			return nil
		}
		log.Debug("converting PDF files", "files", args, "target", targetFmt)
		return convert.PDF(convert.PDFOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			Quality:   quality,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

func init() {
	pdfCmd.AddCommand(pdfCompressCmd)
	pdfCmd.AddCommand(pdfConvertCmd)
}
