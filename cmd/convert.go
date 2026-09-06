package cmd

import (
	"slices"
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/convert"
	"github.com/y3owk1n/uts/internal/format"
)

// convertCmd is the top-level alias: "uts convert image ..." behaves exactly
// like "uts image convert ...".
var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"x"},
	Short:   "Convert between formats",
	Long: `Convert files between different formats (image, video, audio, pdf).

Requires --to <format> to specify the target format.`,
	Example: `  uts convert image photo.heic --to jpg
  uts convert video clip.mov --to mp4
  uts convert audio track.wav --to mp3 -q 96
  uts convert pdf report.pdf --to jpg`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func newImageConvertCmd(use string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "Convert between image formats",
		Long: `Convert image files between formats using ImageMagick (or sips on macOS).

Target formats: ` + strings.Join(format.ImageTargets, ", "),
		Example: `  uts image convert photo.heic --to jpg
  uts image convert screenshot.png --to webp -q high
  uts image convert photo.jpg --to avif -q 70
  uts image convert photo.heic --to jpg --max 1600
  uts image convert ./photos --to jpg -r`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := requireTarget(format.ImageTargets)
			if err != nil {
				return err
			}

			err = checkMaxEdge()
			if err != nil {
				return err
			}

			files, err := inputs(args, format.ImageExts)
			if err != nil {
				return err
			}

			log.Debug("converting image files", "files", files, "target", targetFmt)

			return convert.Image(convert.ImageOptions{
				Files:     files,
				Target:    targetFmt,
				Quality:   quality,
				OutputDir: outputDir,
				InPlace:   inPlace,
				DryRun:    dryRun,
				MaxEdge:   maxEdge,
			})
		},
	}
	addTargetFlag(cmd, format.ImageTargets)
	addMaxFlag(cmd)

	return cmd
}

func newVideoConvertCmd(use string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "Convert between video formats",
		Long: `Convert video files between containers using ffmpeg.

When the source codecs already fit the target container the streams are
copied without re-encoding: instant and lossless. Pass -q or --max to
force a re-encode.

Target formats: ` + strings.Join(format.VideoTargets, ", "),
		Example: `  uts video convert clip.mov --to mp4
  uts video convert recording.mkv --to webm -q medium
  uts video convert presentation.avi --to mkv -q 18
  uts video convert clip.mov --to mp4 --max 1280
  uts video convert ./clips --to mp4`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := requireTarget(format.VideoTargets)
			if err != nil {
				return err
			}

			err = checkMaxEdge()
			if err != nil {
				return err
			}

			files, err := inputs(args, format.VideoExts)
			if err != nil {
				return err
			}

			log.Debug("converting video files", "files", files, "target", targetFmt)

			return convert.Video(convert.VideoOptions{
				Files:      files,
				Target:     targetFmt,
				Quality:    quality,
				QualitySet: cmd.Flags().Changed("quality"),
				OutputDir:  outputDir,
				InPlace:    inPlace,
				DryRun:     dryRun,
				MaxEdge:    maxEdge,
			})
		},
	}
	addTargetFlag(cmd, format.VideoTargets)
	addMaxFlag(cmd)

	return cmd
}

func newAudioConvertCmd(use string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "Convert between audio formats",
		Long: `Convert audio files, or extract the audio track of video files, using ffmpeg.

Target formats: ` + strings.Join(format.AudioTargets, ", "),
		Example: `  uts audio convert track.wav --to mp3
  uts audio convert video.mp4 --to mp3
  uts audio convert song.flac --to m4a -q high
  uts audio convert ./music --to opus -q 96`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := requireTarget(format.AudioTargets)
			if err != nil {
				return err
			}

			files, err := inputs(args, slices.Concat(format.AudioExts, format.VideoExts))
			if err != nil {
				return err
			}

			log.Debug("converting audio files", "files", files, "target", targetFmt)

			return convert.Audio(convert.AudioOptions{
				Files:     files,
				Target:    targetFmt,
				Quality:   quality,
				OutputDir: outputDir,
				InPlace:   inPlace,
				DryRun:    dryRun,
			})
		},
	}
	addTargetFlag(cmd, format.AudioTargets)

	return cmd
}

func newPDFConvertCmd(use string, aliases []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Aliases: aliases,
		Short:   "Convert between PDF and images",
		Long: `Convert PDF documents to images or combine images into a PDF.

PDF → images: writes <basename>/page-*.<ext> next to the PDF (jpg or png)
images → PDF: combines the images, in order, into a single PDF`,
		Example: `  uts pdf convert report.pdf --to jpg
  uts pdf convert slides.pdf --to png -q high
  uts pdf convert page1.png page2.png --to pdf
  uts pdf convert './scans/*.jpg' --to pdf`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := requireTarget(format.PDFTargets)
			if err != nil {
				return err
			}

			files, err := inputs(args, slices.Concat(format.PDFExts, format.ImageExts))
			if err != nil {
				return err
			}

			log.Debug("converting PDF files", "files", files, "target", targetFmt)

			return convert.PDF(convert.PDFOptions{
				Files:     files,
				Target:    targetFmt,
				Quality:   quality,
				OutputDir: outputDir,
				DryRun:    dryRun,
			})
		},
	}
	addTargetFlag(cmd, format.PDFTargets)

	return cmd
}

func init() {
	convertCmd.AddCommand(newImageConvertCmd("image", []string{"img", "i"}))
	convertCmd.AddCommand(newVideoConvertCmd("video", []string{"v"}))
	convertCmd.AddCommand(newAudioConvertCmd("audio", []string{"a"}))
	convertCmd.AddCommand(newPDFConvertCmd("pdf", []string{"p"}))
}
