package main

import (
	"errors"
	"github.com/kercylan98/go-log/log"
	"os"
)

func main() {
	var loggerA, loggerB log.Logger
	var builder = log.GetBuilder()

	loggerA = log.GetDefault()

	loggerA = builder.Build()
	loggerA = builder.Develop()
	loggerA = builder.Test()
	loggerA = builder.Production()
	loggerA = builder.Silent()
	loggerA = builder.DevelopOnGoland()

	loggerB = builder.FromConfiguration(
		log.GetConfigBuilder().Build().
			WithWriter(os.Stdout).
			WithLeveler(log.LevelInfo),
	)

	loggerB = builder.FromConfigurators(log.LoggerConfiguratorFn(func(config log.LoggerConfiguration) {
		config.WithEnableColor(true).WithLevelStr(log.LevelDebug, "debug")
	}))

	loggerC := builder.Multi(loggerA, loggerB)

	loggerC.Error("Example", "err", errors.New("example error"))
	loggerC.Error("Example", log.Err(errors.New("example error")))
}
