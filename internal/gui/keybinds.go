package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/sirupsen/logrus"
)

// Binding - a keybinding mapping a key and modifier to a handler. The keypress
// is only handled if the given view has focus, or handled globally if the view
// is ""
type Binding struct {
	ViewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         any
	Modifier    gocui.Modifier
	Description string
}

func (gui *Gui) GetInitialKeybindings() []*Binding {
	bindings := []*Binding{
		{
			ViewName: "",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		},
		{
			ViewName: "",
			Key:      ']',
			Modifier: gocui.ModNone,
			Handler:  gui.nextTab,
		},
		{
			ViewName: "",
			Key:      '[',
			Modifier: gocui.ModNone,
			Handler:  gui.previousTab,
		},
		{
			ViewName: "seasons",
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.nextItem,
		},
		{
			ViewName: "seasons",
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.previousItem,
		},
	}

	return bindings
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
		gui.Log.WithFields(logrus.Fields{
			"view":     binding.ViewName,
			"key":      binding.Key,
			"modifier": binding.Modifier,
		}).Debug("Registered keybinding")
	}

	return nil
}
