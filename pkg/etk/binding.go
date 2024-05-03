package etk

/*
// Ns provides the etk: module, an Elvish binding for this TUI framework.
var Ns = eval.BuildNsNamed("etk").
	AddVars(map[string]vars.Var{
		"codearea": vars.NewReadOnly(CodeArea),
		"listbox":  vars.NewReadOnly(ListBox),
	}).
	AddGoFns(map[string]any{
		"comp":            comp,
		"vbox":            vbox,
		"handle":          handle,
		"adapt-to-widget": adaptToWidget,

		"text-buffer": func(content string, dot int) tk.CodeBuffer {
			return tk.CodeBuffer{Content: content, Dot: dot}
		},
	}).Ns()

func comp(fm *eval.Frame, fn eval.Callable) Comp {
	return func(c Context) (View, React) {
		subcompViews := map[string]View{}
		subcompReacts := map[string]React{}
		var this = eval.BuildNs().AddVars(map[string]vars.Var{
			"state": stateSubTreeVar(c),
			"subcomp": vars.FromGet(func() any {
				m := vals.EmptyMap
				for k, v := range subcomps {
					m = m.Assoc(k, v)
				}
				return m
			}),
		}).AddGoFns(map[string]any{
			"state": func(name string, _eq string, init any) {
				State(c, name, init)
			},
			"subcomp": func(name string, f Comp, setStatesMap vals.Map) (Scene, error) {
				setStates, err := convertSetStates(setStatesMap)
				if err != nil {
					return Scene{}, err
				}
				el := c.Subcomp(name, WithStates(f, setStates...))
				subcomps[name] = el
				return el, nil
			},
		}).Ns()
		p1, getOut, err := eval.ValueCapturePort()
		if err != nil {
			return errElement(err)
		}
		err = fm.Evaler.Call(fn, eval.CallCfg{Args: []any{this}},
			eval.EvalCfg{Ports: []*eval.Port{nil, p1, nil}})
		if err != nil {
			return errElement(err)
		}
		outs := getOut()
		if len(outs) != 1 {
			return errElement(fmt.Errorf("should only have one output"))
		}
		el, ok := outs[0].(Scene)
		if !ok {
			return errElement(fmt.Errorf("output should be element"))
		}
		return el
	}
}

type stateSubTreeVar Context

func (v stateSubTreeVar) Get() any {
	return getPath(*v.state, v.path)
}

func (v stateSubTreeVar) Set(val any) error {
	valMap, ok := val.(vals.Map)
	if !ok {
		return fmt.Errorf("must be map")
	}
	*v.state = assocPath(*v.state, v.path, valMap)
	return nil
}

func convertSetStates(m vals.Map) ([]any, error) {
	var setStates []any
	for it := m.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		name, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("key should be string")
		}
		setStates = append(setStates, name, v)
	}
	return setStates, nil
}

func vbox(fm *eval.Frame, rowsList vals.List, propsMap vals.Map) Scene {
	var rows []View
	for it := rowsList.Iterator(); it.HasElem(); it.Next() {
		elem, ok := it.Elem().(Scene)
		if !ok {
			return errElement(fmt.Errorf("vbox needs elements"))
		}
		rows = append(rows, elem.View)
	}

	focusAny, ok := propsMap.Index("focus")
	if !ok {
		return errElement(fmt.Errorf("vbox needs focus"))
	}
	focus, ok := focusAny.(int)
	if !ok {
		return errElement(fmt.Errorf("vbox needs int focus"))
	}

	handlerAny, ok := propsMap.Index("handler")
	if !ok {
		return errElement(fmt.Errorf("vbox needs handler"))
	}
	handler, ok := handlerAny.(eval.Callable)
	if !ok {
		return errElement(fmt.Errorf("vbox needs callable handler"))
	}

	return VBoxView{Rows: rows, Focus: focus}.WithHandler(func(e term.Event) Action {
		s := ""
		if ke, ok := e.(term.KeyEvent); ok {
			s = ui.Key(ke).String()
		}

		p1, getOut, err := eval.ValueCapturePort()
		if err != nil {
			// How do I indicate error here 😨
			Notify(ui.T(fmt.Sprintf("value capture port error: %s", err)))
			return Errored
		}
		err = fm.Evaler.Call(handler, eval.CallCfg{Args: []any{e, s}},
			eval.EvalCfg{Ports: []*eval.Port{nil, p1, nil}})
		if err != nil {
			// How do I indicate error here 😨
			var sb strings.Builder
			diag.ShowError(&sb, err)
			Notify(ui.T("handler exception"))
			Notify(ui.ParseSGREscapedText(sb.String()))
			return Errored
		}
		for _, out := range getOut() {
			if action, ok := out.(Action); ok {
				return action
			}
		}
		// TODO: Error when there's no than one Action output
		return Errored
	})
}

func errElement(err error) Scene {
	return Text(ui.T(err.Error(), ui.FgRed)).WithHandler(func(term.Event) Action { return Unused })
}

func handle(el Scene, ev term.Event) Action {
	return el.React(ev)
}

func adaptToWidget(f func(Context) Scene) tk.Widget {
	return AdaptToWidget(f)
}
*/
