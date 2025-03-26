package gui

import (
	"fmt"
	"slices"

	"github.com/awesome-gocui/gocui"
)

const UNKNOWN_VIEW_ERROR_MSG = "unknown view"

type Views struct {
	Preview *gocui.View // Side Panel

	// Tabs
	Search  *gocui.View
	Seasons *gocui.View
	About   *gocui.View

	// Create three views for each "tab" (tab) with titles.
	// Each title is in a dedicated view just for managing active tab.
	// This approach is necessary because gocui seams not provide an easy or native way
	// to create tabs. This will need a serious refactor one day.
	TSearch  *gocui.View
	TSeasons *gocui.View
	TAbout   *gocui.View
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

func (gui *Gui) getTViews() []ViewMap {
	// ViewMap for View Titles
	return []ViewMap{
		{viewPtr: &gui.Views.TSearch, name: "tsearch"},
		{viewPtr: &gui.Views.TSeasons, name: "tseasons"},
		{viewPtr: &gui.Views.TAbout, name: "tabout"},
	}
}

func (gui *Gui) createAllViews() error {
	frameRunes := []rune{'─', '│', '╭', '╮', '╰', '╯'}

	var err error
	allViews := slices.Concat(gui.getViews(), gui.getTViews())
	for _, v := range allViews {
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

	gui.Views.TSearch.Title = "Search"
	gui.Views.TSearch.Highlight = true

	gui.Views.TSeasons.Title = "Seasons"
	gui.Views.TSeasons.Highlight = true

	gui.Views.TAbout.Title = "About"
	gui.Views.TAbout.Highlight = true

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

func (gui *Gui) applyAlwaysOnTop() error {
	// Always on top for view which handle the title display
	for _, v := range gui.getTViews() {
		_, err := gui.setCurrentTabOnTop(v.name)
		if err != nil {
			return nil
		}
	}

	return nil
}
