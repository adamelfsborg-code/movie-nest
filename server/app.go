package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/db"
)

type Server struct {
	router http.Handler
	config config.Environments
}

func New(config config.Environments) *Server {
	server := &Server{
		config: config,
	}

	server.loadRoutes()

	return server
}

func (a *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.config.ServerAddr,
		Handler: a.router,
	}

	ch := make(chan error, 1)

	err := db.Store.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Store.Close()
		if err != nil {
			fmt.Println("Failed to close Database", err)
		}
	}()

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("Failed to start server: %w", err)
		}

		close(ch)
	}()

	fmt.Println("Server started")

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
