package etk

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/eval/vals"
)

func AdaptToWidget(f Comp) tk.Widget {
	state := vals.EmptyMap
	view, react := f(Context{&state, nil, nil})
	return &widget{&state, view, react, f}
}

type widget struct {
	state *vals.Map
	view  View
	react React
	f     Comp
}

func (w *widget) Render(width, height int) *term.Buffer {
	return w.view.Render(width, height)
}

func (w *widget) MaxHeight(width, height int) int {
	return height
}

func (w *widget) Handle(event term.Event) bool {
	action := w.react(event)
	w.view, w.react = w.f(Context{w.state, nil, nil})
	return action != Unused
}
