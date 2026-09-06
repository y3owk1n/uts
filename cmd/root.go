package cmd

import (
	"fmt"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
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

	// Version is the current version of uts, set at build time.
	Version = "dev"
	// GitCommit is the current commit hash of uts.
	GitCommit = "unknown"
	// BuildDate is the date of the current build.
	BuildDate = "unknown"
)

// RootCmd is the root command for uts.
var RootCmd = &cobra.Command{
	Use:   "uts",
	Short: "One CLI for every format",
	Long: `uts — One CLI for every format v` + Version + `

Compress, convert, and inspect any media file without remembering
a dozen different command-line tools.

Inputs may be files, directories, or quoted glob patterns. Directories are
expanded to the files they contain (the whole tree with -r).

Quality presets: low, medium, high, or a numeric value
(CRF 0–51 for video, 1–100 for images, 96k–320k for audio, 72–300 DPI for PDF).

Files are saved as <name>-small.<ext> by default. Use -i to replace in-place.`,
	Example: `  uts image compress screenshot.png -q low
  uts image compress ./photos -r
  uts video compress recording.mp4 -i
  uts convert image photo.heic --to jpg
  uts info video.mp4`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and returns any error.
func Execute() error {
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if verbose {
			log.SetLevel(log.DebugLevel)
		}

		if inPlace && outputDir != "" {
			ui.Message.Warnf("--in-place is ignored when --output is set; originals are kept")

			inPlace = false
		}

		return nil
	}

	return RootCmd.Execute()
}

// inputs expands CLI arguments into the files a command should process.
// Globs are resolved, directories are listed (recursively with -r) and
// filtered to the given extensions (nil keeps everything).
func inputs(args []string, exts []string) ([]string, error) {
	files := util.ExpandInputs(args, recursive, exts)
	if len(files) == 0 {
		return nil, derrors.New(derrors.CodeFileNotFound, "no input files matched")
	}

	log.Debug("resolved inputs", "count", len(files))

	return files, nil
}

// requireTarget validates the --to flag for convert commands.
func requireTarget(valid []string) error {
	if targetFmt != "" {
		return nil
	}

	return derrors.Newf(derrors.CodeInvalidInput, "missing --to <format>. Valid targets: %v", valid)
}

// addTargetFlag attaches the --to flag with shell completion for its values.
func addTargetFlag(cmd *cobra.Command, valid []string) {
	cmd.Flags().StringVar(&targetFmt, "to", "", "Target format: "+fmt.Sprint(valid))
	_ = cmd.RegisterFlagCompletionFunc(
		"to",
		cobra.FixedCompletions(valid, cobra.ShellCompDirectiveNoFileComp),
	)
}

func init() {
	RootCmd.SetVersionTemplate(
		fmt.Sprintf(
			"uts version %s\nGit commit: %s\nBuild date: %s\n",
			Version,
			GitCommit,
			BuildDate,
		),
	)

	RootCmd.PersistentFlags().StringVarP(&quality, "quality", "q", "medium",
		"Quality preset: low, medium, high, or a number (CRF/quality/bitrate/DPI)")
	RootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "",
		"Output directory (default: same as input)")
	RootCmd.PersistentFlags().BoolVarP(&inPlace, "in-place", "i", false,
		"Replace original file with the result")
	RootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false,
		"Show the commands that would run without running them")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Verbose output (logs every external command)")
	RootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false,
		"Recurse into directories and expand '**' glob patterns")

	RootCmd.AddCommand(videoCmd)
	RootCmd.AddCommand(imageCmd)
	RootCmd.AddCommand(pdfCmd)
	RootCmd.AddCommand(audioCmd)
	RootCmd.AddCommand(archiveCmd)
	RootCmd.AddCommand(convertCmd)
	RootCmd.AddCommand(infoCmd)
	RootCmd.AddCommand(doctorCmd)
	RootCmd.AddCommand(envCmd)
}
