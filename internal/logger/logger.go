package logger

import (
	"go.uber.org/zap"
)

func New() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err.Error())
	}
	sugar := logger.Sugar()
	return sugar
}
