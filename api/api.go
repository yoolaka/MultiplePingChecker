package api

import (
	"bufio"
	"github.com/MultiplePingChecker"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os/exec"
	"strconv"
)

func createPing(c echo.Context) error {
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

	go executePing(c, ping_cmd, host_name, channel)

	return c.String(http.StatusCreated, "")
}

func executePing(c echo.Context, ping_cmd *exec.Cmd, host_name string, channel chan<- string) {

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
func getPing(c echo.Context) error {
	json_hosts := make([]ping_host, 0)
	for _, host := range hosts {
		json_hosts = append(json_hosts, *host)
	}
	return c.JSON(http.StatusOK, json_hosts)
}
func getPingStatus(c echo.Context) error {
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
func deletePing(c echo.Context) error {
	host_name := c.Param("hostname")
	if ping_cmd, exist := ping_cmds[host_name]; exist {
		c.Logger().Print("here")
		ping_cmd.Process.Kill()
	}
	return c.String(http.StatusNoContent, "")
}
