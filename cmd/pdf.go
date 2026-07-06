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
	Short:   "PDF documents",
	Long:    `Compress and convert PDF documents using Ghostscript.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var pdfCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress PDFs using Ghostscript",
	Args:    cobra.MinimumNArgs(1),
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
	Args:    cobra.MinimumNArgs(1),
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
