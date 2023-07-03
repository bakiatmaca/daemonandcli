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

package internal

import (
	"context"
	"net"
)

func SockSendRecv(ctx context.Context, sock net.Conn) (chan string, chan string) {

	sendChan := make(chan string)
	recvChan := make(chan string)

	go func() { //read socket
		defer func() {
			sock.Close()
			close(sendChan)
			close(recvChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				buffer := make([]byte, 1024)
				p, err := sock.Read(buffer)

				if err != nil {
					return
				}

				if p > 1 {
					recvChan <- string(buffer[:p-1]) //"\n" ignored
				}
			}
		}
	}()

	go func() { //write socket
		for {
			select {
			case data := <-sendChan:
				_, err := sock.Write([]byte(data))

				if err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return sendChan, recvChan
}
