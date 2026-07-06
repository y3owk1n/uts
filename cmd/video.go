package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/convert"
)

var videoCmd = &cobra.Command{
	Use:     "video",
	Aliases: []string{"v"},
	Short:   "Compress and convert video files",
	Long: `Compress and convert video files using ffmpeg (libx265).

ACTIONS
  compress  Compress videos using ffmpeg (libx265)
  convert   Convert between video formats (mp4, mkv, webm, mov, avi, flv)

SUPPORTED FORMATS
  Input:   mp4, mov, mkv, avi, webm, m4v, flv, wmv
  Output:  mp4, mov, mkv, webm, avi, flv

COMPRESSION EXAMPLES
  uts video compress screen-recording.mp4 -q low
  uts video compress vacation.mov -q high -i
  uts video compress lecture.mkv --dry-run
  uts video compress '*.mp4' -r -q medium

CONVERSION EXAMPLES
  uts video convert clip.mov --to mp4
  uts video convert recording.mkv --to webm -q medium
  uts video convert presentation.avi --to mp4 -q 18
  uts video convert '*.mov' --to mp4 -i`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var videoCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress video files using ffmpeg (libx265)",
	Long: `Compress video files using ffmpeg with libx265 codec.

USAGE
  uts video compress <input...> [options]

SUPPORTED FORMATS
  mp4, mov, mkv, avi, webm, m4v, flv, wmv

QUALITY
  high       crf=23, preset=slow
  medium     crf=28, preset=medium (default)
  low        crf=32, preset=fast
  <0-51>     Raw CRF value (lower = better quality)

OUTPUT
  Files saved as <name>-small.<ext> in the same directory.

EXAMPLES
  uts video compress screen-recording.mp4 -q low
  uts video compress vacation.mov -q high -i
  uts video compress clip1.mp4 clip2.mp4 clip3.mp4 -q medium
  uts video compress lecture.mkv -q 25 --dry-run -v
  uts video compress '*.mp4' -r`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("compressing video files", "files", args)
		return compress.Video(compress.VideoOptions{
			Files:     args,
			Quality:   quality,
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

var videoConvertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between video formats",
	Long: `Convert video files between formats using ffmpeg.

USAGE
  uts video convert <input...> --to <format> [options]

TARGET FORMATS
  mp4   libx264 / aac
  mov   libx264 / aac
  mkv   libx265 / aac
  webm  libvpx-vp9 / libopus
  avi   libx264 / mp3
  flv   libx264 / aac

EXAMPLES
  uts video convert clip.mov --to mp4
  uts video convert recording.mkv --to webm -q medium
  uts video convert presentation.avi --to mkv -q 18
  uts video convert clip1.mov clip2.mov --to mp4 -i
  uts video convert '*.mov' --to mp4`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			log.Error("Missing --to <format>. Examples: --to mp4, --to mkv, --to webm")
			return nil
		}
		log.Debug("converting video files", "files", args, "target", targetFmt)
		return convert.Video(convert.VideoOptions{
			Files:     args,
			Target:    strings.ToLower(targetFmt),
			OutputDir: outputDir,
			InPlace:   inPlace,
			DryRun:    dryRun,
		})
	},
}

func init() {
	videoCmd.AddCommand(videoCompressCmd)
	videoCmd.AddCommand(videoConvertCmd)
}
