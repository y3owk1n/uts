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

var Message = message.Default()

var Banner = bannerAPI{}

var Panel = panelAPI{}

var Style = styleAPI{}

var Table = tableAPI{}

type bannerAPI struct{}

func (bannerAPI) Logo(version string) string { return banner.Logo(style.Default(), version) }
func (bannerAPI) Header(text string) string  { return banner.Header(style.Default(), text) }

type panelAPI struct{}

func (panelAPI) Panel(content string) string {
	return panel.Panel(style.Default(), content)
}

func (panelAPI) Section(title, content string) string {
	return panel.Section(style.Default(), title, content)
}

type styleAPI struct{}

func (styleAPI) Palette() style.Palette { return style.Default() }
func (styleAPI) ColorEnabled() bool     { return style.ColorEnabled() }

type tableAPI struct{}

func (tableAPI) New(headers ...string) *table.Table { return table.New(headers...) }

func IsTTY() bool {
	return term.IsTerminal(os.Stdout.Fd())
}
