package main

import (
	"strconv"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/ui"
)

func Counter(c etk.Context) (etk.View, etk.React) {
	valueVar := etk.State(c, "value", 0)
	buttonView, buttonReact := c.Subcomp("button", etk.WithStates(Button,
		"label", "Count",
		"submit", func() etk.Action {
			valueVar.Swap(func(i int) int { return i + 1 })
			return etk.Consumed
		},
	))

	return etk.HBoxFlex(
			etk.Text(ui.T(strconv.Itoa(valueVar.Get()))),
			buttonView,
		).WithFocus(1).WithGap(1),
		buttonReact
}

func Button(c etk.Context) (etk.View, etk.React) {
	labelVar := etk.State(c, "label", "button")
	submitVar := etk.State(c, "submit", func() etk.Action { return etk.Unused })
	return etk.Text(ui.T("[ "+labelVar.Get()+" ]", ui.Inverse)),
		c.WithBinding("button", func(ev term.Event) etk.Action {
			if ev == term.K(' ') || ev == term.K(ui.Enter) {
				return submitVar.Get()()
			}
			return etk.Unused
		})
}
