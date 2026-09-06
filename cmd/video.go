//nolint:goconst
package cmd

import (
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/format"
)

// videoCmd represents the video command.
var videoCmd = &cobra.Command{
	Use:     "video",
	Aliases: []string{"v"},
	Short:   "Compress and convert video files",
	Long: `Compress and convert video files using ffmpeg.

Input formats: ` + strings.Join(format.VideoExts, ", ") + `
Output formats: ` + strings.Join(format.VideoTargets, ", "),
	Example: `  uts video compress screen-recording.mp4 -q low
  uts video convert clip.mov --to mp4`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// videoCompressCmd represents the video compress command.
var videoCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress video files using ffmpeg",
	Long: `Compress video files using ffmpeg, keeping the container format.

Codecs: mp4/mov/avi/flv use H.264, mkv uses H.265, webm uses VP9.
Quality: high (crf=23, slow), medium (crf=28, medium), low (crf=32, fast), or raw 0-51.`,
	Example: `  uts video compress screen-recording.mp4 -q low
  uts video compress vacation.mov -q high -i
  uts video compress lecture.mkv -q 25 --dry-run
  uts video compress 4k-screen-recording.mov --max 1920
  uts video compress ./recordings -r`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := checkMaxEdge()
		if err != nil {
			return err
		}

		files, err := inputs(args, format.VideoExts)
		if err != nil {
			return err
		}

		log.Debug("compressing video files", "files", files)

		return compress.Video(compress.VideoOptions{
			Files:        files,
			Quality:      quality,
			OutputDir:    outputDir,
			InPlace:      inPlace,
			DryRun:       dryRun,
			Jobs:         jobs,
			SkipExisting: skipExisting,
			Backup:       backup,
			MaxEdge:      maxEdge,
		})
	},
}

func init() {
	addMaxFlag(videoCompressCmd)
	videoCmd.AddCommand(videoCompressCmd)
	videoCmd.AddCommand(newVideoConvertCmd("convert", []string{"x"}))
}
