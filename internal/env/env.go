// Package env displays all environment variables recognized by uts.
package env

import (
	"fmt"
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/style"
)

const (
	sectionControl = "Control"
	sectionTheming = "Theming"
)

// varEntry describes one row in the env table.
type varEntry struct {
	Section string
	Name    string
	Value   string
}

// Run displays all recognized environment variables and their current values.
func Run(version string) {
	palette := ui.Style.Palette()
	vars := collectVars(palette)

	tbl := ui.Table.New("Section", "Variable", "Value")
	for _, v := range vars {
		tbl.Row(v.Section, v.Name, ui.Message.Accent(v.Value))
	}

	_, _ = lipgloss.Fprint(os.Stdout, ui.Banner.Logo(version))
	_, _ = lipgloss.Fprintln(os.Stdout)
	_, _ = lipgloss.Fprint(os.Stdout, tbl.Render(palette))
}

// collectVars resolves every env var uts cares about.
func collectVars(palette style.Palette) []varEntry {
	forceColor := envValue("FORCE_COLOR", "not set")
	noColor := envValue("NO_COLOR", "not set")

	return []varEntry{
		// Color control
		{Section: sectionControl, Name: "FORCE_COLOR", Value: forceColor},
		{Section: sectionControl, Name: "NO_COLOR", Value: noColor},
		// Theming — resolved palette values
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_PRIMARY",
			Value:   adaptiveColorValue(palette.Primary),
		},
		{Section: sectionTheming, Name: "UTS_COLOR_TEXT", Value: adaptiveColorValue(palette.Text)},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_MUTED",
			Value:   adaptiveColorValue(palette.Muted),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_SUBTLE",
			Value:   adaptiveColorValue(palette.Subtle),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_BORDER",
			Value:   adaptiveColorValue(palette.Border),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_ACCENT",
			Value:   adaptiveColorValue(palette.Accent),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_SUCCESS",
			Value:   adaptiveColorValue(palette.Success),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_WARNING",
			Value:   adaptiveColorValue(palette.Warning),
		},
		{
			Section: sectionTheming,
			Name:    "UTS_COLOR_ERROR",
			Value:   adaptiveColorValue(palette.Error),
		},
	}
}

// envValue returns the value of an env var, or the fallback if unset.
func envValue(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}

	return fallback
}

// adaptiveColorValue formats an AdaptiveColor as
// "Light: <hex>, Dark: <hex>" so the user can see at a glance
// what each terminal background will render.
func adaptiveColorValue(c compat.AdaptiveColor) string {
	return "Light: " + colorToHex(c.Light) + ", Dark: " + colorToHex(c.Dark)
}

// colorToHex converts a color.Color to a "#rrggbb" hex string.
func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()

	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8) //nolint:mnd
}
