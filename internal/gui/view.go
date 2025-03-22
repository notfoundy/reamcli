package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	// Side Panel
	Streaming *gocui.View
	Preview   *gocui.View

	// Main Panel
	Main *gocui.View
}

type ViewMap struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) getViews() []ViewMap {
	return []ViewMap{
		{viewPtr: &gui.Views.Preview, name: "preview"},
		{viewPtr: &gui.Views.Main, name: "main"},
	}
}

func (gui *Gui) createAllViews() error {
	frameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}

	var err error
	for _, mapping := range gui.getViews() {
		if mapping.viewPtr == nil {
			return fmt.Errorf("viewPtr is nil for view: %s", mapping.name)
		}

		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		(*mapping.viewPtr).FrameRunes = frameRunes
		(*mapping.viewPtr).FgColor = gocui.ColorDefault
	}

	gui.Views.Preview.Title = "Preview"

	return nil
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10, 0)
}
