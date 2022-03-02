package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rodinv/errors"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/rodinv/geodata/internal/geo"
	"github.com/rodinv/geodata/internal/geo/repository"
	"github.com/rodinv/geodata/internal/pkg/config"
)

// @title Swagger Geodata API
// @version 1.0
// @description This is geodata server.
// @BasePath /v1
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var configPath string

	flag.StringVar(&configPath, "config", "", "path to config")
	flag.Parse()

	if configPath == "" {
		return errors.New("specify the 'config' parameter")
	}

	// init config
	cfg, err := config.Get(configPath)
	if err != nil {
		return errors.Wrapf(err, "loading config")
	}

	// init database
	db, err := repository.Load(cfg.Db.Path)
	if err != nil {
		return errors.Wrapf(err, "loading data file %s", cfg.Db.Path)
	}

	// init handler
	handler := geo.New(db)

	// init Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API V1
	v1 := e.Group("/v1")

	// register handlers
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	v1.GET("/ip/location", handler.GetLocationByIP)
	v1.GET("/city/locations", handler.GetLocations)

	// Start server
	s := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Http.Port),
		ReadTimeout: cfg.Http.Timeout,
	}

	e.Logger.Fatal(e.StartServer(s))

	return nil
}
