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
	var x0, x1 int
	if name == "preview" {
		x0 = 1
		x1 = sideSectionWeight - frameOffset
	} else {
		x0 = sideSectionWeight + frameOffset
		x1 = sideSectionWeight + mainSectionWeight - frameOffset - 1
	}

	return Position{x0, 1, x1, maxY - 2}
}
