package gui

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/sirupsen/logrus"
)

type Gui struct {
	g         *gocui.Gui
	Views     Views
	Log       *logrus.Logger
	Tabs      Tabs
	ErrorChan chan error
}

type Tabs struct {
	Search  Tab
	Seasons Tab
	About   Tab
}

func NewGui(log *logrus.Logger, errorChan chan error) (*Gui, error) {
	gui := &Gui{
		Log:       log,
		ErrorChan: errorChan,
	}

	return gui, nil
}

func (gui *Gui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	gui.g = g

	g.SetManagerFunc(gocui.ManagerFunc(gui.layout))

	gui.setTabs()

	if err := gui.createAllViews(); err != nil {
		return err
	}

	if err = gui.keybindings(g); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

	return nil
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gui *Gui) setTabs() {
	gui.Tabs = Tabs{
		Search:  gui.getSearchTab(),
		Seasons: gui.getSeasonsTab(),
		About:   gui.getAboutTab(),
	}
}
