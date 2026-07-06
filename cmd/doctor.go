package cmd

import (
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/doctor"
)

// doctorCmd represents the doctor command.
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check that required tools are installed",
	Long: `Check that all external tools required by uts are available.

Scans your PATH for ffmpeg, ghostscript, pngquant, and other
tools used by uts, then reports which are found and which are missing.

EXAMPLES
  uts doctor`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		doctor.Run(Version)
	},
}
