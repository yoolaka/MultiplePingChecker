package main

import (
	//"bufio"
	"github.com/MultiplePingChecker/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"net/http"
	//"os/exec"
	//"strconv"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/ping", api.CreatePing)
	e.GET("/ping", api.GetPing)
	e.GET("/ping/:hostname", api.GetPingStatus)
	e.DELETE("/ping/:hostname", api.DeletePing)

	// Start server
	e.Logger.Fatal(e.Start(":9335"))
}
