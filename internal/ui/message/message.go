package message

import (
	"fmt"
	"image/color"
	"io"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/y3owk1n/uts/internal/ui/style"
)

// Icons defines the set of icons used in messages.
type Icons struct {
	Info    string
	Success string
	Warn    string
	Error   string
	Step    string
	Bullet  string
	Arrow   string
}

// DefaultIcons returns the default icon set.
func DefaultIcons() Icons {
	return Icons{
		Info:    "ℹ",
		Success: "✓",
		Warn:    "⚠",
		Error:   "✖",
		Step:    "▸",
		Bullet:  "•",
		Arrow:   "→",
	}
}

// Printer formats and writes styled messages.
type Printer struct {
	palette style.Palette
	types   style.Type
	icons   Icons
	out     io.Writer
	errOut  io.Writer
}

// New creates a new Printer.
func New(palette style.Palette, types style.Type, icons Icons, out, errOut io.Writer) *Printer {
	if out == nil {
		out = os.Stdout
	}

	if errOut == nil {
		errOut = os.Stderr
	}

	return &Printer{
		palette: palette,
		types:   types,
		icons:   icons,
		out:     out,
		errOut:  errOut,
	}
}

// Default returns a Printer with default settings.
func Default() *Printer {
	palette := style.Default()
	types := style.Types(palette)

	return New(palette, types, DefaultIcons(), os.Stdout, os.Stderr)
}

// Infof prints an info message.
func (pr *Printer) Infof(format string, args ...any) {
	pr.line(pr.out, pr.icons.Info, fmt.Sprintf(format, args...), pr.palette.Accent)
}

// Successf prints a success message.
func (pr *Printer) Successf(format string, args ...any) {
	pr.line(pr.out, pr.icons.Success, fmt.Sprintf(format, args...), pr.palette.Success)
}

// Warnf prints a warning message.
func (pr *Printer) Warnf(format string, args ...any) {
	pr.line(pr.out, pr.icons.Warn, fmt.Sprintf(format, args...), pr.palette.Warning)
}

// Errorf prints an error message.
func (pr *Printer) Errorf(format string, args ...any) {
	pr.line(pr.errOut, pr.icons.Error, fmt.Sprintf(format, args...), pr.palette.Error)
}

// Stepf prints a step message.
func (pr *Printer) Stepf(format string, args ...any) {
	pr.line(pr.out, pr.icons.Step, fmt.Sprintf(format, args...), pr.palette.Primary)
}

// Bulletf prints a bullet point.
func (pr *Printer) Bulletf(format string, args ...any) {
	styled := lipgloss.NewStyle().
		Foreground(pr.palette.Muted).
		Render("  " + pr.icons.Bullet + " " + fmt.Sprintf(format, args...))
	//nolint:errcheck
	lipgloss.Fprintln(pr.out, styled)
}

// Mutedf prints a muted message.
func (pr *Printer) Mutedf(format string, args ...any) {
	//nolint:errcheck
	lipgloss.Fprintln(pr.out, pr.types.Muted.Render(fmt.Sprintf(format, args...)))
}

// Pair prints a key-value pair.
func (pr *Printer) Pair(key, value string) {
	styledKey := pr.types.Key.Render(key)
	styledVal := pr.types.Code.Render(value)
	//nolint:errcheck
	lipgloss.Fprintf(pr.out, "%s  %s\n", styledKey, styledVal)
}

// PairLine returns a key-value pair as a string.
func (pr *Printer) PairLine(key, value string) string {
	return pr.types.Key.Render(key) + "  " + pr.types.Code.Render(value) + "\n"
}

// Highlight returns bold primary-colored text.
func (pr *Printer) Highlight(text string) string {
	return lipgloss.NewStyle().Bold(true).Foreground(pr.palette.Primary).Render(text)
}

// Dim returns muted text.
func (pr *Printer) Dim(text string) string {
	return pr.types.Muted.Render(text)
}

// Success returns green success-colored text.
func (pr *Printer) Success(text string) string {
	return lipgloss.NewStyle().Foreground(pr.palette.Success).Render(text)
}

// Warn returns yellow warning-colored text.
func (pr *Printer) Warn(text string) string {
	return lipgloss.NewStyle().Foreground(pr.palette.Warning).Render(text)
}

// Error returns red error-colored text.
func (pr *Printer) Error(text string) string {
	return lipgloss.NewStyle().Foreground(pr.palette.Error).Render(text)
}

// Accent returns accent-colored text.
func (pr *Printer) Accent(text string) string {
	return lipgloss.NewStyle().Foreground(pr.palette.Accent).Render(text)
}

// Text returns standard text-colored text.
func (pr *Printer) Text(text string) string {
	return lipgloss.NewStyle().Foreground(pr.palette.Text).Render(text)
}

// SuccessRow returns a success row with icon.
func (pr *Printer) SuccessRow(text string) string {
	return pr.styledIconLine(
		pr.icons.Success,
		pr.palette.Success,
	) + pr.types.Body.Render(
		text,
	) + "\n"
}

// WarnRow returns a warning row with icon.
func (pr *Printer) WarnRow(text string) string {
	return pr.styledIconLine(pr.icons.Warn, pr.palette.Warning) + pr.types.Body.Render(text) + "\n"
}

// ErrorRow returns an error row with icon.
func (pr *Printer) ErrorRow(text string) string {
	return pr.styledIconLine(pr.icons.Error, pr.palette.Error) + pr.types.Body.Render(text) + "\n"
}

// Detail returns indented muted text.
func (pr *Printer) Detail(text string) string {
	return pr.types.Muted.Render("    "+text) + "\n"
}

func (pr *Printer) styledIconLine(icon string, color color.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(icon) + " "
}

func (pr *Printer) line(
	writer io.Writer,
	icon string,
	message string,
	iconColor color.Color,
) {
	styledIcon := lipgloss.NewStyle().Bold(true).Foreground(iconColor).Render(icon)
	styledMsg := pr.types.Body.Render(message)
	//nolint:errcheck
	lipgloss.Fprintf(writer, "%s %s\n", styledIcon, styledMsg)
}
