package logs

import (
	"github.com/sirupsen/logrus"
	"io"
	"nh-downloader/internal/config"
	"os"
	"path"
)

var logger *logrus.Logger

func Init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	// todo config
	logger.SetLevel(logrus.DebugLevel)

	// set log file
	// todo config
	file, err := os.OpenFile(path.Join(config.LogsPath(), "nhdl.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		multiWriter := io.MultiWriter(file, os.Stdout)
		logger.SetOutput(multiWriter)
	} else {
		logrus.Fatal("Failed to open log file:", err)
	}
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}
