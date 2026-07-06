package cmd

import (
	"github.com/spf13/cobra"

	"github.com/y3owk1n/uts/internal/info"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show file info and suggestions",
	Long: `Display file size, type, and suggest the best compress/convert
command for the detected format.

Examples:
  uts info video.mp4
  uts info *.png
  uts info photo.heic`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		info.Show(info.Options{
			Files: args,
		})
		return nil
	},
}
