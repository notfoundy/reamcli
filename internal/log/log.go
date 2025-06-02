package log

import (
	"fmt"
	"os"

	"github.com/notfoundy/reamcli/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)

	env := viper.GetString("env")
	var logPath string
	if env == "development" {
		logPath = "log/development.log"
	} else {
		logPath, _ = utils.GetLogPath()
	}

	// highly recommended: tail -f development.log | humanlog
	// https://github.com/aybabtme/humanlog
	log.Formatter = &logrus.JSONFormatter{}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Println("Unable to log to file")
		os.Exit(1)
	}
	log.SetOutput(file)

	return log
}
