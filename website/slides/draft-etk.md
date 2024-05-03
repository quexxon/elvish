# The design and implementation of the Etk TUI framework: a study of unusual tradeoffs

Qi Xiao (xiaq)

(Draft slides)

***

# Motivation

-   Elvish's TUI consists of ~22k LoC (pkg/cli, pkg/edit)

    -   Compared to ~18k LoC of the interpreter (pkg/parse, pkg/eval)

-   The good

    -   Heavily modularized into a hierarchy of components

    -   Simple, uniform `Widget` API

    -   High test coverage

-   The bad: hard to develop and test in Go

    -   Adding a simple functionality often requires "propagating" it through
        multiple layers of internal APIs

    -   Tests are verbose and time-consuming to create and update

-   The ugly: very hard to extend from Elvish code

    -   Adhoc customization/extension points and bespoke bindings

    -   Internal components cannot be reused from Elvish code

    -   No API to write a component in Elvish

***

# Immediate mode UI

-   Originally: for each frame, clear the screen and redraw everything

-   Immediate mode emulation over retained mode primitives

    -   Components are "just functions"

    -   React, SwiftUI, Jetpack Compose

***

# Component is "just function"

-   Hide complexity

    -   Bad

-   Encourages fine-grained componentization

    -   Good

***

# State management in immediate mode

-   Local variables don't survive the next call

-   Solution: Store it somewhere else™️

    -   YOLO, just use global variables

    -   Use global variables, but declared inside components (SwiftUI, Jetpack
        Compose)

    -   Use global variables, but hidden by the framework (React)

-   Etk's solution

    -   Use global variables, organized into a tree

***

# State tree in Etk

-   Backed by nested persistent maps

    -   Immutable

    -   But cheap to create variation of

    -   Undo/redo comes (almost) for free

***

# TUI primitives

-   Terminal is mostly text

-   In-band signals: escape codes

    -   Combination and function keys

    -   Text style

    -   Cursor addressing

-   Out-of-band control: signals and `ioctl`

-   Unix's terminal API dates back to the 1960s (TTY = **t**ele**ty**pewritter)

    -   Various unsuccessful reform attempts throughout history

***

# Implementation

-   The Etk core is only ~X00 LoC, or ~Y00 effective LoC

    -   Doesn't include terminal primitives or common components

***

# An open system

-   Inspectable

-   Tinkerable

-   Not a sealed "product"

***

# Unsafety

-   State bindings in Go are not type-safe

-   A conscious choice

    -   Support the desired style

    -   Any state could be changed by Elvish code to any value

-   Subcomponent and state names are just strings

    -   Not ideal

***

# Non-encapsulation

-   A component can mutate the state of any of its descendant

    -   Even which subcomponent a component uses is a mutable state

-   Access to the root allows you to inspect and mutate any point of the state
    tree

-   Again, a conscious choice

    -   The entire state tree is exposed in the Elvish binding

    -   Each TUI app gets an API for free

***

# Sample code (Go)

```go
// TODO: source of the Counter component
```

***

# Sample code (Elvish)

```elvish
# TODO: source of the Counter component
```

***

# Keybinding

-   Different components use different keybindings

-   No problem, just declaring binding as a state:

    ```go
    func MyComp(c etk.Context) etk.Scene {
        bindingVar := etk.State(c, "binding",
            func(term.Event) etk.Action { return etk.Unused })

        return Scene{
            View: someView,
            React: func(ev term.Event) etk.Action {
                action := bindingVar.Get()(ev)
                if action != etk.Unused {
                    return action
                }
                /* Custom logic */
            },
        }
    }
    ```

***

# Keybinding (cont.)

-   But multiple instances of the same component should share a keybinding

    -   If I bind <kbd>Left</kbd> to move the cursor one character left, it
        should apply to all CodeArea components

    -   Some individual instances can have another overlay

***

# Concurrency safety

-   Imperative style is very natural for expressing state changes

    -   And not very good for concurrency safety

-   Guard each mutation with a mutex - simple

    -   Concurrent writes are not properly isolated

-   RETVRN to Elm's purely functional style?

    -   Function signatures become more complex

    -   Event handler can change the state

        -   Need to manually propagate state changes from subcomponents

    -   Initialization can change the state too

-   If mutexes don't solve your problem, you aren't using enough mutexes

    -   Two mutexes

***

# Name safety revisited

-   Simple typos can't be prevented

    -   Let's declare them as Go identifiers instead

    -   But...

-   Compare this:

    ```go
    // TODO
    ```

-   With this:

    ```go
    // TODO
    ```

-   The former is less safe, but highlights the component and state hierarchy
    very clearly

-   Declared keywords when?

***

# Cost vs benefits

-   We give up

    -   Type and name safety

    -   Ability to enforce any invariant at all

-   We get

    -   Concise code

    -   First-class Elvish binding

    -   Undo/redo for free

-   We still keep

    -   Concurrency safety
