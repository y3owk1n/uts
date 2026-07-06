package cmd

import (
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/info"
)

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show file info and suggestions",
	Long: `Show file info and suggestions for compression/conversion.

USAGE
  uts info <input...> [options]

DESCRIPTION
  Displays file size, type, and suggests the best compress/convert
  command for the detected format.

EXAMPLES
  uts info video.mp4
  uts info '*.png'
  uts info photo.heic`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		info.Show(info.Options{
			Files:   args,
			Version: Version,
		})

		return nil
	},
}
