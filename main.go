package main

import (
	"fmt"

	"go-redis/config"
	"go-redis/lib/file"
	"go-redis/lib/logger"
	"go-redis/resp/handler"
	"go-redis/tcp"
)

const configFile string = "etc/config.yaml"

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "tcp-server",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if !file.CheckNotExist(configFile) {
		config.SetupConfig(configFile, ".")
	}

	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)},
		handler.MakeRespHandler())

	if err != nil {
		logger.Fatal(err)
	}
}
