package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	Log.SetLevel(logrus.InfoLevel)

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}
