package app

import (
	"strconv"

	"github.com/notfoundy/reamcli/internal/gui"
	"github.com/notfoundy/reamcli/internal/log"
	"github.com/notfoundy/reamcli/internal/mal"
	"github.com/notfoundy/reamcli/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type App struct {
	Gui       *gui.Gui
	Log       *logrus.Logger
	MalClient *mal.Client
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

	clientId := viper.GetString("mal_client_id")
	callbackUrl := "http://localhost:" + strconv.Itoa(utils.GetCallbackPort()) + "/callback"
	client := mal.NewClient(
		clientId,
		callbackUrl,
	)

	if client.AccessToken == "" {
		if err := client.StartOAuth(); err != nil {
			return nil, err
		}
	}

	app.MalClient = client

	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}
