package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/usace-nsi/consequences-server/handlers"
)

const apiprefix = "/consequences"

func main() {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.RequestLogger())
	handler := handlers.Handler{}

	e.GET(apiprefix+"/version", handler.Version)
	e.POST(apiprefix+"/compute", handler.Compute)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
