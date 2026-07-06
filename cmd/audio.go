package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/convert"
)

var audioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"a"},
	Short:   "Audio files (wav, flac, aac, mp3, m4a, opus, ogg)",
	Long:    `Compress and convert audio files using ffmpeg.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var audioCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress audio files using ffmpeg (aac)",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("compressing audio files", "files", args)
		return compress.Audio(compress.AudioOptions{
			Files:     args,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

var audioConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between audio formats",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to mp3, --to wav, --to flac, --to aac")
			return nil
		}
		log.Debug("converting audio files", "files", args, "target", targetFmt)
		return convert.Audio(convert.AudioOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			Quality:   quality,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

func init() {
	audioCmd.AddCommand(audioCompressCmd)
	audioCmd.AddCommand(audioConvertCmd)
}
