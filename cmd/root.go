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
	Short: "All-in-one utility toolkit",
	Long: `uts — All-in-one utility toolkit v` + Version + `

A modular CLI tool with category-based subcommands for compressing,
converting, and managing media files.

USAGE
  uts <category> <action> <input...> [options]

CATEGORIES
  video     Video files (mp4, mov, mkv, avi, webm)
  image     Images (png, jpg, webp, gif, bmp, tiff, heic, avif)
  pdf       PDF documents
  audio     Audio files (wav, flac, aac, mp3, m4a, opus)
  archive   Directories/files into archives

ACTIONS
  compress  Compress files (available for all categories)
  convert   Convert between formats (image, video, audio, pdf)
  extract   Extract archives (archive only)
  list      List archive contents (archive only)

TOP-LEVEL COMMANDS
  info      Show file info and suggestions
  convert   Convert between formats directly

QUALITY
  high      Best quality, larger files      (crf=23, 192k audio, 300dpi PDF)
  medium    Balanced quality and size       (crf=28, 128k audio, 150dpi PDF)
  low       Smallest files, lower quality   (crf=32, 96k audio, 72dpi PDF)
  <number>  Numeric value (CRF, quality %, kbps, or DPI)

QUICK EXAMPLES
  uts image compress screenshot.png -q low
  uts video compress recording.mp4 -i
  uts convert image photo.heic --to jpg
  uts convert image screenshot.png --to webp -q 85
  uts info video.mp4

OUTPUT
  Files are saved as <name>-small.<ext> in the same directory by default.
  Use -o to specify a different output directory.
  Use -i to replace the original file in-place.`,
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
