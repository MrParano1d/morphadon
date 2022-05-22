package engine

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	if viper.GetString("MARLA_ENV") == "production" {
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{
			"./logs/marla.log",
		}
		zLogger, err := cfg.Build()
		if err != nil {
			panic(err)
		}

		Logger = zLogger
	} else {
		zLogger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		Logger = zLogger
	}
}
