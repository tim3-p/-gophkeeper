package main

import (
	"log"
	"os"

	"github.com/tim3-p/gophkeeper/cmd/client/internal/action"
	"github.com/tim3-p/gophkeeper/cmd/client/internal/config"
)

const defaultConfigFile = "gosecret.cfg"

func main() {
	var err error
	configFile, ok := os.LookupEnv("GOSECRET_CFG")
	if !ok {
		configFile = defaultConfigFile
	}

	ok, err = config.CheckFileMode(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("wrong config file mode")
	}

	err = config.ParseConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = config.ParseFlags()
	if err == config.ErrUnknownMode {
		config.Usage("unknown mode is choosen")
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = action.ChooseAct()
	if err != nil {
		log.Fatal(err)
	}
}
