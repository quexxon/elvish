package main

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/ui"
)

func Flight(c etk.Context) (etk.View, etk.React) {
	// TODO: Horizontal
	typeView, typeReact := c.Subcomp("type",
		etk.WithStates(etk.ListBox, "items", etk.StringItems("one-way", "return")))
	outboundView, outboundReact := c.Subcomp("outbound",
		etk.WithStates(etk.CodeArea, "prompt", ui.T("outbound: ")))
	// TODO: Disable inbound for one-way
	inboundView, inboundReact := c.Subcomp("inbound",
		etk.WithStates(etk.CodeArea, "prompt", ui.T("inbound:  ")))
	bookView, bookReact := c.Subcomp("book",
		etk.WithStates(Button, "label", "Book", "submit", func() {
		}))

	focusVar := etk.State(c, "focus", 0)
	focus := focusVar.Get()
	return etk.VBox(
			typeView, outboundView, inboundView, bookView,
		).WithFocus(focus),
		func(ev term.Event) etk.Action {
			action := []etk.React{
				typeReact, outboundReact, inboundReact, bookReact,
			}[focus](ev)
			if action == etk.Unused {
				switch ev {
				case term.K(ui.Down), term.K(ui.Tab):
					if focus < 3 {
						focusVar.Set(focus + 1)
					}
				case term.K(ui.Up), term.K(ui.Tab, ui.Shift):
					if focus > 0 {
						focusVar.Set(focus - 1)
					}
				}
			}
			return action
		}
}
