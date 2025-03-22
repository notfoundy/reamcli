package gui

import (
	"github.com/awesome-gocui/gocui"
)

func (gui *Gui) layout(g *gocui.Gui) error {
	g.Highlight = true

	for _, v := range gui.getViews() {
		frameOffset := 1
		pos := gui.getPositionByViewName(v.name, frameOffset)
		g.SetView(v.name, pos.x0, pos.y0, pos.x1, pos.y1, 0)
	}

	// here is a good place log some stuff
	// if you download humanlog and do tail -f development.log | humanlog
	// this will let you see these branches as prettified json
	// gui.Log.Info(utils.AsJson(gui.State.Branches[0:4]))
	return nil
}
