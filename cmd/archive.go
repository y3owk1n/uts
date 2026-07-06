package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/compress"
	"github.com/y3owk1n/uts/internal/archive"
)

var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"arc", "ar"},
	Short:   "Archive files and directories",
	Long:    `Compress, extract, and list archive files.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var archiveCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Create compressed archives from files/directories",
	Args:    cobra.MinimumNArgs(1),
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

var archiveExtractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"x"},
	Short:   "Extract archive contents",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debug("extracting archives", "files", args)
		return archive.Extract(archive.ExtractOptions{
			Files:     args,
			OutputDir: outputDir,
			DryRun:    dryRun,
		})
	},
}

var archiveListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List archive contents",
	Args:    cobra.MinimumNArgs(1),
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
