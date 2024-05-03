package main

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/ui"
)

type Task struct {
	Name        string
	Description string
	Code        string
}

type Tasks []Task

func (ts Tasks) Show(i int) ui.Text { return ui.T(ts[i].Name) }

func (ts Tasks) Len() int { return len(ts) }

var tasks = Tasks{
	{"Set up carapace", "Carapace provides a lot of completions.", "sudo brew install carapace"},
	{"Use readline binding", "Keybindings like:\nCtrl-N to next line\nCtrl-P to previous line\nCtrl-F to next character\nCtrl-B to previous character", "use readline-binding"},
}

func Wizard(c etk.Context) (etk.View, etk.React) {
	listView, listReact := c.Subcomp("list", etk.WithStates(etk.ListBox, "items", tasks))
	selectedVar := etk.BindState[int](c, "list", "selected")
	selected := selectedVar.Get()
	description := etk.Text(ui.T(tasks[selected].Description))
	code := etk.Text(ui.T("\n" + tasks[selected].Code))

	return etk.VBox(etk.HBox(listView, description), code).WithFocus(0),
		func(e term.Event) etk.Action {
			if e == term.K(ui.Enter) {
				// TODO: Show notification?
				return etk.Consumed
			}
			return listReact(e)
		}
}
