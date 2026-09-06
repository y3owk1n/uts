//nolint:goconst
package cmd

import (
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/format"
)

// audioCmd represents the audio command.
var audioCmd = &cobra.Command{
	Use:     "audio",
	Aliases: []string{"a"},
	Short:   "Compress and convert audio files",
	Long: `Compress and convert audio files using ffmpeg.

Input formats: ` + strings.Join(format.AudioExts, ", ") + `
Output formats: ` + strings.Join(format.AudioTargets, ", "),
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
	Short:   "Compress audio files using ffmpeg",
	Long: `Compress audio files using ffmpeg.

Lossy inputs (mp3, ogg, opus, m4a) are re-encoded in their own format.
Lossless inputs (wav, flac) and everything else become AAC in .m4a.

Quality: high (192k), medium (128k), low (96k), or raw kbps.`,
	Example: `  uts audio compress podcast.wav -q low
  uts audio compress voice-memo.m4a -q high
  uts audio compress voice.wav -q 256 --dry-run
  uts audio compress ./recordings -r`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := inputs(args, format.AudioExts)
		if err != nil {
			return err
		}

		log.Debug("compressing audio files", "files", files)

		return compress.Audio(compress.AudioOptions{
			Files:     files,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

func init() {
	audioCmd.AddCommand(audioCompressCmd)
	audioCmd.AddCommand(newAudioConvertCmd("convert", []string{"x"}))
}
