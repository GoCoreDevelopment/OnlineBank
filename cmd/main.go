package main

import (
	"api/internal/config"
	"api/internal/db"
	"api/internal/handlers"
	"api/internal/repository"
	"api/internal/server"
	"api/internal/servises/jwtservice"
	"api/internal/servises/userservice"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	DB, err := db.InitDB(*cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseDB(DB)

	repo := repository.NewRepository(DB)
	userService := userservice.NewUserService(repo)
	jwtService := jwtservice.NewJWTService(cfg)
	userHandlers := handlers.NewUserHandlers(userService, jwtService)

	e := echo.New()
	userHandlers.InitHandlers(e)
	
	if err := server.StartServer(e, cfg); err != nil {
		log.Fatal(err)
	}
}
