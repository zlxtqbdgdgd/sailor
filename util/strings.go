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

package util

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func ContainsStr(s []string, t string) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}

func EscapeChars(s, chars string) string {
	if !strings.ContainsAny(s, chars) {
		return s
	}
	buf := make([]byte, 0, len(s)+32)
	for prev, hit := 0, 0; ; prev += hit + 1 {
		if hit = strings.IndexAny(s[prev:], chars); hit < 0 {
			buf = append(buf, s[prev:]...)
			break
		}
		buf = append(buf, s[prev:prev+hit]...)
		buf = append(buf, 0x5C) // 0x5C = \
		buf = append(buf, s[prev+hit])
	}
	var r string
	// conversions will not copy the []byte values since Go v1.8
	r = (" " + string(buf))[1:]
	return r
}

func IsStartedWithNumber(s string) bool {
	for i, v := range s {
		if i == 0 && unicode.IsNumber(v) {
			return true
		} else {
			return false
		}
	}
	return false
}

// meet condition which only a bit case in large cases
func RmNonLetters(s string) (r string) {
	r = s
	for _, v := range s {
		// 0x005F == "_", 0x0030 = 0, 0x0039 = 9
		if !unicode.IsLetter(v) && !(v >= 0x0030 && v <= 0x0039) && v != 0x005F {
			r = ""
			break
		}
	}
	if r != "" {
		return s
	}
	for _, v := range s {
		if unicode.IsLetter(v) || (v >= 0x0030 && v <= 0x0039) || v == 0x005F {
			r += string(v)
		}
	}
	return r
}

func Conv2Str(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return "NULL"
	case bool:
		if x {
			return "TRUE"
		}
		return "FALSE"
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case uint:
		return strconv.FormatInt(int64(v.(uint)), 10)
	case uint32:
		return strconv.FormatInt(int64(v.(uint32)), 10)
	case uint64:
		return strconv.FormatInt(int64(v.(uint64)), 10)
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(v.(float64), 'G', -1, 64)
	case string:
		return v.(string)
	default:
		return ""
	}
}

func SqueezeStrings(s []string) []string {
	if len(s) < 2 {
		return s
	}
	j := 0
	for i, v := range s {
		if i == 0 {
			continue
		}
		if v != s[j] && i != j+1 {
			j++
			s[j] = v
		}
	}
	return s[:j+1]
}

func Conv2StrSlice(si []interface{}) []string {
	var ss []string
	for _, v := range si {
		ss = append(ss, fmt.Sprint(v))
	}
	return ss
}

func Conv2FloatSlice(fi []interface{}) []float64 {
	var ff []float64
	for _, v := range fi {
		if f, ok := v.(float64); ok {
			ff = append(ff, f)
		} else {
			return nil
		}
	}
	return ff
}
