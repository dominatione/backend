package main

import (
	"context"
	"github.com/dominati-one/backend/internal/app/backend"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Stamp})

	ctx := context.Background()

	parameters := backend.AppParameters{
		GrpcApiListenAddress: "127.0.0.1",
		GrpcApiListenPort:    3009,
	}

	app := backend.NewApp(parameters)

	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}

	select {}

	os.Exit(0)
}
