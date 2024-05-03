package etk

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/ui"
)

// Items is an interface for accessing multiple items.
type Items interface {
	// Show renders the item at the given zero-based index.
	Show(i int) ui.Text
	// Len returns the number of items.
	Len() int
}

type stringItems []string

func StringItems(items ...string) Items   { return stringItems(items) }
func (si stringItems) Show(i int) ui.Text { return ui.T(si[i]) }
func (si stringItems) Len() int           { return len(si) }

func ListBox(c Context) (View, React) {
	itemsVar := State(c, "items", Items(nil))
	selectedVar := State(c, "selected", 0)

	// This is not used anywhere, but declared here for ease of use from a
	// keybinding.
	_ = State(c, "submit", func(Items, int) {})

	selected := selectedVar.Get()
	focus := 0
	var spans []ui.Text
	if items := itemsVar.Get(); items != nil {
		for i := 0; i < items.Len(); i++ {
			if i > 0 {
				spans = append(spans, ui.T("\n"))
			}
			if i == selected {
				focus = len(spans)
				spans = append(spans, ui.StyleText(items.Show(i), ui.Inverse))
			} else {
				spans = append(spans, items.Show(i))
			}
		}
	}

	return Text(spans...).WithDotBefore(focus),
		c.WithBinding("listbox", func(e term.Event) Action {
			selected := selectedVar.Get()
			items := itemsVar.Get()
			switch e {
			case term.K(ui.Up):
				if selected > 0 {
					selectedVar.Set(selected - 1)
					return Consumed
				}
			case term.K(ui.Down):
				if selected < items.Len()-1 {
					selectedVar.Set(selected + 1)
					return Consumed
				}
			}
			return Unused
		})
}
