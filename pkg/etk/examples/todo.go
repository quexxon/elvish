package main

import (
	"fmt"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/ui"
)

type todoItem struct {
	text string
	done bool
}

type todoItems []todoItem

func (ti todoItems) Len() int { return len(ti) }
func (ti todoItems) Show(i int) ui.Text {
	done := ' '
	if ti[i].done {
		done = 'X'
	}
	return ui.T(fmt.Sprintf("[%c] %s", done, ti[i].text))
}

func Todo(c etk.Context) (etk.View, etk.React) {
	// TODO: API to combine init and bind
	listView, listReact := c.Subcomp("list", etk.WithStates(etk.ListBox, "items", todoItems{}))
	itemsVar := etk.BindState[todoItems](c, "list", "items")
	selectedVar := etk.BindState[int](c, "list", "selected")

	newItemView, newItemReact := c.Subcomp("new-item", etk.WithStates(etk.CodeArea, "prompt", ui.T("new item: ")))
	bufferVar := etk.BindState[tk.CodeBuffer](c, "new-item", "buffer")

	focusVar := etk.State(c, "focus", 1)
	focus := focusVar.Get()

	return etk.VBox(listView, newItemView).WithFocus(focus), func(e term.Event) etk.Action {
		if e == term.K(ui.Tab) {
			focusVar.Set(1 - focus)
			return etk.Consumed
		}
		if focus == 0 {
			action := listReact(e)
			if action == etk.Unused {
				switch e {
				case term.K(ui.Down):
					focusVar.Set(1)
					return etk.Consumed
				case term.K(' '):
					item := &itemsVar.Get()[selectedVar.Get()]
					item.done = !item.done
					return etk.Consumed
				}
			}
			return action
		} else {
			action := newItemReact(e)
			if action == etk.Unused {
				switch e {
				case term.K(ui.Up):
					focusVar.Set(0)
					return etk.Consumed
				case term.K(ui.Enter):
					itemsVar.Set(append(itemsVar.Get(), todoItem{text: bufferVar.Get().Content}))
					bufferVar.Set(tk.CodeBuffer{})
					return etk.Consumed
				}
			}
			return action
		}
	}
}
