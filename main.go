package main

import (
	"github.com/MultiplePingChecker/api"
	"github.com/MultiplePingChecker/temp_db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	a := &api.Api{temp_db.InitDB(100, "")}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/ping", a.CreatePing)
	e.GET("/ping", a.GetPing)
	e.GET("/ping/:hostname", a.GetPingStatus)
	e.DELETE("/ping/:hostname", a.DeletePing)

	// Start server
	e.Logger.Fatal(e.Start("localhost:9335"))
}
