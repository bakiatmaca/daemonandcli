//This program is free software; you can redistribute it and/or
//modify it under the terms of the GNU General Public License
//as published by the Free Software Foundation; either version 2
//of the License, or (at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program; If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"bakiatmaca.com/poc/daemonandcli/cmd/server/dummy"
	usock "bakiatmaca.com/poc/daemonandcli/internal"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

var (
	fakeService *dummy.FakeService
)

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HideBanner = true

	//remove exist sock
	os.Remove(usock.ADDRESS)

	//unix server sock
	svrSock, err := net.Listen(usock.PROTOCOL, usock.ADDRESS)

	if err != nil {
		log.Println("Connection Error", err.Error())
		os.Exit(-1)
	}

	//unix domain socket bind to http server
	e.Listener = svrSock

	//create fakeservice
	fakeService = dummy.NewFakeService()

	e.GET("/get/:cmd", getValue)
	e.POST("/set", setValue)

	log.Println("Daemon service started")
	e.Logger.Fatal(e.Start(""))
}

func getValue(c echo.Context) error {

	cmd := c.Param("cmd")

	switch cmd {
	case "cap":
		return c.JSON(http.StatusOK, fakeService.GetCapacity())
	case "uptime":
		return c.JSON(http.StatusOK, dummy.FormatDuration(fakeService.GetUptime()))
	case "constr":
		return c.JSON(http.StatusOK, fakeService.ConnectionString)
	case "thredcount":
		return c.JSON(http.StatusOK, fakeService.MaxThreadCount)
	default:
	}

	return c.NoContent(http.StatusBadRequest)
}

func setValue(c echo.Context) error {

	cmdb := new(usock.CmdBundle)

	if err := c.Bind(cmdb); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	switch cmdb.Cmd {
	case "setconstr":
		v, _ := cmdb.Value.(string)
		fakeService.ConnectionString = v
		return c.JSON(http.StatusOK, "Ok")
	case "setthredcount":
		v, _ := cmdb.Value.(int)
		fakeService.MaxThreadCount = v
		return c.JSON(http.StatusOK, "Ok")
	default:
	}

	return c.NoContent(http.StatusBadRequest)
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
