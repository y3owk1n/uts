package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/convert"
)

var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img", "i"},
	Short:   "Images (png, jpg, webp, gif, bmp, tiff, heic, avif)",
	Long:    `Compress and convert image files using format-specific tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var imageCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress images using format-specific tools",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("compressing image files", "files", args)
		return compress.Image(compress.ImageOptions{
			Files:     args,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

var imageConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between image formats",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to jpg, --to webp, --to png")
			return nil
		}
		log.Debug("converting image files", "files", args, "target", targetFmt)
		return convert.Image(convert.ImageOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

func init() {
	imageCmd.AddCommand(imageCompressCmd)
	imageCmd.AddCommand(imageConvertCmd)
}
