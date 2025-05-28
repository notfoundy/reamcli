package gui

import (
	"fmt"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/notfoundy/reamcli/internal/ani"
	"github.com/notfoundy/reamcli/internal/player"
)

type Tab struct {
	Index         int
	Key           string
	Render        func() error
	Data          []*ani.Anime
	SelectedIndex int
	IsActive      bool
	Episodes      Episodes
}

func (gui *Gui) setSearchTab() Tab {
	return Tab{
		Index: 0,
		Key:   "search",
		Render: func() error {
			return gui.renderString(gui.g, "search", "Hello world from search")
		},
		IsActive: true,
	}
}

func (gui *Gui) setSeasonsTab() Tab {
	// data, err := ani.GetSeasonNow()
	seasonName := "Spring" // ou "Spring", "Fall", "Winter"
	year := 2025
	data, err := ani.GetSeasonAnimes(seasonName, year)
	if err != nil {
		gui.Log.Error(err)
		return Tab{
			Index: 1,
			Key:   "seasons",
			Render: func() error {
				return gui.renderString(gui.g, "seasons", "Error can't load animes")
			},
			IsActive: true,
		}
	}
	gui.Log.Debug(fmt.Sprintf("nombre animes : %d", len(data)))

	tab := Tab{
		Index: 1,
		Key:   "seasons",
		Data:  data,
		Render: func() error {
			return gui.renderTabList("seasons")
		},
		IsActive: false,
	}

	return tab
}

func (gui *Gui) setAboutTab() Tab {
	return Tab{
		Index: 2,
		Key:   "about",
		Render: func() error {
			return gui.renderString(gui.g, "about", "Hello world from about")
		},
		IsActive: false,
	}
}

func (gui *Gui) allTabs() []*Tab {
	return []*Tab{&gui.Tabs.Search, &gui.Tabs.Seasons, &gui.Tabs.About}
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

func (gui *Gui) nextItem(g *gocui.Gui, v *gocui.View) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if tab.SelectedIndex < len(tab.Data)-1 {
		tab.SelectedIndex++
	}
	return gui.renderTabList("seasons")
}

func (gui *Gui) previousItem(g *gocui.Gui, v *gocui.View) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if tab.SelectedIndex > 0 {
		tab.SelectedIndex--
	}
	return gui.renderTabList("seasons")
}

func (gui *Gui) handleEnterTab(g *gocui.Gui, v *gocui.View) error {
	fromView := v
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	selectedItem := tab.Data[tab.SelectedIndex]
	malId, err := strconv.Atoi(selectedItem.MalId)
	if err != nil {
		return err
	}
	tab.Episodes = gui.setEpisodesList(malId)

	// BUG: need to lock other views when we play smething
	return gui.createEpisodesPopup("List of episodes", func(g *gocui.Gui, v *gocui.View) error {
		msg := fmt.Sprintf("Playing episode %d : %s", tab.Episodes.SelectedIndex+1, tab.Episodes.Data.Data[tab.Episodes.SelectedIndex].Title)
		gui.renderString(g, v.Name(), msg)
		url, err := ani.GetEpisodeUrl()
		if err != nil {
			return err
		}
		go func() {
			player.Play(url)
			tab.Episodes.Render()
		}()
		return nil
	}, func(g *gocui.Gui, v *gocui.View) error {
		g.DeleteView(v.Name())
		gui.g.SetCurrentView(fromView.Name())

		return nil
	})
}
