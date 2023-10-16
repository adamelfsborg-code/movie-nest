package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/server"
)

func main() {
	env, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	server := server.New(*env)
	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
