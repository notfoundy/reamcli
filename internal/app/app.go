package app

import (
	"github.com/notfoundy/reamcli/internal/gui"
	"github.com/notfoundy/reamcli/internal/log"
	"github.com/sirupsen/logrus"
)

type App struct {
	Gui       *gui.Gui
	Log       *logrus.Logger
	ErrorChan chan error
}

func NewApp() (*App, error) {
	app := &App{
		ErrorChan: make(chan error),
	}

	var err error
	app.Log = log.NewLogger()
	app.Gui, err = gui.NewGui(app.Log, app.ErrorChan)
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}
