package main

import (
	//"bytes"
	"bufio"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os/exec"
	"strconv"
	//	"golang.org/x/net/icmp"
	//	"golang.org/x/net/ipv4"
	"net/http"
)

type (
	ping_host struct {
		HostName string `json:"hostname"`
		Count    int    `json:"count"`
		Channel  chan string
		Buffer   string
	}
)

var (
	hosts = map[string]*ping_host{}
	//buffers       = map[string]string{}
	//channels      = map[string](chan string){}
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
	e.GET("/ping/:hostname", getStatus)
	e.DELETE("/ping/:hostname", deletePing)

	// Start server
	e.Logger.Fatal(e.Start(":9335"))
}
func createPing(c echo.Context) error {
	host_name := c.FormValue("server")
	count_string := c.FormValue("count")
	count, err := strconv.Atoi(count_string)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invaild count")
	}

	var buffer string

	channel := make(chan string)
	host := &ping_host{
		host_name,
		count,
		channel,
		buffer,
	}

	hosts[host.HostName] = host

	go executePing(c, host_name, count_string, channel)

	return c.String(http.StatusCreated, "")
}

func executePing(c echo.Context, host_name string, count_string string, channel chan<- string) {

	ping_cmd := exec.Command("ping", "-c", count_string, host_name)
	stdout, err := ping_cmd.StdoutPipe()
	if err != nil {
		c.Logger().Error(err)
	}
	scanner := bufio.NewReader(stdout)

	if ping_err := ping_cmd.Start(); ping_err != nil {
		c.Logger().Error(ping_err)
	}
	var line string

	if host, exist := hosts[host_name]; exist {
		for {
			line, err = scanner.ReadString('\n')
			//c.Logger().Print(line)
			channel <- line
			host.Buffer += line
			if err != nil {
				break
			}

		}
	}
	//channel <- buf.String()
	ping_cmd.Wait()

}
func getPing(c echo.Context) error {
	for host := range hosts {

	}
}
func getStatus(c echo.Context) error {
	host_name := c.Param("hostname")
	if host, exist := hosts[host_name]; exist {
		wait := c.QueryParam("wait")
		if wait == "true" {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
			c.Response().WriteHeader(http.StatusOK)
			for buf := range host.Channel {
				if _, err := c.Response().Write([]byte(buf)); err != nil {
					return err
				}
				c.Response().Flush()
			}
			return nil
		} else {
			return c.String(http.StatusOK, host.Buffer)
		}
	} else {
		return c.String(http.StatusNotFound, "No registered hostname")
	}
}
func deletePing(c echo.Context) error {

	return c.String(http.StatusOK, "")
}
