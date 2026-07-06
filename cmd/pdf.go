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

ACTIONS
  compress  Compress PDFs using Ghostscript
  convert   Convert between PDF and images (pdftoppm, ImageMagick)

SUPPORTED CONVERSIONS
  PDF -> images:  jpg, png
  images -> PDF:  jpg, png, webp, gif, bmp, tiff

COMPRESSION EXAMPLES
  uts pdf compress thesis.pdf -q low
  uts pdf compress report.pdf -q medium -o ./web/
  uts pdf compress '*.pdf' -r
  uts pdf compress slides.pdf --dry-run

CONVERSION EXAMPLES
  uts pdf convert report.pdf --to jpg
  uts pdf convert slides.pdf --to png -q high
  uts pdf convert document.pdf --to jpg -q 200
  uts pdf convert '*.jpg' --to pdf
  uts pdf convert '*.jpg' '*.png' --to pdf`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pdfCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress PDFs using Ghostscript",
	Long: `Compress PDF documents using Ghostscript.

USAGE
  uts pdf compress <input...> [options]

QUALITY
  high       /printer  (300 DPI)
  medium     /ebook    (150 DPI) (default)
  low        /screen   (72 DPI)
  <dpi>      Numeric DPI (e.g. 150, 300, 400)

OUTPUT
  Files saved as <name>-small.pdf in the same directory.

EXAMPLES
  uts pdf compress thesis.pdf -q low
  uts pdf compress report.pdf -q medium -o ./web/
  uts pdf compress doc1.pdf doc2.pdf doc3.pdf -q low
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

USAGE
  uts pdf convert <input...> --to <format> [options]

DIRECTIONS
  PDF -> images:   pdftoppm or ImageMagick (outputs page-1.jpg, page-2.jpg, ...)
  images -> PDF:   ImageMagick (combines into single PDF)

QUALITY (PDF->images only)
  high       400 DPI
  medium     300 DPI (default)
  low        150 DPI
  <dpi>      Numeric DPI (e.g. 150, 300, 400)

OUTPUT
  PDF -> images:  Creates <basename>/ directory with page-*.ext files
  images -> PDF:  Creates <first-image-name>.pdf

EXAMPLES
  uts pdf convert report.pdf --to jpg
  uts pdf convert slides.pdf --to png -q high
  uts pdf convert document.pdf --to jpg -q 200
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
