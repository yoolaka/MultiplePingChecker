package api

import (
	"bufio"
	"github.com/MultiplePingChecker/temp_db"
	"github.com/labstack/echo"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"
)

type (
	ping_host struct {
		HostName string `json:"hostname"`
		Count    int    `json:"count"`
	}
	Api struct {
		TempDB temp_db.TempDB
	}
)

func (a *Api) CreatePing(c echo.Context) error {
	var buffer string
	host_name := c.FormValue("server")
	count_string := c.FormValue("count")
	count, err := strconv.Atoi(count_string)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "Invaild count")
	} else {
		if count > 1024 {
			c.Logger().Error(err)
			return c.String(http.StatusBadRequest, "Exceed count limitation")
		}
	}

	channel := make(chan string, 1024)

	ping_cmd := exec.Command("ping", "-c", count_string, host_name)
	//ping_cmd := exec.Command("./ping.sh", count_string, host_name)
	ping_cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	host_entry := &temp_db.PingHostEntry{
		host_name,
		count,
		channel,
		buffer,
		ping_cmd,
	}
	if res := a.TempDB.InsertPingHost(host_name, host_entry); res != temp_db.SUCCESS {
		return c.String(http.StatusBadRequest, "Duplicated host")
	}

	go a.ExecutePing(c, ping_cmd, host_name, channel)

	json_host := &ping_host{host_name, count}

	return c.JSON(http.StatusCreated, json_host)
}

func (a *Api) ExecutePing(c echo.Context, ping_cmd *exec.Cmd, host_name string, channel chan<- string) {

	stdout, err := ping_cmd.StdoutPipe()
	if err != nil {
		c.Logger().Error(err)
	}
	scanner := bufio.NewReader(stdout)

	if ping_err := ping_cmd.Start(); ping_err != nil {
		c.Logger().Error(ping_err)
	}
	var line string

	if host_entry, exist := a.TempDB.SearchHost(host_name); exist {
		for {
			line, err = scanner.ReadString('\n')
			c.Logger().Print(line)
			channel <- line
			host_entry.Buffer += line
			if err != nil {
				break
			}

		}
	}
	ping_cmd.Wait()
	close(channel)

}
func (a *Api) GetPing(c echo.Context) error {
	json_hosts := make([]ping_host, 0)
	for _, host_entry := range *a.TempDB.GetPingHost() {
		host := &ping_host{host_entry.HostName, host_entry.Count}
		json_hosts = append(json_hosts, *host)
	}
	return c.JSON(http.StatusOK, json_hosts)
}
func (a *Api) GetPingStatus(c echo.Context) error {
	host_name := c.Param("hostname")
	if host_entry, exist := a.TempDB.SearchHost(host_name); exist {
		wait := c.QueryParam("wait")
		if wait == "true" {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlain)
			c.Response().WriteHeader(http.StatusOK)
			for buf := range host_entry.Channel {
				if _, err := c.Response().Write([]byte(buf)); err != nil {
					return err
				}
				c.Response().Flush()
			}
			return nil
		} else {
			return c.String(http.StatusOK, host_entry.Buffer)
		}
	} else {
		return c.String(http.StatusNotFound, "No registered hostname")
	}
}
func (a *Api) DeletePing(c echo.Context) error {
	host_name := c.Param("hostname")
	if host_entry, exist := a.TempDB.SearchHost(host_name); exist {
		c.Logger().Print("here")
		host_entry.PingCmd.Process.Kill()
		if res := a.TempDB.DeletePingHost(host_name); res != temp_db.SUCCESS {
			return c.String(http.StatusNotFound, "")
		} else {
			return c.String(http.StatusNoContent, "")
		}
	} else {
		return c.String(http.StatusNotFound, "")
	}
}
