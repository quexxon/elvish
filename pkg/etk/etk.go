// Package etk implements an [immediate mode] TUI framework with managed states.
//
// Each component in the TUI is implemented by a [Comp]: a function taking a
// [Context] and returning a [Scene]:
//
//   - The [Context] provides access to states associated with the component and
//     supports creating sub-components.
//
//   - The [Scene] is a snapshot of the UI that reflects the current state.
//
// Whenever there is an update in the states, the function is called again to
// generate a new [Scene].
//
// The state is organized into a tree, with individual state variables as leaf
// nodes and components as inner nodes. The [Context] provides access to the
// current level and all descendant levels, allowing a component to manipulate
// not just its own state, but also that of any descendant. This is the only way
// of passing information between components: if a component has any
// customizable property, it is modelled as a state that its parent can modify.
//
// # Design notes
//
// Immediate mode is an alternative to the more common [retained mode] style of
// graphics API. Some GUI frameworks using this style are [Dear ImGui] and [Gio
// UI]. [React], [SwiftUI] and [Jetpack Compose] also provide immediate mode
// APIs above an underlying [retained mode] API.
//
// Immediate mode libraries differ a lot in how component structure and state
// are managed. Etk is used to implement Elvish's terminal UI, so the choices
// made by etk is driven largely by how easy it is to create an Elvish binding
// for the framework that is maximally programmable:
//
//   - The open nature of the state tree makes it easy to inspect and mutate the
//     terminal UI as it is running.
//
//   - The managed nature of the state tree gives us concurrency safety and
//     undo/redo almost for free.
//
//   - The use of [vals.Map] to back the state tree sacrifices type safety in
//     the Go version of the framework, but makes Elvish integration much
//     easier.
//
// [immediate mode]: https://en.wikipedia.org/wiki/Immediate_mode_(computer_graphics)
// [retained mode]: https://en.wikipedia.org/wiki/Retained_mode
// [Dear ImGui]: https://github.com/ocornut/imgui
// [Gio UI]: https://gioui.org
// [React]: https://react.dev
// [SwiftUI]: https://developer.apple.com/xcode/swiftui/
// [Jetpack Compose]: https://developer.android.com/compose
package etk

import (
	"slices"
	"strings"

	"src.elv.sh/pkg/cli"
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/must"
	"src.elv.sh/pkg/ui"
)

// TODO: Automatically remove state that hasn't been referenced.

// For debugging - remove later ❗️
var Notify func(ui.Text)

// Comp is the type for a component. It is called every time the state changes
// to generate a new [Scene].
type Comp func(Context) (View, React)

// WithStates returns a variation of the component that overrides the initial
// value of some state variables. The variadic arguments must come in (key,
// value) pairs, and the keys must be strings.
func WithStates(f Comp, setStates ...any) Comp {
	return func(c Context) (View, React) {
		for i := 0; i < len(setStates); i += 2 {
			key, value := setStates[i].(string), setStates[i+1]
			// TODO: This is not optimal; we shouldn't have to start from the root
			// every time.
			if stateVar := BindState[any](c, strings.Split(key, "/")...); stateVar.GetAny() == nil {
				stateVar.Set(value)
			}
		}
		return f(c)
	}
}

// Scene is a snapshot of the UI that reflects a fixed state.
//
// This type is isomorphic to [tk.Widget], but the latter represents a
// long-living object with a mutable state instead. The difference is best seen
// in how [tk.Widget.Handle] and [Scene.Handler] behaves:
//
//   - Calling the Handle method of a [tk.Widget] causes its internal state to
//     be mutated. The next call to its Render method will reflect the changed
//     state.
//
//   - Calling the Handler field of a Scene causes a state managed elsewhere to
//     be mutated. The next call to the Render method of the Layout field still
//     reflects the previous state. Instead, a new Scene needs to be generated
//     at this point to reflect the changed state.

type View tk.Renderer

type React func(term.Event) Action

type Action uint32

