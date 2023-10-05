package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/adamelfsborg-code/movie-nest/server"
)

func main() {
	env, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	store := db.New(
		env.DatabaseAddr,
		env.DatabaseName,
		env.DatabaseUser,
		env.DatabasePassword,
	)
	store.CreateDatabase()

	db.Store.Exec("SET search_path TO movie_nest")

	db.Store.AddQueryHook(&db.QueryLogger{})

	server := server.New(*env)
	err = server.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
