// Package env displays all environment variables recognized by uts.
package env

import (
	"fmt"
	"os"
	"sort"

	"charm.land/lipgloss/v2"
	"github.com/y3owk1n/uts/internal/ui"
)

// varEntry describes one environment variable.
type varEntry struct {
	Name        string
	Default     string
	Description string
}

// vars returns the full list of recognized variables.
func vars() []varEntry {
	return []varEntry{
		// Color palette
		{
			Name:        "UTS_COLOR_PRIMARY",
			Default:     "#6f4d8c / #c9a0e9",
			Description: "Logo, titles, active elements",
		},
		{
			Name:        "UTS_COLOR_PRIMARY_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of PRIMARY",
		},
		{
			Name:        "UTS_COLOR_PRIMARY_DARK",
			Default:     "",
			Description: "Dark-terminal variant of PRIMARY",
		},
		{Name: "UTS_COLOR_TEXT", Default: "#2a2738 / #e0def4", Description: "Body text"},
		{Name: "UTS_COLOR_TEXT_LIGHT", Default: "", Description: "Light-terminal variant of TEXT"},
		{Name: "UTS_COLOR_TEXT_DARK", Default: "", Description: "Dark-terminal variant of TEXT"},
		{
			Name:        "UTS_COLOR_MUTED",
			Default:     "#5a5672 / #9a96b5",
			Description: "Dimmed / secondary text",
		},
		{
			Name:        "UTS_COLOR_MUTED_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of MUTED",
		},
		{Name: "UTS_COLOR_MUTED_DARK", Default: "", Description: "Dark-terminal variant of MUTED"},
		{
			Name:        "UTS_COLOR_SUBTLE",
			Default:     "#9a96b5 / #5a5672",
			Description: "Hints and subtle labels",
		},
		{
			Name:        "UTS_COLOR_SUBTLE_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of SUBTLE",
		},
		{
			Name:        "UTS_COLOR_SUBTLE_DARK",
			Default:     "",
			Description: "Dark-terminal variant of SUBTLE",
		},
		{
			Name:        "UTS_COLOR_BORDER",
			Default:     "#2a2738 / #3a364d",
			Description: "Panel borders and separators",
		},
		{
			Name:        "UTS_COLOR_BORDER_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of BORDER",
		},
		{
			Name:        "UTS_COLOR_BORDER_DARK",
			Default:     "",
			Description: "Dark-terminal variant of BORDER",
		},
		{
			Name:        "UTS_COLOR_ACCENT",
			Default:     "#4068a0 / #80b8e8",
			Description: "Info highlights and links",
		},
		{
			Name:        "UTS_COLOR_ACCENT_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of ACCENT",
		},
		{
			Name:        "UTS_COLOR_ACCENT_DARK",
			Default:     "",
			Description: "Dark-terminal variant of ACCENT",
		},
		{Name: "UTS_COLOR_SUCCESS", Default: "#5a9b65 / #abe9b3", Description: "Success messages"},
		{
			Name:        "UTS_COLOR_SUCCESS_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of SUCCESS",
		},
		{
			Name:        "UTS_COLOR_SUCCESS_DARK",
			Default:     "",
			Description: "Dark-terminal variant of SUCCESS",
		},
		{Name: "UTS_COLOR_WARNING", Default: "#b89556 / #f9e2af", Description: "Warning messages"},
		{
			Name:        "UTS_COLOR_WARNING_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of WARNING",
		},
		{
			Name:        "UTS_COLOR_WARNING_DARK",
			Default:     "",
			Description: "Dark-terminal variant of WARNING",
		},
		{Name: "UTS_COLOR_ERROR", Default: "#b86080 / #f28fad", Description: "Error messages"},
		{
			Name:        "UTS_COLOR_ERROR_LIGHT",
			Default:     "",
			Description: "Light-terminal variant of ERROR",
		},
		{Name: "UTS_COLOR_ERROR_DARK", Default: "", Description: "Dark-terminal variant of ERROR"},

		// Color control
		{
			Name:        "NO_COLOR",
			Default:     "",
			Description: "Disable all color output (see no-color.org)",
		},
		{Name: "FORCE_COLOR", Default: "", Description: "Force color output (e.g. when piping)"},
	}
}

// Run displays all recognized environment variables and their current values.
func Run(version string) {
	palette := ui.Style.Palette()
	entries := vars()

	// Sort by name for stable output.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	_, _ = fmt.Fprint(os.Stdout, ui.Banner.Logo(version))

	nameStyle := lipgloss.NewStyle().Foreground(palette.Accent).Bold(true)
	valStyle := lipgloss.NewStyle().Foreground(palette.Text)
	unsetStyle := lipgloss.NewStyle().Foreground(palette.Muted).Italic(true)
	descStyle := lipgloss.NewStyle().Foreground(palette.Subtle)

	for _, entry := range entries {
		val, ok := os.LookupEnv(entry.Name)

		var valuePart string
		switch {
		case ok:
			valuePart = valStyle.Render(val)
		case entry.Default != "":
			valuePart = unsetStyle.Render(entry.Default + " (default)")
		default:
			valuePart = unsetStyle.Render("not set")
		}

		_, _ = fmt.Fprintf(
			os.Stdout, "  %s  %s  %s\n",
			nameStyle.Render(entry.Name),
			valuePart,
			descStyle.Render(entry.Description),
		)
	}
}
