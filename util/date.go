// Copyright 2018 JXB. All Rights Reserved.
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

package util

import "time"

// input: e.g.: 5月13日，5月14日，5月15日，5月18日
// output: e.g.: 5月13日至5月15日，5月18日
// days is sorted in character
func SqueezeDaysText(inLayout string, days []string) ([]string, error) {
	var err error
	// convert string days to time days
	tDays := make([]time.Time, len(days))
	for i, v := range days {
		tDays[i], err = time.Parse(inLayout, v)
		if err != nil {
			continue // skip wrong format days
			//return nil, err
		}
	}
	var last string
	var out []string
	for i, _ := range tDays {
		if i == 0 {
			out = append(out, days[0])
			continue
		}
		diff := DaysBetween(tDays[i-1], tDays[i])
		if diff == 1 {
			last = "至" + days[i]
		} else if diff > 1 {
			if last != "" {
				out[len(out)-1] += last
			}
			out = append(out, days[i])
			last = ""
		}
		if i == len(tDays)-1 {
			out[len(out)-1] += last
		}
	}
	return out, nil
}

func FormatTime(date, inLayout, outLayout string) (string, error) {
	tm, err := time.Parse(inLayout, date)
	if err != nil {
		return "", err
	}
	return tm.Format(outLayout), nil
}

func GetOralDate(layout, d string) string {
	tDate, err := time.ParseInLocation(layout, d, time.Local)
	if err != nil {
		return d
	}
	days := DaysBetween(time.Now(), tDate)
	switch days {
	default:
		return tDate.Format("1月2日")
	case -3:
		return "大前天"
	case -2:
		return "前天"
	case -1:
		return "昨天"
	case 0:
		return "今天"
	case 1:
		return "明天"
	case 2:
		return "后天"
	case 3:
		return "大后天"
	}
}

func DaysBetween(start, end time.Time) int {
	dStart, _ := time.Parse("2006-01-02", start.Format("2006-01-02"))
	dEnd, _ := time.Parse("2006-01-02", end.Format("2006-01-02"))
	dur := dEnd.Sub(dStart)
	days := int(dur.Hours()) / 24
	return days
}
