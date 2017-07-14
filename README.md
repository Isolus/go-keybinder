go-keybinder [![GoDoc](https://godoc.org/github.com/Isolus/go-keybinder?status.svg)](https://godoc.org/github.com/Isolus/go-keybinder)
=====

The go-keybinder project provides Go bindings for keybinder-3.0.

## Sample Use

The following example can be found in [Examples](https://github.com/Isolus/go-keybinder/examples/).

```Go
package main

import (
	"log"
	"unsafe"

	"github.com/Isolus/go-keybinder"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var win *gtk.Window

func keybinderHandler(keystring string, data unsafe.Pointer) {

	if win.IsVisible() {
		glib.IdleAdd(win.Hide)
	} else {
		glib.IdleAdd(win.Show)
	}
	log.Printf("Pressed: %s", keystring)
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	keybinder.Init()
	keybinder.Bind("<Ctrl>space", keybinderHandler, nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	var err error
	win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	l, err := gtk.LabelNew("Hello, gotk3!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Add the label to the window.
	win.Add(l)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
```

To build the example:

```
$ go build example.go

```

## Documentation

The internal `go doc` style documentation can be viewed
online without installing this package by using the [GoDoc site](http://godoc.org/github.com/Isolus/go-keybinder).

You can also view the documentation locally once the package is
installed with the `godoc` tool by running `godoc -http=":6060"` and
pointing your browser to
http://localhost:6060/pkg/github.com/Isolus/go-keybinder

## Installation

go-keybinder currently requires keybinder-3.0. A recent Go is also required.

## TODO
- Add bindings for keybinder\_bind\_full().

## License

Package go-keybinder is licensed under the MIT License.