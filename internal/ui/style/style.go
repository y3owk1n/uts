package style

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
	"github.com/charmbracelet/colorprofile"
)

// init bridges the uts-specific FORCE_COLOR env var to CLICOLOR_FORCE
// so that colorprofile (used internally by lipgloss.Fprint) correctly
// forces color output when piping. Without this, lipgloss strips ANSI
// escapes when stdout is not a TTY, even with FORCE_COLOR=1 set.
func init() {
	if _, ok := os.LookupEnv("FORCE_COLOR"); ok {
		_ = os.Setenv("CLICOLOR_FORCE", "1")
	}
}

// Palette defines the color palette for the CLI.
type Palette struct {
	Primary compat.AdaptiveColor
	Text    compat.AdaptiveColor
	Muted   compat.AdaptiveColor
	Subtle  compat.AdaptiveColor
	Border  compat.AdaptiveColor
	Accent  compat.AdaptiveColor
	Success compat.AdaptiveColor
	Warning compat.AdaptiveColor
	Error   compat.AdaptiveColor
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
		Primary: compat.AdaptiveColor{
			Light: lipgloss.Color("#6f4d8c"),
			Dark:  lipgloss.Color(base0E),
		},
		Text:   compat.AdaptiveColor{Light: lipgloss.Color(base01), Dark: lipgloss.Color(base05)},
		Muted:  compat.AdaptiveColor{Light: lipgloss.Color(base03), Dark: lipgloss.Color(base04)},
		Subtle: compat.AdaptiveColor{Light: lipgloss.Color(base04), Dark: lipgloss.Color(base03)},
		Border: compat.AdaptiveColor{Light: lipgloss.Color(base01), Dark: lipgloss.Color(base02)},
		Accent: compat.AdaptiveColor{
			Light: lipgloss.Color("#4068a0"),
			Dark:  lipgloss.Color(base0D),
		},
		Success: compat.AdaptiveColor{
			Light: lipgloss.Color("#5a9b65"),
			Dark:  lipgloss.Color(base0B),
		},
		Warning: compat.AdaptiveColor{
			Light: lipgloss.Color("#b89556"),
			Dark:  lipgloss.Color(base0A),
		},
		Error: compat.AdaptiveColor{
			Light: lipgloss.Color("#b86080"),
			Dark:  lipgloss.Color(base08),
		},
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

func overrideColor(color compat.AdaptiveColor, name string) compat.AdaptiveColor {
	if v, ok := envColor("UTS_COLOR_" + name); ok {
		color.Light = lipgloss.Color(v)
		color.Dark = lipgloss.Color(v)
	}

	if v, ok := envColor("UTS_COLOR_" + name + "_LIGHT"); ok {
		color.Light = lipgloss.Color(v)
	}

	if v, ok := envColor("UTS_COLOR_" + name + "_DARK"); ok {
		color.Dark = lipgloss.Color(v)
	}

	return color
}

var hexColorPattern = regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var namedColors = map[string]struct{}{
	"black":   {},
	"red":     {},
	"green":   {},
	"yellow":  {},
	"blue":    {},
	"magenta": {},
	"cyan":    {},
	"white":   {},
	"gray":    {},
	"grey":    {},
}

func isValidColor(value string) bool {
	lower := strings.ToLower(value)

	n, err := strconv.Atoi(lower)
	if err == nil {
		return n >= 0 && n <= 255
	}

	if hexColorPattern.MatchString(value) {
		return true
	}

	if _, ok := namedColors[lower]; ok {
		return true
	}

	return false
}

var warnedInvalidColors sync.Map

func warnInvalidColor(envVar, value string) {
	key := envVar + "\x00" + value
	if _, already := warnedInvalidColors.LoadOrStore(key, struct{}{}); already {
		return
	}

	fmt.Fprintf(
		os.Stderr,
		"uts: %s=%q is not a valid color (expected #abc, #abcdef, #abcdef12, a named color like \"red\", or 0-255); using default\n",
		envVar,
		value,
	)
}

func envColor(envName string) (string, bool) {
	raw := strings.TrimSpace(os.Getenv(envName))
	if raw == "" {
		return "", false
	}

	if !isValidColor(raw) {
		warnInvalidColor(envName, raw)

		return "", false
	}

	if !strings.HasPrefix(raw, "#") && hexColorPattern.MatchString(raw) {
		_, err := strconv.Atoi(raw)
		if err != nil {
			raw = "#" + raw
		}
	}

	return raw, true
}

// ColorEnabled reports whether color output is enabled.
func ColorEnabled() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}

	if _, ok := os.LookupEnv("FORCE_COLOR"); ok {
		return true
	}

	return colorprofile.Detect(os.Stdout, os.Environ()) != colorprofile.Ascii
}
