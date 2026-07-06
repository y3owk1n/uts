// Package cmd contains the uts CLI.
//
//nolint:goconst
package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/archive"
	"github.com/y3owk1n/uts/internal/compress"
)

// archiveCmd represents the archive command.
var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"arc", "ar"},
	Short:   "Compress, extract, and list archives",
	Long: `Compress, extract, and list archives.

Supported algorithms: gzip, zstd, xz, brotli, zip
Archive formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2`,
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
	Short:   "Create compressed archives from files/directories",
	Long: `Create compressed archives with the specified algorithm.

Algorithms: zip (default), gzip, zstd, xz, brotli.
Output saved as <name>.tar.<algo> or <name>.zip.`,
	Example: `  uts archive compress ./project/ --algorithm zstd
  uts archive compress ./data/ --algorithm zip
  uts archive compress ./src/ --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("compressing archives", "files", args, "algorithm", algorithm)

		return compress.Archive(compress.ArchiveOptions{
			Files:     args,
			Algorithm: strings.ToLower(algorithm),
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
	Long: `Extract archive files to the specified directory.

Supported formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2`,
	Example: `  uts archive extract backup.zip
  uts archive extract project.tar.gz
  uts archive extract '*.tar.zst' -o ./output/
  uts archive extract backup.zip --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("extracting archives", "files", args)

		return archive.Extract(archive.ExtractOptions{
			Files:     args,
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
	Long: `List the contents of archive files without extracting.

Supported formats: zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2`,
	Example: `  uts archive list backup.zip
  uts archive list project.tar.gz
  uts archive list '*.tar.zst'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("listing archives", "files", args)

		return archive.List(archive.ListOptions{
			Files: args,
		})
	},
}

func init() {
	archiveCmd.AddCommand(archiveCompressCmd)
	archiveCmd.AddCommand(archiveExtractCmd)
	archiveCmd.AddCommand(archiveListCmd)
}
