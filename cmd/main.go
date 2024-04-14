package main

import (
	"banner_service/cmd/server"
	"banner_service/internal/config"
	"banner_service/internal/logger"
	"fmt"
	"os"
)

const configPath = "configs/local.yaml"

func main() {
	cfg, err := config.New(configPath)
	exitOnError(err)

	app, err := server.New(cfg, logger.New(os.Stdout))
	exitOnError(err)

	err = app.Run()
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
