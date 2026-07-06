package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/convert"
)

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between formats",
	Long: `Convert files between different formats.

USAGE
  uts convert <subcategory> <input...> --to <format> [options]

SUBCATEGORIES
  image   Image format conversion (ImageMagick / sips)
  video   Video format conversion (ffmpeg)
  audio   Audio format conversion (ffmpeg)
  pdf     PDF <-> image conversion (pdftoppm / ImageMagick)

IMAGE EXAMPLES
  uts convert image photo.heic --to jpg
  uts convert image screenshot.png --to webp -q 85
  uts convert image '*.heic' --to jpg

VIDEO EXAMPLES
  uts convert video clip.mov --to mp4
  uts convert video recording.mkv --to webm -q 20

AUDIO EXAMPLES
  uts convert audio track.wav --to mp3 -q 96
  uts convert audio song.flac --to m4a

PDF EXAMPLES
  uts convert pdf report.pdf --to jpg
  uts convert pdf '*.jpg' '*.png' --to pdf`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var convertImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img", "i"},
	Short:   "Convert between image formats",
	Long: `Convert image files between formats.

USAGE
  uts convert image <input...> --to <format> [options]

TARGET FORMATS
  jpg, png, webp, gif, bmp, tiff, avif

EXAMPLES
  uts convert image photo.heic --to jpg
  uts convert image screenshot.png --to webp -q 85
  uts convert image '*.heic' --to jpg`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>")
			return nil
		}
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

var convertVideoCmd = &cobra.Command{
	Use:     "video",
	Aliases: []string{"v"},
	Short:   "Convert between video formats",
	Long: `Convert video files between formats.

USAGE
  uts convert video <input...> --to <format> [options]

TARGET FORMATS
  mp4, mkv, webm, mov, avi, flv

EXAMPLES
  uts convert video clip.mov --to mp4
  uts convert video recording.mkv --to webm -q 20`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>")
			return nil
		}
		return convert.Video(convert.VideoOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

var convertAudioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"a"},
	Short:   "Convert between audio formats",
	Long: `Convert audio files between formats.

USAGE
  uts convert audio <input...> --to <format> [options]

TARGET FORMATS
  mp3, aac, m4a, wav, flac, opus, ogg

EXAMPLES
  uts convert audio track.wav --to mp3 -q 96
  uts convert audio song.flac --to m4a`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>")
			return nil
		}
		return convert.Audio(convert.AudioOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			Quality:   quality,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

var convertPDFCmd = &cobra.Command{
	Use:     "pdf",
	Aliases: []string{"p"},
	Short:   "Convert PDF to/from images",
	Long: `Convert PDF documents to images or combine images into PDF.

USAGE
  uts convert pdf <input...> --to <format> [options]

TARGET FORMATS
  jpg, png (PDF -> images) or pdf (images -> PDF)

EXAMPLES
  uts convert pdf report.pdf --to jpg
  uts convert pdf '*.jpg' '*.png' --to pdf`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>")
			return nil
		}
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
	convertCmd.AddCommand(convertImageCmd)
	convertCmd.AddCommand(convertVideoCmd)
	convertCmd.AddCommand(convertAudioCmd)
	convertCmd.AddCommand(convertPDFCmd)
}
