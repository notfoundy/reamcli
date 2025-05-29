package gui

import (
	"fmt"
	"strings"
	"time"

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

func (gui *Gui) renderTabList(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if len(tab.Data) == 0 {
		return gui.renderString(gui.g, viewName, "Loading season anime...")
	}

	view, err := gui.g.View(viewName)
	if err != nil {
		return err
	}
	view.Clear()
	view.SelFgColor = gocui.ColorRed | gocui.AttrBold
	view.SelBgColor = gocui.ColorBlack

	for i, s := range tab.Data {
		if i == tab.SelectedIndex {
			fmt.Fprintf(view, "\033[31m➤\033[0m  %s\n", s.AnimeDetails.Title)
		} else {
			fmt.Fprintf(view, "   %s\n", s.AnimeDetails.Title)
		}
	}
	view.SetCursor(0, tab.SelectedIndex)
	err = gui.scrollToSelectedItem(view, tab.SelectedIndex)
	if err != nil {
		return err
	}

	return nil
}

func (gui *Gui) scrollToSelectedItem(v *gocui.View, selectedIndex int) error {
	ox, oy := v.Origin()
	_, y := v.Size()

	if selectedIndex < oy {
		return v.SetOrigin(ox, selectedIndex)
	}

	if selectedIndex >= oy+y {
		newOrigin := max(selectedIndex-y+1, 0)
		return v.SetOrigin(ox, newOrigin)
	}
	return nil
}

func (gui *Gui) renderEpisodesList(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}
	episodes := tab.Episodes
	var builder strings.Builder

	for i, item := range episodes.Data {
		indexStr := fmt.Sprintf("%d - ", i+1)
		airedSince := formatAiredDate(item.AiredDate)

		if i == episodes.SelectedIndex {
			builder.WriteString("\033[31m➤\033[0m ")
		} else {
			builder.WriteString("   ")
		}

		line := fmt.Sprintf("%s - %s  (%s)\n", indexStr, item.Title, airedSince)
		builder.WriteString(line)
	}

	return gui.renderString(gui.g, viewName, builder.String())
}

func formatAiredDate(t time.Time) string {
	now := time.Now()
	if t.After(now) {
		return ""
	}
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return "today"
	}
	days := int(now.Sub(t).Hours() / 24)

	return fmt.Sprintf("%d days ago", days)
}
