package ui

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/y3owk1n/uts/internal/ui/banner"
	"github.com/y3owk1n/uts/internal/ui/message"
	"github.com/y3owk1n/uts/internal/ui/panel"
	"github.com/y3owk1n/uts/internal/ui/style"
	"github.com/y3owk1n/uts/internal/ui/table"
)

// Message is the default message printer.
var Message = message.Default()

// Banner provides banner rendering helpers.
var Banner = bannerAPI{}

// Panel provides panel rendering helpers.
var Panel = panelAPI{}

// Style provides style helpers.
var Style = styleAPI{}

// Table provides table helpers.
var Table = tableAPI{}

// bannerAPI implements banner operations.
type bannerAPI struct{}

func (bannerAPI) Logo(version string) string { return banner.Logo(style.Default(), version) }
func (bannerAPI) Header(text string) string  { return banner.Header(style.Default(), text) }

// panelAPI implements panel operations.
type panelAPI struct{}

func (panelAPI) Panel(content string) string {
	return panel.Panel(style.Default(), content)
}

func (panelAPI) Section(title, content string) string {
	return panel.Section(style.Default(), title, content)
}

// styleAPI implements style operations.
type styleAPI struct{}

func (styleAPI) Palette() style.Palette { return style.Default() }
func (styleAPI) ColorEnabled() bool     { return style.ColorEnabled() }

// tableAPI implements table operations.
type tableAPI struct{}

func (tableAPI) New(headers ...string) *table.Table { return table.New(headers...) }

// IsTTY reports whether the output is a terminal.
func IsTTY() bool {
	return term.IsTerminal(os.Stdout.Fd())
}
