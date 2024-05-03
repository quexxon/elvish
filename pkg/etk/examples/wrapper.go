package main

import (
	"strings"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/ui"
)

func Wrapper(c etk.Context) (etk.View, etk.React) {
	innerView, innerReact := c.Subcomp("inner", nop)
	innerStateVar := etk.BindState[vals.Map](c, "inner")

	stateText := ui.T(strings.ReplaceAll(vals.Repr(innerStateVar.Get(), 0), "\t", " "))
	return etk.VBox(innerView, etk.Text(stateText)).WithFocus(0),
		func(e term.Event) etk.Action {
			action := innerReact(e)
			if action == etk.Unused && (e == term.K('[', ui.Ctrl)) {
				_ = ui.Tab
				return etk.Exit
			}
			return action
		}
}

func nop(c etk.Context) (etk.View, etk.React) {
	return etk.Empty, func(term.Event) etk.Action { return etk.Unused }
}
