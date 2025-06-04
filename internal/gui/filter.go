package gui

import (
	"fmt"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/notfoundy/reamcli/internal/ani"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

const VIEW_NAME = "filter"

func (gui *Gui) openFilter(g *gocui.Gui, v *gocui.View) error {
	return gui.createFilterModal()
}

func (gui *Gui) createFilterModal() error {
	gui.g.Update(func(g *gocui.Gui) error {
		tab, err := gui.getCurrentTabOnTop()
		if err != nil {
			return err
		}

		pos := gui.getPositionByViewName(VIEW_NAME, 1)
		title := fmt.Sprintf("%s (type to search)", VIEW_NAME)
		filterView, err := gui.prepareFilterModal(title, pos)
		if err != nil {
			return err
		}
		filterView.Editable = true
		filterView.Editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			gui.filterOnBufferChangeEditor(v, key, ch, mod)
		})

		if cached, ok := gui.CacheFilter[tab.Key]; ok {
			filterView.Clear()
			fmt.Fprint(filterView, cached)
		}

		if _, err := gui.g.SetCurrentView(VIEW_NAME); err != nil {
			return err
		}

		// NOTE: We use the same action to confirm or close the filter modal
		return gui.setFilterModalKeyBindings(gui.g, gui.closeFilter(tab.Key), gui.closeFilter(tab.Key))
	})
	return nil
}

func (gui *Gui) closeFilter(fromViewName string) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		g.DeleteView(v.Name())
		_, err := gui.g.SetCurrentView(fromViewName)
		return err
	}
}

func (gui *Gui) filterOnBufferChangeEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) error {
	gocui.DefaultEditor.Edit(v, key, ch, mod)
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	gui.g.Update(func(g *gocui.Gui) error {
		searchContent := v.Buffer()
		gui.CacheFilter[tab.Key] = searchContent
		if searchContent == "" {
			tab.Data = tab.OriginalData
			return nil
		}

		var combinedTitles []string
		for i, a := range tab.OriginalData {
			combined := strconv.Itoa(i+1) + " " + a.AnimeDetails.Title + " " + a.AnimeDetails.TitleEnglish
			combinedTitles = append(combinedTitles, combined)
		}

		matches := fuzzy.Find(searchContent, combinedTitles)
		filtered := lo.Map(matches, func(m fuzzy.Match, _ int) *ani.Anime {
			return tab.OriginalData[m.Index]
		})

		tab.Data = filtered
		tab.Render()
		return nil
	})

	return nil
}

func (gui *Gui) prepareFilterModal(title string, pos Position) (*gocui.View, error) {
	v, err := gui.g.SetView(VIEW_NAME, pos.x0, pos.y0, pos.x1, pos.y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return nil, err
	}

	gui.Views.Episodes = v
	v.Title = title
	v.Visible = true
	v.Clear()

	return v, nil
}

func (gui *Gui) setFilterModalKeyBindings(g *gocui.Gui, handleConfirm, handleClose func(*gocui.Gui, *gocui.View) error) error {
	bindings := []*Binding{
		{
			ViewName: VIEW_NAME,
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedConfirmationFunction(handleClose),
		},
		{
			ViewName: VIEW_NAME,
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedConfirmationFunction(handleClose),
		},
		{
			ViewName: VIEW_NAME,
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  gui.wrappedConfirmationFunction(handleConfirm),
		},
	}
	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}
	return nil
}
