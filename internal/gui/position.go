package gui

type Position struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

func (gui *Gui) getMidSectionWeights() (int, int) {
	maxX, _ := gui.g.Size()
	sidePanelWidthRatio := 0.33
	sideSectionWeight := int(float64(maxX) * sidePanelWidthRatio)
	mainSectionWeight := maxX - sideSectionWeight

	return sideSectionWeight, mainSectionWeight
}

func (gui *Gui) getPositionByViewName(name string, frameOffset int) Position {
	_, maxY := gui.g.Size()
	sideSectionWeight, mainSectionWeight := gui.getMidSectionWeights()
	var x0, x1, y0, y1 int
	switch name {
	case "preview":
		x0 = 1
		x1 = sideSectionWeight - frameOffset
		y0 = 1
		y1 = maxY - 2
	default: // we assume the default case only contain tabs for the main panel
		x0 = sideSectionWeight + frameOffset
		x1 = sideSectionWeight + mainSectionWeight - frameOffset - 1
		y0 = 1
		y1 = maxY - 2
	}

	return Position{x0, y0, x1, y1}
}
