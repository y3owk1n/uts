package cmd

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/archive"
	"github.com/y3owk1n/uts/internal/compress"
)

var archiveCmd = &cobra.Command{
	Use:     "archive",
	Aliases: []string{"arc", "ar"},
	Short:   "Compress, extract, and list archives",
	Long: `Compress, extract, and list archive files.

ACTIONS
  compress  Create compressed archive (auto-selects best algorithm)
  extract   Extract archive contents (zip, tar, gz, zst, xz, bz2)
  list      List archive contents without extracting

ALGORITHMS
  auto    Auto-select best algorithm (default)
  gzip    gzip compression
  zstd    Zstandard compression (fast + good ratio)
  xz      LZMA2 compression (best ratio)
  brotli  Brotli compression
  zip     ZIP archive (widely compatible)

COMPRESSION EXAMPLES
  uts archive compress ./project/ --algorithm zstd
  uts archive compress ./data/ --algorithm gzip
  uts archive compress ./src/ --dry-run
  uts archive compress ./docs/ --algorithm brotli
  uts archive compress ./photos/ --algorithm zip

EXTRACTION EXAMPLES
  uts archive extract backup.zip
  uts archive extract project.tar.gz
  uts archive extract '*.tar.zst'
  uts archive extract backup.zip -o ./output/

LIST EXAMPLES
  uts archive list backup.zip
  uts archive list project.tar.gz
  uts archive list '*.tar.zst'`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var archiveCompressCmd = &cobra.Command{
	Use:     "compress",
	Aliases: []string{"c"},
	Short:   "Create compressed archives from files/directories",
	Long: `Create compressed archives with the specified algorithm.

USAGE
  uts archive compress <input...> [options]

ALGORITHMS
  auto    Auto-select best algorithm (default)
  gzip    gzip compression
  zstd    Zstandard (fast + good ratio)
  xz      LZMA2 (best ratio)
  brotli  Brotli compression
  zip     ZIP archive (widely compatible)

OUTPUT
  Creates <name>.tar.<algorithm> or <name>.zip in the output directory.

EXAMPLES
  uts archive compress ./project/ --algorithm zstd
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

var archiveExtractCmd = &cobra.Command{
	Use:     "extract",
	Aliases: []string{"x"},
	Short:   "Extract archive contents",
	Long: `Extract archive files to the specified directory.

USAGE
  uts archive extract <archive...> [options]

SUPPORTED FORMATS
  zip     ZIP archives
  tar     Plain tar archives
  gz/tgz  gzip-compressed tar
  zst     Zstandard-compressed tar
  xz/txz  XZ-compressed tar
  bz2     bzip2-compressed tar

OUTPUT
  Extracts archive contents into the output directory.

EXAMPLES
  uts archive extract backup.zip
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

var archiveListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List archive contents",
	Long: `List the contents of archive files without extracting.

USAGE
  uts archive list <archive...> [options]

SUPPORTED FORMATS
  zip     ZIP archives
  tar     Plain tar archives
  gz/tgz  gzip-compressed tar
  zst     Zstandard-compressed tar
  xz/txz  XZ-compressed tar
  bz2     bzip2-compressed tar

EXAMPLES
  uts archive list backup.zip
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
