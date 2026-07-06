package cmd

import (
	"fmt"
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

Input formats: mp4, mov, mkv, avi, webm, m4v, flv, wmv
Output formats: mp4, mov, mkv, webm, avi, flv`,
	Example: `  uts video compress screen-recording.mp4 -q low
  uts video convert clip.mov --to mp4`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var videoCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress video files using ffmpeg (libx265)",
	Long: `Compress video files using ffmpeg with libx265 codec.

Quality: high (crf=23, slow), medium (crf=28, medium), low (crf=32, fast), or raw 0-51.`,
	Example: `  uts video compress screen-recording.mp4 -q low
  uts video compress vacation.mov -q high -i
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

Target formats: mp4, mov, mkv, webm, avi, flv`,
	Example: `  uts video convert clip.mov --to mp4
  uts video convert recording.mkv --to webm -q medium
  uts video convert presentation.avi --to mkv -q 18
  uts video convert '*.mov' --to mp4`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if targetFmt == "" {
			return fmt.Errorf("missing --to <format>. Examples: --to mp4, --to mkv, --to webm")
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
