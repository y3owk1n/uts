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

var rootCmd = &cobra.Command{
	Use:   "uts",
	Short: "All-in-one utility toolkit",
	Long: `uts is an all-in-one utility toolkit for compressing and converting
media files including video, image, audio, PDF, and archives.

Usage: uts <category> <action> <input...> [options]`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&quality, "quality", "q", "medium",
		"Quality preset: low, medium, high, or a number")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "",
		"Output directory (default: same as input)")
	rootCmd.PersistentFlags().BoolVarP(&inPlace, "in-place", "i", false,
		"Replace original file with compressed version")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false,
		"Show what would be done without doing it")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false,
		"Enable recursive glob patterns")
	rootCmd.PersistentFlags().StringVar(&algorithm, "algorithm", "auto",
		"Archive algorithm (auto, gzip, zstd, xz, brotli, zip)")
	rootCmd.PersistentFlags().StringVar(&targetFmt, "to", "",
		"Target format for conversion")

	rootCmd.AddCommand(videoCmd)
	rootCmd.AddCommand(imageCmd)
	rootCmd.AddCommand(pdfCmd)
	rootCmd.AddCommand(audioCmd)
	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(infoCmd)
}
