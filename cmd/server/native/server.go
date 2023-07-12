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
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"bakiatmaca.com/poc/daemonandcli/cmd/server/dummy"
	usock "bakiatmaca.com/poc/daemonandcli/internal"
)

var (
	service *dummy.FakeService
)

func main() {

	//remove exist sock
	os.Remove(usock.ADDRESS)

	//unix server svrSock
	svrSock, err := net.Listen(usock.PROTOCOL, usock.ADDRESS)

	if err != nil {
		log.Println("Connection Error", err.Error())
		os.Exit(-1)
	} else {
		log.Println("Daemon service started")
	}

	//create a fake service
	service = dummy.NewFakeService()

	go func() {
		for {
			con, err := svrSock.Accept()

			if err != nil {
				log.Println("Accept Error", err.Error())
				os.Exit(-1)
			}

			ctx := context.TODO()

			sendCh, recvCh := usock.SockSendRecv(ctx, con)
			hook(ctx, sendCh, recvCh)
		}
	}()

	ex := make(chan os.Signal, 2)
	signal.Notify(ex, syscall.SIGINT, syscall.SIGTERM)
	<-ex

	svrSock.Close()
}

func hook(ctx context.Context, sendChan chan string, recvChan chan string) {
	go func() {
		for {
			select {
			case cmd, ok := <-recvChan:
				if ok {
					go func() {
						result := ""

						if cmd == "Get-Capacity" {
							result = service.GetCapacity()
						} else if cmd == "Get-Uptime" {
							result = dummy.FormatDuration(service.GetUptime())
						} else if cmd == "Get-ConStr" {
							result = service.ConnectionString
						} else if strings.HasPrefix(cmd, "setconstr=") {
							service.ConnectionString = strings.Split(cmd, "=")[1]
							result = "OK"
						} else if cmd == "Get-MaxThread" {
							result = strconv.Itoa(service.MaxThreadCount)
						} else if strings.HasPrefix(cmd, "Set-MaxThread=") {
							c, _ := strconv.Atoi(strings.Split(cmd, "=")[1])
							if c > 0 {
								service.MaxThreadCount = c
								result = "OK"
							} else {
								result = "wrong argument"
							}
						}

						if len(result) > 0 {
							sendChan <- fmt.Sprintf("%s\n", result)
						}
					}()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
