package main

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/tim3-p/gophkeeper/cmd/server/internal/config"
	"github.com/tim3-p/gophkeeper/internal/server"
)

func InitConfig() error {
	err := env.Parse(&config.EnvConfig)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = server.StartServer(
		config.Cfg.ServerPort,
		config.Cfg.StoreFile,
		config.Cfg.ServerKey,
		config.Cfg.ServerCRT,
	)
	if err != nil {
		log.Fatal(err)
	}
}
