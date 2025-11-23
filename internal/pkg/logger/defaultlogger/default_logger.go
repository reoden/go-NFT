package defaultLogger

import (
	"os"

	"github.com/reoden/go-echo-template/pkg/constants"
	"github.com/reoden/go-echo-template/pkg/logger"
	"github.com/reoden/go-echo-template/pkg/logger/config"
	"github.com/reoden/go-echo-template/pkg/logger/logrous"
	"github.com/reoden/go-echo-template/pkg/logger/models"
	"github.com/reoden/go-echo-template/pkg/logger/zap"
)

var l logger.Logger

func initLogger() {
	logType := os.Getenv("LogConfig_LogType")

	switch logType {
	case "Zap", "":
		l = zap.NewZapLogger(
			&config.LogOptions{LogType: models.Zap, CallerEnabled: false},
			constants.Dev,
		)
		break
	case "Logrus":
		l = logrous.NewLogrusLogger(
			&config.LogOptions{LogType: models.Logrus, CallerEnabled: false},
			constants.Dev,
		)
		break
	default:
	}
}

func GetLogger() logger.Logger {
	if l == nil {
		initLogger()
	}

	return l
}
