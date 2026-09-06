package cmd

import (
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/info"
)

var infoJSON bool

// infoCmd represents the info command.
var infoCmd = &cobra.Command{
	Use:   "info <input...>",
	Short: "Show file info and suggestions",
	Long: `Show file info and suggestions for compression/conversion.

Displays size, type and category. For video, audio and images ffprobe adds
duration, resolution, codecs and bit rate; PDFs show their page count.
Suggests the best compress and convert command for the detected format.

--json prints one JSON array with the same details for scripting.`,
	Example: `  uts info video.mp4
  uts info '*.png'
  uts info ./downloads -r
  uts info clip.mp4 --json | jq '.[0].durationSeconds'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		files, err := inputs(args, nil)
		if err != nil {
			return err
		}

		return info.Show(info.Options{
			Files:   files,
			Version: Version,
			JSON:    infoJSON,
		})
	},
}

func init() {
	infoCmd.Flags().BoolVar(&infoJSON, "json", false, "Print a JSON array instead of panels")
}
