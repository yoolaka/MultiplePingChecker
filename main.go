package main

import (
	"bufio"
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
	e.POST("/ping", CreatePing)
	e.GET("/ping", GetPing)
	e.GET("/ping/:hostname", GetPingStatus)
	e.DELETE("/ping/:hostname", DeletePing)

	// Start server
	e.Logger.Fatal(e.Start(":9335"))
}
func CreatePing(c echo.Context) error {
	host_name := c.FormValue("server")
	count_string := c.FormValue("count")
	count, err := strconv.Atoi(count_string)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invaild count")
	}
	if _, exist := hosts[host_name]; exist {
		c.Logger().Error("Host already exists")
		return c.String(http.StatusBadRequest, "")
	}

	var buffer string

	channel := make(chan string, 1)
	host := &ping_host{
		host_name,
		count,
	}
	host_attr := &ping_host_attr{
		channel,
		buffer,
	}

	ping_cmd := exec.Command("ping", "-c", count_string, host_name)
	hosts[host.HostName] = host
	host_attrs[host.HostName] = host_attr
	ping_cmds[host.HostName] = ping_cmd

	go ExecutePing(c, ping_cmd, host_name, channel)

	return c.String(http.StatusCreated, "")
}

func ExecutePing(c echo.Context, ping_cmd *exec.Cmd, host_name string, channel chan<- string) {

	stdout, err := ping_cmd.StdoutPipe()
	if err != nil {
		c.Logger().Error(err)
	}
	scanner := bufio.NewReader(stdout)

	if ping_err := ping_cmd.Start(); ping_err != nil {
		c.Logger().Error(ping_err)
	}
	var line string

	if host_attr, exist := host_attrs[host_name]; exist {
		for {
			line, err = scanner.ReadString('\n')
			c.Logger().Print(line)
			channel <- line
			host_attr.Buffer += line
			if err != nil {
				break
			}

		}
	}
	ping_cmd.Wait()

}
func GetPing(c echo.Context) error {
	json_hosts := make([]ping_host, 0)
	for _, host := range hosts {
		json_hosts = append(json_hosts, *host)
	}
	return c.JSON(http.StatusOK, json_hosts)
}
func GetPingStatus(c echo.Context) error {
	host_name := c.Param("hostname")
	if host_attr, exist := host_attrs[host_name]; exist {
		wait := c.QueryParam("wait")
		if wait == "true" {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
			c.Response().WriteHeader(http.StatusOK)
			for buf := range host_attr.Channel {
				if _, err := c.Response().Write([]byte(buf)); err != nil {
					return err
				}
				c.Response().Flush()
			}
			return nil
		} else {
			return c.String(http.StatusOK, host_attr.Buffer)
		}
	} else {
		return c.String(http.StatusNotFound, "No registered hostname")
	}
}
func DeletePing(c echo.Context) error {
	host_name := c.Param("hostname")
	if ping_cmd, exist := ping_cmds[host_name]; exist {
		c.Logger().Print("here")
		ping_cmd.Process.Kill()
	}
	return c.String(http.StatusNoContent, "")
}
