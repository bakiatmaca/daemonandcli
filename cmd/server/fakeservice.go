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
	"math/rand"
	"time"
)

type FakeService struct {
	SeviceID         string
	ConnectionString string
	MaxThreadCount   int
	UpTime           time.Time
}

func NewFakeService() *FakeService {
	return &FakeService{
		SeviceID:         "0101d235-50b5-4abf-9e6f-8b1e1662f851",
		ConnectionString: "N/A",
		MaxThreadCount:   10,
		UpTime:           time.Now(),
	}
}

func (s *FakeService) GetUptime() time.Duration {
	return time.Since(s.UpTime)
}

func (s *FakeService) GetCapacity() string {
	return fmt.Sprintf("capacity %%%d", rand.Intn(99))
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
