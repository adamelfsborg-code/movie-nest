package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adamelfsborg-code/movie-nest/config"
	"github.com/adamelfsborg-code/movie-nest/db"
	"github.com/go-pg/pg/v10"
)

type Server struct {
	router  http.Handler
	config  config.Environments
	datbase pg.DB
}

func New(config config.Environments) *Server {
	server := &Server{
		config: config,
	}

	d := pg.Connect(&pg.Options{
		Addr:     config.DatabaseAddr,
		Database: config.DatabaseName,
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
	})

	server.datbase = *d
	server.loadRoutes()

	return server
}

func (a *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.config.ServerAddr,
		Handler: a.router,
	}

	err := a.datbase.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Failed to connect to repo: %w", err)
	}

	defer func() {
		err := a.datbase.Close()
		if err != nil {
			fmt.Println("Failed to close Repo", err)
		}
	}()

	a.datbase.AddQueryHook(&db.QueryLogger{})

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := a.datbase.Ping(ctx)
				var searchPath string
				_, err = a.datbase.QueryOne(pg.Scan(&searchPath), "SHOW search_path")
				if err != nil {
					fmt.Println("Error getting search path:", err)
					os.Exit(1)
				}

				if err != nil {
					log.Println("Database connection lost:", err)
				}
			}
		}
	}()

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("Failed to start server: %w", err)
		}

		close(ch)
	}()

	fmt.Println("Server started")

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
