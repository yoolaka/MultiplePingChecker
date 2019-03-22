package MultiplePingChecker

import (
	"bufio"
	"github.com/MultiplePingChecker/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os/exec"
	"strconv"
)

type (
	ping_host struct {
		HostName string `json:"hostname"`
		Count    int    `json:"count"`
	}
	ping_host_attr struct {
		Channel chan string
		Buffer  string
	}
)

var (
	hosts         = map[string]*ping_host{}
	host_attrs    = map[string]*ping_host_attr{}
	ping_cmds     = map[string]*exec.Cmd{}
	index         = 1
	default_count = 100
	default_name  = ""
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/ping", createPing)
	e.GET("/ping", getPing)
	e.GET("/ping/:hostname", getPingStatus)
	e.DELETE("/ping/:hostname", deletePing)

	// Start server
	e.Logger.Fatal(e.Start(":9335"))
}
