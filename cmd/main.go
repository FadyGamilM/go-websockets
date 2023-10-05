package main

import (
	"fmt"
	"log"
	"os"

	"github.com/FadyGamilM/go-websockets/configs"
	paseto "github.com/FadyGamilM/go-websockets/internal/business/auth/paseto"
	userService "github.com/FadyGamilM/go-websockets/internal/business/user"
	"github.com/FadyGamilM/go-websockets/internal/database/postgres"
	userRepo "github.com/FadyGamilM/go-websockets/internal/repository"
	"github.com/FadyGamilM/go-websockets/internal/transport"
	"github.com/FadyGamilM/go-websockets/internal/transport/handlers"
)

func main() {
	//! 1. connect to postgres
	db, err := postgres.SetupPostgresConnection()
	if err != nil {
		log.Printf("error trying to connect to database : %v \n", err)
		os.Exit(1)
	}

	// instantiate repos
	userRepository := userRepo.New(db)

	// instantiate the concrete impl of the infrastructure dependnecies (such as paseto token auth mechanism)
	pasetoConfigs, err := configs.LoadPasetoTokenConfig("./configs")
	if err != nil {
		fmt.Println("error trying to load config variables", err)
		os.Exit(1)
	}
	pasetoTokenAuth, err := paseto.NewPaseto(pasetoConfigs.Paseto.SymmetricKey)
	if err != nil {
		log.Printf("error trying to create paseto token auth imp | %v \n", err)
		os.Exit(1)
	}

	// instantiate services
	userService := userService.NewUserService(&userService.UserServiceConfig{UserRepo: userRepository, TokenAuth: pasetoTokenAuth})

	// instantiate router
	router := transport.CreateRouter()

	// instantiate handler
	handlers.NewUserHandler(&handlers.UserHandlerConfig{
		R:           router,
		UserService: userService,
	})

	server := transport.CreateServer(router)
	server.Run()
}
