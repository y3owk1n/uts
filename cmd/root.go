package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	quality   string
	outputDir string
	inPlace   bool
	dryRun    bool
	verbose   bool
	recursive bool
	algorithm string
	targetFmt string

	Version = "1.0.0"
)

var RootCmd = &cobra.Command{
	Use:   "uts",
	Short: "One CLI for every format",
	Long: `uts — One CLI for every format v` + Version + `

Compress, convert, and inspect any media file without remembering
a dozen different command-line tools.

Quality presets: low, medium, high, or a numeric value
(CRF 0–51 for video, 1–100 for images, 96k–320k for audio, 72–300 DPI for PDF).

Files are saved as <name>-small.<ext> by default. Use -i to replace in-place.`,
	Example: `  uts image compress screenshot.png -q low
  uts video compress recording.mp4 -i
  uts convert image photo.heic --to jpg
  uts info video.mp4`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&quality, "quality", "q", "medium",
		"Quality preset: low, medium, high, or a number (CRF/quality/bitrate/DPI)")
	RootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "",
		"Output directory (default: same as input)")
	RootCmd.PersistentFlags().BoolVarP(&inPlace, "in-place", "i", false,
		"Replace original file with compressed version")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false,
		"Show what would be done without doing it")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose output")
	RootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false,
		"Enable recursive glob patterns")
	RootCmd.PersistentFlags().StringVar(&algorithm, "algorithm", "auto",
		"Archive algorithm (auto, gzip, zstd, xz, brotli, zip)")
	RootCmd.PersistentFlags().StringVar(&targetFmt, "to", "",
		"Target format for conversion")

	RootCmd.AddCommand(videoCmd)
	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(pdfCmd)
	RootCmd.AddCommand(audioCmd)
	RootCmd.AddCommand(archiveCmd)
	RootCmd.AddCommand(convertCmd)
	RootCmd.AddCommand(infoCmd)
}
