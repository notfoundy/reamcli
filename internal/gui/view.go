package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	// Side Panel
	Preview *gocui.View

	// Tabs
	Search  *gocui.View
	Seasons *gocui.View
	About   *gocui.View

	// Modal
	Episodes *gocui.View
	Filter   *gocui.View
}

type ViewMap struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) getViews() []ViewMap {
	return []ViewMap{
		{viewPtr: &gui.Views.Preview, name: "preview"},

		{viewPtr: &gui.Views.Search, name: "search"},
		{viewPtr: &gui.Views.Seasons, name: "seasons"},
		{viewPtr: &gui.Views.About, name: "about"},
	}
}

func (gui *Gui) createAllViews() error {
	frameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}

	var err error
	for _, v := range gui.getViews() {
		if v.viewPtr == nil {
			return fmt.Errorf("viewPtr is nil for view: %s", v.name)
		}

		*v.viewPtr, err = gui.prepareView(v.name)
		if err != nil && err.Error() != UNKNOWN_VIEW_ERROR_MSG {
			return err
		}

		(*v.viewPtr).FrameRunes = frameRunes
		(*v.viewPtr).FgColor = gocui.ColorDefault
	}

	gui.Views.Preview.Title = "Preview"
	gui.Views.Preview.Wrap = true

	gui.Views.Search.Title = "Search"
	gui.Views.Search.Highlight = true
	gui.Views.Search.Autoscroll = true
	gui.Views.Search.Wrap = true

	gui.Views.Seasons.Title = "Seasons"
	gui.Views.Seasons.Highlight = true
	gui.Views.Seasons.Autoscroll = false
	gui.Views.Seasons.Wrap = false

	gui.Views.About.Title = "About"
	gui.Views.About.Highlight = true
	gui.Views.About.Autoscroll = true
	gui.Views.About.Wrap = true

	defaultViewOnTop, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}
	gui.g.SetViewOnTop(defaultViewOnTop.Key)

	return nil
}

func (gui *Gui) prepareView(viewName string) (*gocui.View, error) {
	// arbitrarily giving the view enough size so that we don't get an error, but
	// it's expected that the view will be given the correct size before being shown
	return gui.g.SetView(viewName, 0, 0, 10, 10, 0)
}
