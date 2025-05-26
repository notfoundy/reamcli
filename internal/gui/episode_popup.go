package gui

import (
	"github.com/awesome-gocui/gocui"
	"github.com/darenliang/jikan-go"
	"github.com/notfoundy/reamcli/internal/ani"
)

type Episodes struct {
	Data          *jikan.AnimeEpisodes
	Render        func() error
	SelectedIndex int
}

func (gui *Gui) setEpisodesList(malId int) Episodes {
	ep, _ := ani.GetAiredEpisode(malId)
	return Episodes{
		Data: ep,
		Render: func() error {
			return gui.renderEpisodesList("episodes")
		},
		SelectedIndex: 0,
	}
}

func (gui *Gui) createEpisodesPopup(title string, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	gui.g.Update(func(g *gocui.Gui) error {
		err := gui.prepareEpisodesPopup(title)
		if err != nil {
			return err
		}
		gui.Views.Episodes.Editable = false
		return gui.setKeyBindings(g, handleConfirm, handleClose)
	})
	return nil
}

func (gui *Gui) prepareEpisodesPopup(title string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return nil
	}
	// DEGUG: still have eps

	x0, y0, x1, y1 := gui.getEpisodesPopupDimensions(tab.Episodes.Data)
	v, err := gui.g.SetView("episodes", x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	gui.Views.Episodes = v
	v.Title = title
	v.Visible = true
	v.Wrap = true
	v.Clear()

	tab.Episodes.Render()

	if _, err := gui.g.SetCurrentView("episodes"); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) getEpisodesPopupDimensions(episodes *jikan.AnimeEpisodes) (int, int, int, int) {
	width, height := gui.g.Size()
	panelWidth := width / 2
	panelHeight := int(float32(len(episodes.Data)) * 1.25)
	return width/2 - panelWidth/2,
		height/2 - panelHeight/2 - panelHeight%2 - 1,
		width/2 + panelWidth/2,
		height/2 + panelHeight/2
}

func (gui *Gui) setKeyBindings(g *gocui.Gui, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	bindings := []*Binding{
		{
			ViewName: "episodes",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedConfirmationFunction(handleClose),
		},
		{
			ViewName: "episodes",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedConfirmationFunction(handleConfirm),
		},
		{
			ViewName: "episodes",
			Key:      'j',
			Modifier: gocui.ModNone,
			Handler:  gui.nextEpisode,
		},
		{
			ViewName: "episodes",
			Key:      'k',
			Modifier: gocui.ModNone,
			Handler:  gui.PreviousEpisode,
		},
	}
	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}

func (gui *Gui) wrappedConfirmationFunction(function func(*gocui.Gui, *gocui.View) error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		if err := gui.handleEpisodePopupClose(); err != nil {
			return err
		}

		if function != nil {
			if err := function(g, v); err != nil {
				return err
			}
		}

		return nil
	}
}

func (gui *Gui) handleEpisodePopupClose() error {
	gui.g.DeleteKeybindings("episodes")
	return nil
}

func (gui *Gui) nextEpisode(g *gocui.Gui, v *gocui.View) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if tab.Episodes.SelectedIndex < len(tab.Episodes.Data.Data)-1 {
		tab.Episodes.SelectedIndex++
	}
	return tab.Episodes.Render()
}

func (gui *Gui) PreviousEpisode(g *gocui.Gui, v *gocui.View) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if tab.Episodes.SelectedIndex > 0 {
		tab.Episodes.SelectedIndex--
	}
	return tab.Episodes.Render()
}
