package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) renderString(g *gocui.Gui, viewName string, s string) error {
	g.Update(func(*gocui.Gui) error {
		v, err := g.View(viewName)
		if err != nil {
			return nil
		}
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
		return gui.setViewContent(v, s)
	})
	return nil
}

func (gui *Gui) setViewContent(v *gocui.View, s string) error {
	v.Clear()
	fmt.Fprint(v, s)
	return nil
}

// BUG: the first item is white and bold
func (gui *Gui) renderTabList(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if len(tab.Data.Data) == 0 {
		return gui.renderString(gui.g, viewName, "Loading season anime...")
	}

	view, err := gui.g.View(viewName)
	if err != nil {
		return err
	}

	// Clear the view and prepare it for rendering.
	view.Clear()

	for i, s := range tab.Data.Data {
		if i == tab.SelectedIndex {
			fmt.Fprintf(view, "\033[31m➤\033[0m  %s\n", s.Title)
		} else {
			fmt.Fprintf(view, "   %s\n", s.Title)
		}
	}

	return nil
}
