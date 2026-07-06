package style

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Palette defines the color palette for the CLI.
type Palette struct {
	Primary lipgloss.AdaptiveColor
	Text    lipgloss.AdaptiveColor
	Muted   lipgloss.AdaptiveColor
	Subtle  lipgloss.AdaptiveColor
	Border  lipgloss.AdaptiveColor
	Accent  lipgloss.AdaptiveColor
	Success lipgloss.AdaptiveColor
	Warning lipgloss.AdaptiveColor
	Error   lipgloss.AdaptiveColor
}

const (
	base00 = "#1f1d2e"
	base01 = "#2a2738"
	base02 = "#3a364d"
	base03 = "#5a5672"
	base04 = "#9a96b5"
	base05 = "#e0def4"
	base08 = "#f28fad"
	base0A = "#f9e2af"
	base0B = "#abe9b3"
	base0D = "#80b8e8"
	base0E = "#c9a0e9"
)

func basePalette() Palette {
	return Palette{
		Primary: lipgloss.AdaptiveColor{Light: "#6f4d8c", Dark: base0E},
		Text:    lipgloss.AdaptiveColor{Light: base01, Dark: base05},
		Muted:   lipgloss.AdaptiveColor{Light: base03, Dark: base04},
		Subtle:  lipgloss.AdaptiveColor{Light: base04, Dark: base03},
		Border:  lipgloss.AdaptiveColor{Light: base01, Dark: base02},
		Accent:  lipgloss.AdaptiveColor{Light: "#4068a0", Dark: base0D},
		Success: lipgloss.AdaptiveColor{Light: "#5a9b65", Dark: base0B},
		Warning: lipgloss.AdaptiveColor{Light: "#b89556", Dark: base0A},
		Error:   lipgloss.AdaptiveColor{Light: "#b86080", Dark: base08},
	}
}

// Default returns the default palette with optional env overrides.
func Default() Palette {
	palette := basePalette()
	palette.Primary = overrideColor(palette.Primary, "PRIMARY")
	palette.Text = overrideColor(palette.Text, "TEXT")
	palette.Muted = overrideColor(palette.Muted, "MUTED")
	palette.Subtle = overrideColor(palette.Subtle, "SUBTLE")
	palette.Border = overrideColor(palette.Border, "BORDER")
	palette.Accent = overrideColor(palette.Accent, "ACCENT")
	palette.Success = overrideColor(palette.Success, "SUCCESS")
	palette.Warning = overrideColor(palette.Warning, "WARNING")
	palette.Error = overrideColor(palette.Error, "ERROR")

	return palette
}

func overrideColor(color lipgloss.AdaptiveColor, name string) lipgloss.AdaptiveColor {
	if v, ok := envColor("UTS_COLOR_" + name); ok {
		color.Light = v
		color.Dark = v
	}

	if v, ok := envColor("UTS_COLOR_" + name + "_LIGHT"); ok {
		color.Light = v
	}

	if v, ok := envColor("UTS_COLOR_" + name + "_DARK"); ok {
		color.Dark = v
	}

	return color
}

func envColor(envName string) (string, bool) {
	raw := strings.TrimSpace(os.Getenv(envName))

	return raw, raw != ""
}

// ColorEnabled reports whether color output is enabled.
func ColorEnabled() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}

	if _, ok := os.LookupEnv("FORCE_COLOR"); ok {
		return true
	}

	return lipgloss.DefaultRenderer().ColorProfile().Name() != "ascii"
}
