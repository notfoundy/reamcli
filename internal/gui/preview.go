package gui

import (
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/darenliang/jikan-go"
)

const NO_DATA = "No data to load..."

func (gui *Gui) renderPreview(viewName string) error {
	tab, err := gui.getCurrentTabOnTop()
	if err != nil {
		return err
	}

	if tab.Data == nil || len(tab.Data.Data) == 0 {
		return gui.renderString(gui.g, viewName, NO_DATA)
	}

	selectedItem := tab.Data.Data[tab.SelectedIndex]

	v, err := gui.g.View(viewName)
	if err != nil {
		return err
	}
	v.Clear()

	gui.renderTitle(v, selectedItem)
	fmt.Fprintln(v)
	fmt.Fprintf(v, "Information:\n")
	gui.renderInfo(v, selectedItem)
	fmt.Fprintln(v)
	gui.renderSynopsis(v, selectedItem)

	return nil
}

func (gui *Gui) renderTitle(v *gocui.View, selectedItem jikan.AnimeBase) {
	fmt.Fprintf(v, "\033[1;32m%s\033[0m\n", selectedItem.Title)

	if len(selectedItem.TitleEnglish) > 0 && selectedItem.Title != selectedItem.TitleEnglish {
		fmt.Fprintf(v, "%s\n", selectedItem.TitleEnglish)
	}
}

func (gui *Gui) renderInfo(v *gocui.View, selectedItem jikan.AnimeBase) {
	fmt.Fprintf(v, "Type: %s\n", selectedItem.Type)
	fmt.Fprintf(v, "Episodes: %d\n", selectedItem.Episodes)
	fmt.Fprintf(v, "Status: %s\n", selectedItem.Status)
	displayMalDateRange(v, "Aried: %s to %s\n", selectedItem.Aired)
	fmt.Fprintf(v, "Source: %s\n", selectedItem.Source)
	displayMalItems(v, "Genre: %s\n", selectedItem.Genres)
	displayMalItems(v, "Theme: %s\n", selectedItem.Themes)
	fmt.Fprintf(v, "Duration: %s\n", selectedItem.Duration)
	fmt.Fprintf(v, "Rating: %s\n", selectedItem.Rating)
}

func (gui *Gui) renderSynopsis(v *gocui.View, selectedItem jikan.AnimeBase) {
	fmt.Fprintf(v, "Synopsis:\n%s\n", selectedItem.Synopsis)
}

func displayMalDateRange(v *gocui.View, m string, date jikan.DateRange) {
	fromDate := time.Date(date.Prop.From.Year, time.Month(date.Prop.From.Month), date.Prop.From.Day, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(date.Prop.To.Year, time.Month(date.Prop.To.Month), date.Prop.To.Day, 0, 0, 0, 0, time.UTC)
	layout := "Jan 2, 2006"
	fmt.Fprintf(v, m, fromDate.Format(layout), toDate.Format(layout))
}

func displayMalItems(v *gocui.View, m string, items []jikan.MalItem) {
	for _, i := range items {
		fmt.Fprint(v, m, i.Name)
	}
}
