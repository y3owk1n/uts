package cmd

import (
	"github.com/spf13/cobra"
	"github.com/y3owk1n/uts/internal/env"
)

// envCmd represents the env command.
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show recognized environment variables",
	Long: `Show all environment variables recognized by uts and their current values.

Displays every UTS_COLOR_* variable, NO_COLOR, and FORCE_COLOR together
with their defaults and a short description.

EXAMPLES
  uts env`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		env.Run(Version)

		return nil
	},
}
