package ui

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/y3owk1n/uts/internal/ui/message"
	"github.com/y3owk1n/uts/internal/ui/style"
	"github.com/y3owk1n/uts/internal/ui/table"
)

var Message = message.Default()

var Style = styleAPI{}

var Table = tableAPI{}

type styleAPI struct{}

func (styleAPI) Palette() style.Palette { return style.Default() }
func (styleAPI) ColorEnabled() bool     { return style.ColorEnabled() }

type tableAPI struct{}

func (tableAPI) New(headers ...string) *table.Table { return table.New(headers...) }

func IsTTY() bool {
	return term.IsTerminal(os.Stdout.Fd())
}
