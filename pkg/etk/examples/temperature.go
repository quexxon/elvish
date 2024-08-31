package main

import (
	"fmt"
	"strconv"

	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/ui"
)

func Temperature(c etk.Context) (etk.View, etk.React) {
	celsiusView, celsiusReact := c.Subcomp("celsius", etk.WithStates(etk.CodeArea, "prompt", ui.T("Celsius: ")))
	celsiusBufferVar := etk.BindState[tk.CodeBuffer](c, "celsius", "buffer")

	fahrenheitView, fahrenheitReact := c.Subcomp("fahrenheit", etk.WithStates(etk.CodeArea, "prompt", ui.T("Fahrenheit: ")))
	fahrenheitBufferVar := etk.BindState[tk.CodeBuffer](c, "fahrenheit", "buffer")

	focusVar := etk.State(c, "focus", 0)

	return etk.VBox(celsiusView, fahrenheitView).WithFocus(focusVar.Get()),
		func(e term.Event) etk.Action {
			focus := focusVar.Get()
			if e == term.K(ui.Tab) {
				focusVar.Set(1 - focus)
				return etk.Consumed
			}
			if focus == 0 {
				if celsiusReact(e) == etk.Consumed {
					if c, err := strconv.ParseFloat(celsiusBufferVar.Get().Content, 64); err == nil {
						f := fmt.Sprintf("%.2f", c*9/5+32)
						fahrenheitBufferVar.Set(tk.CodeBuffer{Content: f, Dot: len(f)})
					}
					return etk.Consumed
				}
			} else {
				if fahrenheitReact(e) == etk.Consumed {
					if f, err := strconv.ParseFloat(fahrenheitBufferVar.Get().Content, 64); err == nil {
						c := fmt.Sprintf("%.2f", (f-32)*5/9)
						celsiusBufferVar.Set(tk.CodeBuffer{Content: c, Dot: len(c)})
					}
					return etk.Consumed
				}
			}
			return etk.Unused
		}
}
