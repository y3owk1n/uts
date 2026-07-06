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
	Short:   "Compress and convert audio files",
	Long: `Compress and convert audio files using ffmpeg.

ACTIONS
  compress  Compress audio using ffmpeg (aac)
  convert   Convert between audio formats (mp3, wav, flac, opus, ...)

SUPPORTED FORMATS
  Input:   wav, flac, aac, mp3, m4a, opus, ogg, wma
  Output:  mp3, aac, m4a, wav, flac, opus, ogg

COMPRESSION EXAMPLES
  uts audio compress podcast.wav -q low
  uts audio compress voice-memo.m4a -q high
  uts audio compress '*.wav' -r
  uts audio compress track.flac --dry-run

CONVERSION EXAMPLES
  uts audio convert track.wav --to mp3
  uts audio convert song.flac --to m4a -q high
  uts audio convert '*.wav' --to mp3 -q 96
  uts audio convert podcast.opus --to mp3 -q 256`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var audioCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress audio files using ffmpeg (aac)",
	Long: `Compress audio files using ffmpeg with AAC codec.

USAGE
  uts audio compress <input...> [options]

SUPPORTED FORMATS
  wav, flac, aac, mp3, m4a, opus, ogg, wma

QUALITY
  high       192k aac
  medium     128k aac (default)
  low        96k aac
  <kbps>     Numeric bitrate (e.g. 256, 320)

OUTPUT
  Files saved as <name>-small.m4a in the same directory.

EXAMPLES
  uts audio compress podcast.wav -q low
  uts audio compress voice-memo.m4a -q high
  uts audio compress track1.wav track2.flac track3.m4a -q medium
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

var audioConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between audio formats",
	Long: `Convert audio files between formats using ffmpeg.

USAGE
  uts audio convert <input...> --to <format> [options]

TARGET FORMATS
  mp3   libmp3lame
  aac   aac
  m4a   aac
  wav   pcm_s16le
  flac  flac
  opus  libopus
  ogg   libvorbis

QUALITY (bitrate)
  high       192k
  medium     128k (default)
  low        96k
  <kbps>     Numeric bitrate (e.g. 256, 320)

EXAMPLES
  uts audio convert track.wav --to mp3
  uts audio convert song.flac --to m4a -q high
  uts audio convert track1.wav track2.flac --to mp3
  uts audio convert '*.wav' --to mp3 -q 96
  uts audio convert lecture.wav --to mp3 -q 256`,
	Args: cobra.MinimumNArgs(1),
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