const (
	Unused Action = iota
	Consumed
	Errored
	Exit
)

// Context provides access to the state tree at the current level and all
// descendant levels.
type Context struct {
	state   *vals.Map
	binding func(ev term.Event, c Context, tag string, r React) Action
	path    []string
}

func (c Context) descPath(path ...string) []string {
	return slices.Concat(c.path, path)
}

// Subcomp does the following:
//
//   - Create a map state variable with the given name
//
//   - Create a Comp state variable with the given name plus "-comp", using f as
//     the initial value
//
// It then invokes the component with the map as the context.
func (c Context) Subcomp(name string, f Comp) (View, React) {
	State(c, name, vals.EmptyMap)
	compVar := State(c, name+"-comp", f)
	return compVar.Get()(Context{c.state, c.binding, c.descPath(name)})
}

// TODO: How to make this mandatory?
func (c Context) WithBinding(tag string, f React) React {
	return func(ev term.Event) Action {
		if c.binding != nil {
			return c.binding(ev, c, tag, f)
		}
		return f(ev)
	}
}

// Need to support two things:
//
// - Bind a variable to a path, without initializing it
// - And initialize it
//
// TODO: If the variable has been stored with an incompatible type, augment the
// component layout with an error message

// State returns a state variable with the given name under the current level,
// initializing it to a given value if it doesn't exist yet.
func State[T any](c Context, name string, initial T) StateVar[T] {
	sv := BindState[T](c, name)
	// TODO: Detect when the value exists but is of the wrong type, and log it
	// as an error somewhere (where?)
	if sv.GetAny() == nil {
		sv.Set(initial)
	}
	return sv
}

// BindState returns a state variable with the given path from the current
// level. It doesn't initialize the variable.
//
// This should only be used if the variable is initialized elsewhere, most
// typically for accessing the state of a subcomponent after the subcomponent
// has been called.
func BindState[T any](c Context, path ...string) StateVar[T] {
	return StateVar[T]{c.state, c.descPath(path...)}
}

// StateVar provides access to a state variable, a node in the state tree.
type StateVar[T any] struct {
	state *vals.Map
	path  []string
}

func (sv StateVar[T]) Get() T {
	val := getPath(*sv.state, sv.path)
	var dst T
	must.OK(vals.ScanToGo(val, &dst))
	return dst
}

func (sv StateVar[T]) GetAny() any { return getPath(*sv.state, sv.path) }

func (sv StateVar[T]) Set(t T)          { *sv.state = assocPath(*sv.state, sv.path, t) }
func (sv StateVar[T]) Swap(f func(T) T) { sv.Set(f(sv.Get())) }

func getPath(m vals.Map, path []string) any {
	if len(path) == 0 {
		return m
	}
	for len(path) > 1 {
		v, _ := m.Index(path[0])
		m = v.(vals.Map)
		path = path[1:]
	}
	v, _ := m.Index(path[0])
	return v
}

func assocPath(m vals.Map, path []string, newVal any) vals.Map {
	if len(path) == 0 {
		return newVal.(vals.Map)
	}

	if len(path) == 1 {
		return m.Assoc(path[0], newVal)
	}
	v, _ := m.Index(path[0])
	return m.Assoc(path[0], assocPath(v.(vals.Map), path[1:], newVal))
}

func Run(tty cli.TTY, f Comp, gr func(term.Event, Context, string, React) Action) (vals.Map, error) {
	restore, err := tty.Setup()
	if err != nil {
		return nil, err
	}
	defer restore()

	state := vals.EmptyMap
	view, react := f(Context{&state, gr, nil})

	for {
		h, w := tty.Size()
		buf := view.Render(w, h)
		tty.UpdateBuffer(nil, buf, false)
		event, err := tty.ReadEvent()
		if err != nil {
			return nil, err
		}
		action := react(event)
		if action == Exit {
			tty.UpdateBuffer(nil, term.NewBuffer(w), false)
			break
		}
		view, react = f(Context{&state, gr, nil})
	}
	return state, nil
}
