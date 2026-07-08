//nolint:goconst
package cmd

import (
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/convert"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
)

// audioCmd represents the audio command.
var audioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"a"},
	Short:   "Compress and convert audio files",
	Long: `Compress and convert audio files using ffmpeg.

Input formats: wav, flac, aac, mp3, m4a, opus, ogg, wma
Output formats: mp3, aac, m4a, wav, flac, opus, ogg`,
	Example: `  uts audio compress podcast.wav -q low
  uts audio convert track.wav --to mp3`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// audioCompressCmd represents the audio compress command.
var audioCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress audio files using ffmpeg (aac)",
	Long: `Compress audio files using ffmpeg with AAC codec.

Quality: high (192k), medium (128k), low (96k), or raw kbps.
Output saved as <name>-small.m4a.`,
	Example: `  uts audio compress podcast.wav -q low
  uts audio compress voice-memo.m4a -q high
  uts audio compress voice.wav -q 256 --dry-run
  uts audio compress '*.wav' -r`,
	Args: cobra.MinimumNArgs(1),
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

// audioConvertCmd represents the audio convert command.
var audioConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between audio formats",
	Long: `Convert audio files (or extract audio from video) using ffmpeg.

Target formats: mp3, aac, m4a, wav, flac, opus, ogg`,
	Example: `  uts audio convert track.wav --to mp3
  uts audio convert video.mp4 --to mp3
  uts audio convert song.flac --to m4a -q high
  uts audio convert '*.wav' --to mp3 -q 96
  uts audio convert lecture.wav --to mp3`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			return derrors.New(
				derrors.CodeInvalidInput,
				"missing --to <format>. Examples: --to mp3, --to wav, --to flac, --to aac",
			)
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
