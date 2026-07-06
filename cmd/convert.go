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
Subcategories: image, video, audio, pdf

Examples:
  uts convert image photo.heic --to jpg
  uts convert video clip.mov --to mp4
  uts convert audio track.wav --to mp3
  uts convert pdf report.pdf --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var convertImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"img", "i"},
	Short:   "Convert between image formats",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to jpg, --to webp, --to png")
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
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to mp4, --to mkv, --to webm")
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
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to mp3, --to wav, --to flac")
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
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to jpg, --to png, --to pdf")
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
