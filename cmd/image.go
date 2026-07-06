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
	Short:   "Compress and convert image files",
	Long: `Compress and convert images using format-specific tools.

Input formats: png, jpg, jpeg, webp, gif, bmp, tiff, heic, heif, avif
Output formats: jpg, png, webp, gif, bmp, tiff, avif`,
	Example: `  uts image compress screenshot.png -q medium
  uts image convert photo.heic --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var imageCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress images using format-specific tools",
	Long: `Compress images using the best available tool for each format.

Tools by format: png (pngquant+optipng), jpg (jpegoptim),
webp (cwebp), gif (gifsicle), bmp/tiff (ImageMagick),
heic (heif-convert), avif (cavif/avifenc).

HEIC files are converted to JPEG.`,
	Example: `  uts image compress screenshot.png -q medium
  uts image compress logo.jpg -q high -i
  uts image compress '*.png' -r
  uts image compress photo.webp -q 75 --dry-run -v`,
	Args: cobra.MinimumNArgs(1),
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
	Long: `Convert image files between formats using ImageMagick or sips.

Target formats: jpg, png, webp, gif, bmp, tiff, avif`,
	Example: `  uts image convert photo.heic --to jpg
  uts image convert screenshot.png --to webp -q high
  uts image convert photo.jpg --to avif -q 70
  uts image convert '*.heic' --to jpg`,
	Args: cobra.MinimumNArgs(1),
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
