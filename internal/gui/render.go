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

func (gui *Gui) renderEpisodesList(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}
	episodes := tab.Episodes
	var builder strings.Builder

	for i, item := range episodes.Data.Data {
		indexStr := fmt.Sprintf("%d - ", i+1)
		airedSince := formatAiredDate(item.Aired)

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
