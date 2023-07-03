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
	"os"
	"time"

	usock "bakiatmaca.com/poc/daemonandcli/internal"
)

var (
	capacity     = flag.Bool("cap", false, "get capacity")
	uptime       = flag.Bool("uptime", false, "get uptime")
	getConstr    = flag.Bool("constr", false, "get connection string")
	conStr       = flag.String("setconstr", "", "set connection string")
	getMaxThread = flag.Bool("thredcount", false, "get maximum thread count")
	maxThread    = flag.Int("setthredcount", 0, "set maximum thread count")
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

	sendCh, recvCh := usock.SockSendRecv(ctx, sock)
	showResponse(ctx, cancel, recvCh)

	//sendCmd
	sendCmd(sendCh)

	<-ctx.Done()
	sock.Close()

}

func showResponse(ctx context.Context, cancel func(), recvChan chan string) {
	go func() {
		select {
		case data, ok := <-recvChan:
			if ok {
				fmt.Println(data)
				cancel()
			}
		case <-ctx.Done():
			return
		}
	}()
}

func sendCmd(sendChan chan string) {
	if *capacity {
		sendChan <- "Get-Capacity\n"
	} else if *uptime {
		sendChan <- "Get-Uptime\n"
	} else if *getConstr {
		sendChan <- "Get-ConStr\n"
	} else if len(*conStr) > 0 {
		sendChan <- fmt.Sprintf("setconstr=%s\n", *conStr)
	} else if *getMaxThread {
		sendChan <- "Get-MaxThread\n"
	} else if *maxThread > 0 {
		sendChan <- fmt.Sprintf("Set-MaxThread=%d\n", *maxThread)
	} else {
		fmt.Println("usage to help -help")
		os.Exit(0)
	}
}
