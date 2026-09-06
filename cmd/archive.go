// Package cmd contains the uts CLI.
//
//nolint:goconst
package cmd

import (
	"strings"

	"charm.land/log/v2"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/archive"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/util"
)

// archiveCmd represents the archive command.
var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"arc", "ar"},
	Short:   "Compress, extract, and list archives",
	Long: `Compress, extract, and list archives.

Algorithms: ` + strings.Join(archive.Algorithms, ", ") + `
Archive formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2, tar.br`,
	Example: `  uts archive compress ./project/ --algorithm zstd
  uts archive extract backup.zip
  uts archive list project.tar.gz`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// archiveCompressCmd represents the archive compress command.
var archiveCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Create a compressed archive from files/directories",
	Long: `Create a compressed archive with the chosen algorithm.

Algorithms: zip (default), gzip, zstd, xz, brotli.
Output is named after the input: <name>.tar.<algo> or <name>.zip.`,
	Example: `  uts archive compress ./project/ --algorithm zstd
  uts archive compress ./data/ --algorithm zip
  uts archive compress notes.md photo.jpg -o ./backups/
  uts archive compress ./src/ --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("creating archive", "files", args, "algorithm", algorithm)

		return archive.Create(archive.CreateOptions{
			Files:     args,
			Algorithm: algorithm,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

// archiveExtractCmd represents the archive extract command.
var archiveExtractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"x"},
	Short:   "Extract archive contents",
	Long: `Extract archives into the output directory (default: current directory).

Supported formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2, tar.br`,
	Example: `  uts archive extract backup.zip
  uts archive extract project.tar.gz -o ./project/
  uts archive extract '*.tar.zst' -o ./output/
  uts archive extract backup.zip --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := archiveInputs(args)
		if err != nil {
			return err
		}

		log.Debug("extracting archives", "files", files)

		return archive.Extract(archive.ExtractOptions{
			Files:     files,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

// archiveListCmd represents the archive list command.
var archiveListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List archive contents",
	Long: `List the contents of archives without extracting.

Supported formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2, tar.br`,
	Example: `  uts archive list backup.zip
  uts archive list project.tar.gz
  uts archive list '*.tar.zst'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := archiveInputs(args)
		if err != nil {
			return err
		}

		log.Debug("listing archives", "files", files)

		return archive.List(archive.ListOptions{Files: files})
	},
}

// archiveInputs resolves globs but never expands directories: an archive
// command's inputs are the archives themselves.
func archiveInputs(args []string) ([]string, error) {
	files := util.ResolveGlobs(args, recursive)
	if len(files) == 0 {
		return nil, derrors.New(derrors.CodeFileNotFound, "no input files matched")
	}

	return files, nil
}

func init() {
	archiveCompressCmd.Flags().StringVar(&algorithm, "algorithm", "zip",
		"Archive algorithm: "+strings.Join(archive.Algorithms, ", "))
	_ = archiveCompressCmd.RegisterFlagCompletionFunc("algorithm",
		cobra.FixedCompletions(archive.Algorithms, cobra.ShellCompDirectiveNoFileComp))

	archiveCmd.AddCommand(archiveCompressCmd)
	archiveCmd.AddCommand(archiveExtractCmd)
	archiveCmd.AddCommand(archiveListCmd)
}
