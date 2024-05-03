package edit

import (
	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/ui"
)

type Editor struct {
	tty cli.TTY
}

func NewEditor(tty cli.TTY) *Editor {
	return &Editor{tty}
}

func (ed *Editor) ReadCode() (string, error) {
	m, err := etk.Run(ed.tty, etk.WithStates(etk.CodeArea, "prompt", ui.T("etkedit> ")),
		func(ev term.Event, c etk.Context, tag string, f etk.React) etk.Action {
			if ev == term.K('[', ui.Ctrl) {
				return etk.Exit
			}
			return f(ev)
		})
	if err != nil {
		return "", err
	}
	buf, _ := m.Index("buffer")
	return buf.(etk.CodeBuffer).Content, nil
}

func (ed *Editor) RunAfterCommandHooks(src parse.Source, duration float64, err error) {
	// TODO
}
