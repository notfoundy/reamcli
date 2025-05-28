package gui

import (
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/notfoundy/reamcli/internal/ani"
)

const NO_DATA = "No data to load..."

func (gui *Gui) renderPreview(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if len(tab.Data) == 0 {
		return gui.renderString(gui.g, viewName, NO_DATA)
	}

	selectedItem := tab.Data[tab.SelectedIndex]

	v, err := gui.g.View(viewName)
	if err != nil {
		return err
	}
	v.Clear()

	gui.renderTitle(v, selectedItem.AnimeDetails)
	fmt.Fprintln(v)
	fmt.Fprintf(v, "Information:\n")
	gui.renderInfo(v, selectedItem.AnimeDetails)
	fmt.Fprintln(v)
	gui.renderSynopsis(v, selectedItem.AnimeDetails)

	return nil
}

func (gui *Gui) renderTitle(v *gocui.View, a ani.AnimeDetails) {
	fmt.Fprintf(v, "\033[1;32m%s\033[0m\n", a.Title)

	if len(a.TitleEnglish) > 0 && a.Title != a.TitleEnglish {
		fmt.Fprintf(v, "%s\n", a.TitleEnglish)
	}
}

func (gui *Gui) renderInfo(v *gocui.View, a ani.AnimeDetails) {
	fmt.Fprintf(v, "Type: %s\n", a.Type)
	fmt.Fprintf(v, "Episodes: %d\n", a.NumberEpisodes)
	fmt.Fprintf(v, "Status: %s\n", a.Status)
	displayDateRange(v, "Aried: %s to %s\n", a.AiredFrom, a.AiredTo)
	fmt.Fprintf(v, "Source: %s\n", a.Source)
	displayItems(v, "Genre: %s\n", a.Genres)
	displayItems(v, "Theme: %s\n", a.Themes)
	fmt.Fprintf(v, "Rating: %s\n", a.Rating)
}

func (gui *Gui) renderSynopsis(v *gocui.View, a ani.AnimeDetails) {
	fmt.Fprintf(v, "Synopsis:\n%s\n", a.Synopsis)
}

func displayDateRange(v *gocui.View, m string, from time.Time, to time.Time) {
	layout := "Jan 2, 2006"
	fmt.Fprintf(v, m, from.Format(layout), to.Format(layout))
}

func displayItems(v *gocui.View, m string, items []string) {
	for _, i := range items {
		fmt.Fprint(v, m, i)
	}
}
