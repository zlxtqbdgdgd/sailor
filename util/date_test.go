// Copyright 2018 ROOBO. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util_test

import (
	"time"

	"github.com/zlxtqbdgdgd/sailor/util"
	gc "gopkg.in/check.v1"
)

type dateSuite struct {
}

var _ = gc.Suite(&dateSuite{})

func runSqueezeDays(c *gc.C, layout string, days, comp []string) {
	out, err := util.SqueezeDaysText(layout, days)
	if err != nil {
		c.Fatal(err)
	}
	for i, _ := range out {
		if out[i] != comp[i] {
			c.Fatalf("got: %v, want: %v", out, comp)
		}
	}
}

func (*dateSuite) TestSqueezeDays(c *gc.C) {
	//
	days := []string{"5月13日", "5月14日", "5月15日", "5月18日"}
	comp := []string{"5月13日至5月15日", "5月18日"}
	runSqueezeDays(c, "1月2日", days, comp)
	//
	days = []string{"2018-5-13", "2018-5-14", "2018-5-15", "2018-5-18"}
	comp = []string{"2018-5-13至2018-5-15", "2018-5-18"}
	runSqueezeDays(c, "2006-1-2", days, comp)
	//
	days = []string{"5月13日", "5月15日", "5月16日", "5月17日"}
	comp = []string{"5月13日", "5月15日至5月17日"}
	runSqueezeDays(c, "1月2日", days, comp)
	//
	days = []string{"5月13日", "5月15日", "5月16日", "5月18日", "5月19日"}
	comp = []string{"5月13日", "5月15日至5月16日", "5月18日至5月19日"}
	runSqueezeDays(c, "1月2日", days, comp)
}

func runDaysBetweenTest(c *gc.C, want int, start, end time.Time) {
	//t.Logf("test -- start: %v, end: %v", start, end)
	days := util.DaysBetween(start, end)
	if days != want {
		c.Fatalf("want: %d, got: %d", want, days)
	}
}

func (*dateSuite) TestDaysBetween(c *gc.C) {
	start, _ := time.Parse("2006-01-02", "2018-05-10")
	runDaysBetweenTest(c, 0, start, start.Add(time.Minute))
	runDaysBetweenTest(c, 1, start, start.Add(24*time.Hour))
	runDaysBetweenTest(c, 2, start, start.Add(49*time.Hour))
	runDaysBetweenTest(c, -1, start, start.Add(-1*time.Minute))
	runDaysBetweenTest(c, -2, start, start.Add(-25*time.Hour))

	start, _ = time.Parse("2006-01-02 15:04", "2018-05-10 10:00")
	runDaysBetweenTest(c, 0, start, start.Add(13*time.Hour))
	runDaysBetweenTest(c, 1, start, start.Add(14*time.Hour))
	runDaysBetweenTest(c, 2, start, start.Add(49*time.Hour))
	runDaysBetweenTest(c, -0, start, start.Add(-10*time.Hour))
	runDaysBetweenTest(c, -1, start, start.Add(-25*time.Hour))
}
