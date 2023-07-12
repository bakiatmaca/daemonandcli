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
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	usock "bakiatmaca.com/poc/daemonandcli/internal"
)

func main_for_base_http_server() {
	//remove exist sock
	os.Remove(usock.ADDRESS)

	//unix server sock
	svrSock, err := net.Listen(usock.PROTOCOL, usock.ADDRESS)

	if err != nil {
		log.Println("Connection Error", err.Error())
		os.Exit(-1)
	}

	http.HandleFunc("/hello", hello)

	//unix domain socket bind to http server
	errb := http.Serve(svrSock, nil)
	if errb != nil {
		log.Fatalln("UDS bind Error", errb.Error())
	}

	ex := make(chan os.Signal, 2)
	signal.Notify(ex, syscall.SIGINT, syscall.SIGTERM)
	<-ex
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}
