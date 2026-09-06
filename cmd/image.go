//nolint:goconst
package cmd

import (
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/format"
)

// imageCmd represents the image command.
var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img", "i"},
	Short:   "Compress and convert image files",
	Long: `Compress and convert images using format-specific tools.

Input formats: ` + strings.Join(format.ImageExts, ", ") + `
Output formats: ` + strings.Join(format.ImageTargets, ", "),
	Example: `  uts image compress screenshot.png -q medium
  uts image convert photo.heic --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// imageCompressCmd represents the image compress command.
var imageCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress images using format-specific tools",
	Long: `Compress images using the best available tool for each format.

Tools by format: png (pngquant+optipng), jpg (jpegoptim), webp (cwebp),
gif (gifsicle), heic (heif-convert), bmp/tiff/avif (ImageMagick).
ImageMagick is the fallback whenever a dedicated tool is missing.

HEIC files are converted to JPEG. Results that are not smaller than the
original are reported and, with -i, the original is kept.`,
	Example: `  uts image compress screenshot.png -q medium
  uts image compress logo.jpg -q high -i
  uts image compress ./photos -r
  uts image compress photo.jpg --max 2000
  uts image compress '*.png' -q 75 --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := checkMaxEdge()
		if err != nil {
			return err
		}

		files, err := inputs(args, format.ImageExts)
		if err != nil {
			return err
		}

		log.Debug("compressing image files", "files", files)

		return compress.Image(compress.ImageOptions{
			Files:     files,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
			MaxEdge:   maxEdge,
		})
	},
}

func init() {
	addMaxFlag(imageCompressCmd)
	imageCmd.AddCommand(imageCompressCmd)
	imageCmd.AddCommand(newImageConvertCmd("convert", []string{"x"}))
}
