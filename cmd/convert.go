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
	Long: `Convert files between different formats (image, video, audio, pdf).

Requires --to <format> flag to specify the target format.`,
	Example: `  uts convert image photo.heic --to jpg
  uts convert video clip.mov --to mp4
  uts convert audio track.wav --to mp3 -q 96
  uts convert pdf report.pdf --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var convertImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img", "i"},
	Short:   "Convert between image formats",
	Long: `Convert image files between formats.

Target formats: jpg, png, webp, gif, bmp, tiff, avif`,
	Example: `  uts convert image photo.heic --to jpg
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

Target formats: mp4, mkv, webm, mov, avi, flv`,
	Example: `  uts convert video clip.mov --to mp4
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

Target formats: mp3, aac, m4a, wav, flac, opus, ogg`,
	Example: `  uts convert audio track.wav --to mp3 -q 96
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

Target formats: jpg, png (PDF→images), pdf (images→PDF)`,
	Example: `  uts convert pdf report.pdf --to jpg
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
