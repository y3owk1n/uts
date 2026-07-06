// Package doctor checks that external tools required by uts are available.
//
//nolint:goconst
package doctor

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// tool describes an external dependency.
type tool struct {
	Name     string
	Required bool
	UsedBy   string
	Version  string // populated at runtime
}

var baseTools = []tool{
	{Name: "ffmpeg", Required: true, UsedBy: "video, audio"},
	{Name: "gs", Required: false, UsedBy: "pdf compress"},
	{Name: "pngquant", Required: false, UsedBy: "image compress (PNG)"},
	{Name: "optipng", Required: false, UsedBy: "image compress (PNG)"},
	{Name: "jpegoptim", Required: false, UsedBy: "image compress (JPEG)"},
	{Name: "cwebp", Required: false, UsedBy: "image compress/convert (WebP)"},
	{Name: "gifsicle", Required: false, UsedBy: "image compress (GIF)"},
	{Name: "heif-convert", Required: false, UsedBy: "image compress (HEIC)"},
	{Name: "cavif", Required: false, UsedBy: "image compress (AVIF)"},
	{Name: "avifenc", Required: false, UsedBy: "image compress (AVIF)"},
	{Name: "magick", Required: false, UsedBy: "image (ImageMagick)"},
	{Name: "tar", Required: false, UsedBy: "archive"},
	{Name: "zip", Required: false, UsedBy: "archive compress (zip)"},
	{Name: "unzip", Required: false, UsedBy: "archive extract (zip)"},
	{Name: "zstd", Required: false, UsedBy: "archive compress/extract (zstd)"},
}

// Run executes the doctor check and prints the results.
func Run(version string) {
	palette := ui.Style.Palette()

	_, _ = lipgloss.Fprint(os.Stdout, ui.Banner.Logo(version))

	missing := 0
	found := 0

	for i := range baseTools {
		t := &baseTools[i]
		probe(t)

		if t.Version != "" {
			found++
		} else {
			missing++
		}
	}

	// Print results grouped: found first, then missing.
	_, _ = lipgloss.Fprint(os.Stdout, ui.Panel.Section("Found", renderFound(palette)))

	if missing > 0 {
		_, _ = lipgloss.Fprint(os.Stdout, ui.Panel.Section("Missing", renderMissing(palette)))
	}

	// Summary.
	if missing == 0 {
		ui.Message.Successf("All tools available — you're good to go!")
	} else {
		ui.Message.Warnf(
			"%d tool(s) missing — some features may not work.\n  Install with your package manager (e.g. brew install <tool>).",
			missing,
		)
	}
}

func probe(tool *tool) {
	path, err := exec.LookPath(tool.Name)
	if err != nil {
		return
	}

	tool.Version = guessVersion(path, tool.Name)
}

func guessVersion(path, _ string) string {
	// Try common version flags.
	for _, flag := range []string{"--version", "-version", "-v", "version"} {
		out, err := exec.CommandContext(context.Background(), path, flag).CombinedOutput()
		if err != nil {
			continue
		}

		line := firstLine(string(out))
		if line != "" {
			return line
		}
	}

	return "installed"
}

func firstLine(value string) string {
	value = strings.TrimSpace(value)
	if idx := strings.IndexAny(value, "\n\r"); idx > 0 {
		return strings.TrimSpace(value[:idx])
	}

	return value
}

// renderFound returns the styled output for found tools.
func renderFound(palette style.Palette) string {
	var builder strings.Builder

	okStyle := lipgloss.NewStyle().Foreground(palette.Success)
	nameStyle := lipgloss.NewStyle().Foreground(palette.Text)
	verStyle := lipgloss.NewStyle().Foreground(palette.Accent)
	usedByStyle := lipgloss.NewStyle().Foreground(palette.Muted)

	for _, tool := range baseTools {
		if tool.Version == "" {
			continue
		}

		builder.WriteString(okStyle.Render("  ✓ "))
		builder.WriteString(nameStyle.Render(tool.Name))
		builder.WriteString(verStyle.Render("  " + tool.Version))
		builder.WriteString(usedByStyle.Render("  (" + tool.UsedBy + ")"))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderMissing returns the styled output for missing tools.
func renderMissing(palette style.Palette) string {
	var builder strings.Builder

	errStyle := lipgloss.NewStyle().Foreground(palette.Error)
	nameStyle := lipgloss.NewStyle().Foreground(palette.Text)
	reqStyle := lipgloss.NewStyle().Foreground(palette.Warning)
	usedByStyle := lipgloss.NewStyle().Foreground(palette.Muted)
	hintStyle := lipgloss.NewStyle().Foreground(palette.Accent)

	installHint := installHint()

	for _, tool := range baseTools {
		if tool.Version != "" {
			continue
		}

		builder.WriteString(errStyle.Render("  ✗ "))
		builder.WriteString(nameStyle.Render(tool.Name))
		builder.WriteString(usedByStyle.Render("  (" + tool.UsedBy + ")"))

		if tool.Required {
			builder.WriteString(reqStyle.Render("  [required]"))
		}

		builder.WriteString("\n")
	}

	builder.WriteString("\n")
	builder.WriteString(hintStyle.Render("  " + installHint))

	return builder.String()
}

func installHint() string {
	switch runtime.GOOS {
	case "darwin":
		return "Install with: brew install " + brewList()
	case "linux":
		return "Install with your package manager (apt, dnf, pacman, etc.)"
	default:
		return "Install the missing tools using your platform's package manager."
	}
}

func brewList() string {
	var names []string
	for _, t := range baseTools {
		if t.Version == "" {
			names = append(names, t.Name)
		}
	}

	// Map some tool names to brew package names.
	remap := map[string]string{
		"gs":           "ghostscript",
		"magick":       "imagemagick",
		"heif-convert": "libheif",
		"avifenc":      "libavif",
		"cavif":        "cavif",
	}

	for i, n := range names {
		if mapped, ok := remap[n]; ok {
			names[i] = mapped
		}
	}

	return strings.Join(names, " ")
}
