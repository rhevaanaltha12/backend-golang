package helper

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Logger(logType string, message string) {

	time := time.Now().Format("02-01-2006")

	file, _ := os.OpenFile("logs/go-"+time+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		},
	}

	logger.SetOutput(io.MultiWriter(file, os.Stdout))

	switch logType {
	case "info":
		logger.Info(message + "\n")
	case "fail":
		logger.Error(message + "\n")
	}
}
