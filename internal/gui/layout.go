package gui

import (
	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	gui.setViews(gui.g, gui.getViews())
	gui.setViews(gui.g, gui.getTViews())

	gui.highlighActiveTitleViewTab()
	gui.applyAlwaysOnTop()

	// here is a good place log some stuff
	// if you download humanlog and do tail -f development.log | humanlog
	// this will let you see these branches as prettified json
	// gui.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return nil
}

func (gui *Gui) setViews(g *gocui.Gui, views []ViewMap) error {
	frameOffset := 1
	for _, v := range views {
		pos := gui.getPositionByViewName(v.name, frameOffset)
		_, err := g.SetView(v.name, pos.x0, pos.y0, pos.x1, pos.y1, 0)
		if err != nil {
			return err
		}
	}

	return nil
}
