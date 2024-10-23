package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.SetFormatter(&logrus.JSONFormatter{})
	// Set level based on configuration
	Log.SetLevel(logrus.InfoLevel)
}
