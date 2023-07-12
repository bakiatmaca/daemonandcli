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
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	usock "bakiatmaca.com/poc/daemonandcli/internal"
	"github.com/go-resty/resty/v2"
)

var (
	capacity     = flag.Bool("cap", false, "get capacity")
	uptime       = flag.Bool("uptime", false, "get uptime")
	getConstr    = flag.Bool("constr", false, "get connection string")
	conStr       = flag.String("setconstr", "", "set connection string")
	getMaxThread = flag.Bool("thredcount", false, "get maximum thread count")
	maxThread    = flag.Int("setthredcount", 0, "set maximum thread count")
)

var (
	client *resty.Client
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*time.Second))

	//unix client sock connect to server sock
	sock, err := net.Dial(usock.PROTOCOL, usock.ADDRESS)

	if err != nil {
		log.Println("Connection Error", err.Error())
		os.Exit(-1)
	}

	//unix domain socket bind to http client
	hc := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return sock, err
			},
		},
	}
	client = resty.NewWithClient(hc)
	client.SetTimeout(3 * time.Second)

	//invoke
	invoke(cancel)

	<-ctx.Done()
	sock.Close()
}

func invoke(cancel func()) {
	var b string

	if *capacity {
		b = getRequest("cap")
	} else if *uptime {
		b = getRequest("uptime")
	} else if *getConstr {
		b = getRequest("constr")
	} else if *getMaxThread {
		b = getRequest("thredcount")
	} else if len(*conStr) > 0 {
		b = doRequest(&usock.CmdBundle{
			Cmd:   "setconstr",
			Value: *conStr,
		}, false)
	} else if *maxThread > 0 {
		b = doRequest(&usock.CmdBundle{
			Cmd:   "setthredcount",
			Value: *maxThread,
		}, false)
	} else {
		fmt.Println("usage to help -help")
		os.Exit(0)
	}

	if len(b) > 0 {
		fmt.Print(strings.ReplaceAll(b, "\"", ""))
	}

	cancel()
}

func doRequest(cmdb *usock.CmdBundle, isget bool) string {
	var resp *resty.Response
	var err error

	if isget {
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			Get(fmt.Sprintf("http://uds/get/%s", cmdb.Cmd))
	} else {
		resp, err = client.R().SetHeader("Content-Type", "application/json").
			SetBody(*cmdb).
			Post("http://uds/set")
	}

	if err != nil {
		fmt.Println("Connection error", err.Error())
	}

	if resp.StatusCode() != http.StatusOK {
		fmt.Println("execution error")
	}

	return string(resp.Body())
}

func getRequest(cmd string) string {
	return doRequest(&usock.CmdBundle{
		Cmd: cmd,
	}, true)
}
