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
	Short:   "Video files (mp4, mov, mkv, avi, webm)",
	Long: `Compress and convert video files using ffmpeg.
Supported formats: mp4, mov, mkv, avi, webm, m4v, flv, wmv`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var videoCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Compress video files using ffmpeg (libx265)",
	Args:    cobra.MinimumNArgs(1),
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
	Long:    `Convert video files between formats. Requires --to flag.`,
	Args:    cobra.MinimumNArgs(1),
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
