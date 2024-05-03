package etk

import (
	"src.elv.sh/pkg/cli/term"
)

func ComboBox(c Context) (View, React) {
	filterView, filterReact := c.Subcomp("filter", CodeArea)
	filterBufferVar := BindState[CodeBuffer](c, "filter", "buffer")
	listView, listReact := c.Subcomp("list", ListBox)
	listItemsVar := BindState[Items](c, "list", "items")
	listSelectedVar := BindState[int](c, "list", "selected")

	genListVar := State(c, "gen-list", func(string) (Items, int) {
		return nil, -1
	})
	lastFilterContentVar := State(c, "-last-filter-content", "")

	return VBox(filterView, listView).WithFocus(0),
		c.WithBinding("combobox", func(ev term.Event) Action {
			if action := filterReact(ev); action != Unused {
				filterContent := filterBufferVar.Get().Content
				if filterContent != lastFilterContentVar.Get() {
					lastFilterContentVar.Set(filterContent)
					items, selected := genListVar.Get()(filterContent)
					listItemsVar.Set(items)
					listSelectedVar.Set(selected)
				}
				return action
			} else {
				return listReact(ev)
			}
		})
}
