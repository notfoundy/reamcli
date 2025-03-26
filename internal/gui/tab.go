package gui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type Tab struct {
	Index    int
	Key      string
	Title    string
	Render   string
	IsActive bool
}

func (gui *Gui) getSearchTab() Tab {
	return Tab{
		Index:    0,
		Key:      "search",
		Title:    "Search",
		Render:   gui.renderTest("SEARCH"),
		IsActive: true,
	}
}

func (gui *Gui) getSeasonsTab() Tab {
	return Tab{
		Index:    1,
		Key:      "seasons",
		Title:    "Seasons",
		Render:   gui.renderTest("SEASONS"),
		IsActive: false,
	}
}

func (gui *Gui) getAboutTab() Tab {
	return Tab{
		Index:    2,
		Key:      "about",
		Title:    "About",
		Render:   gui.renderTest("ABOUT"),
		IsActive: false,
	}
}

func (gui *Gui) renderTest(str string) string {
	return str
}

func (gui *Gui) getCurrentTabOnTop() (*Tab, error) {
	tabs := []*Tab{&gui.Tabs.Search, &gui.Tabs.Seasons, &gui.Tabs.About}
	for _, t := range tabs {
		if t.IsActive {
			return t, nil
		}
	}
	return nil, fmt.Errorf("no tab found")
}

func (gui *Gui) setCurrentTabOnTop(name string) (*gocui.View, error) {
	if _, err := gui.g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return gui.g.SetViewOnTop(name)
}

func (gui *Gui) nextTab(g *gocui.Gui, v *gocui.View) error {
	var nextTabIndex int
	var nextTabName string
	tabs := []*Tab{&gui.Tabs.Search, &gui.Tabs.Seasons, &gui.Tabs.About}
	for _, t := range tabs {
		if t.IsActive {
			nextTabIndex = (t.Index + 1) % len(tabs)
			nextTabName = tabs[nextTabIndex].Key
			t.IsActive = false
			break
		}
	}

	tabs[nextTabIndex].IsActive = true
	if _, err := gui.setCurrentTabOnTop(nextTabName); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) previousTab(g *gocui.Gui, v *gocui.View) error {
	var nextTabIndex int
	var nextTabName string
	tabs := []*Tab{&gui.Tabs.Search, &gui.Tabs.Seasons, &gui.Tabs.About}
	for _, t := range tabs {
		if t.IsActive {
			if t.Index == 0 {
				nextTabIndex = len(tabs) - 1
			} else {
				nextTabIndex = (t.Index - 1) % len(tabs)
			}
			nextTabName = tabs[nextTabIndex].Key
			t.IsActive = false
			break
		}
	}

	tabs[nextTabIndex].IsActive = true
	if _, err := gui.setCurrentTabOnTop(nextTabName); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) highlighActiveTitleViewTab() error {
	tabs := []*Tab{&gui.Tabs.Search, &gui.Tabs.Seasons, &gui.Tabs.About}
	for _, t := range tabs {
		if t.IsActive {
			v, err := gui.g.View("t" + t.Key)
			if err != nil {
				return err
			}
			v.Highlight = true
			v.TitleColor = gocui.ColorGreen
		} else {
			v, err := gui.g.View("t" + t.Key)
			if err != nil {
				return err
			}
			v.Highlight = false
			v.TitleColor = gocui.ColorDefault
		}
	}

	return nil
}
