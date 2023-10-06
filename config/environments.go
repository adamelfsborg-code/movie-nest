package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Environments struct {
	ServerAddr       string
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	MovieDBApiKey    string
	MovieDBAuthToken string
}

var Env *Environments

func New() (*Environments, error) {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}

	serverAddr, exists := os.LookupEnv("SERVER_ADDR")
	if exists == false {
		return nil, fmt.Errorf("SERVER_ADDR not found.")
	}

	databaseAddr, exists := os.LookupEnv("DATABASE_ADDR")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_ADDR not found.")
	}

	databaseUser, exists := os.LookupEnv("DATABASE_USER")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_USER not found.")
	}

	databasePassword, exists := os.LookupEnv("DATABASE_PASSWORD")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_PASSWORD not found.")
	}

	databaseName, exists := os.LookupEnv("DATABASE_NAME")
	if exists == false {
		return nil, fmt.Errorf("DATABASE_NAME not found.")
	}

	movieDBApiKey, exists := os.LookupEnv("MOVIEDB_API_KEY")
	if exists == false {
		return nil, fmt.Errorf("MOVIEDB_API_KEY not found")
	}

	movieDBAuthToken, exists := os.LookupEnv("MOVIEDB_AUTH_TOKEN")
	if exists == false {
		return nil, fmt.Errorf("MOVIEDB_AUTH_TOKEN not found")
	}

	return &Environments{
		ServerAddr:       serverAddr,
		DatabaseAddr:     databaseAddr,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseName:     databaseName,
		MovieDBApiKey:    movieDBApiKey,
		MovieDBAuthToken: movieDBAuthToken,
	}, nil
}
