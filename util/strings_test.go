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

package util_test

import (
	gc "gopkg.in/check.v1"
	"github.com/zlxtqbdgdgd/sailor/util"
)

type stringSuite struct {
}

var _ = gc.Suite(&stringSuite{})

func (*sliceSuite) TestEscapeChars(c *gc.C) {
	chars := `"',`
	s := `AAAAAA`
	r := util.EscapeChars(s, chars)
	if r != `AAAAAA` {
		c.Fatalf(`EscapeChars failed, want: "AAAAAA", got: %s`, r)
	}

	s = `'AA"A'AA,A`
	r = util.EscapeChars(s, chars)
	if r != `\'AA\"A\'AA\,A` {
		c.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A", got: %s`, r)
	}

	s = `AA"A'AA,A"`
	chars = `"',`
	r = util.EscapeChars(s, chars)
	if r != `AA\"A\'AA\,A\"` {
		c.Fatalf(`EscapeChars failed, want: "AA\"A\'AA\,A\"", got: %s`, r)
	}

	s = `'AA"A'AA,A"`
	r = util.EscapeChars(s, chars)
	if r != `\'AA\"A\'AA\,A\"` {
		c.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A\"", got: %s`, r)
	}

	s = `'AA"刘德华A'AA,A哈哈"`
	r = util.EscapeChars(s, chars)
	if r != `\'AA\"刘德华A\'AA\,A哈哈\"` {
		c.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A\"", got: %s`, r)
	}
}

func (*sliceSuite) TestRmNonLetters(c *gc.C) {
	s := `\Hello 世界`
	r := util.RmNonLetters(s)
	if r != `Hello世界` {
		c.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello  世界`
	r = util.RmNonLetters(s)
	if r != `Hello世界` {
		c.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello12`
	r = util.RmNonLetters(s)
	if r != `Hello12` {
		c.Fatalf(`RmNonLetters failed, want: "Hello12", got: %s`, r)
	}

	s = `Hello，。、|世界!`
	r = util.RmNonLetters(s)
	if r != `Hello世界` {
		c.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello~!@#$%^&*()-+=<>?，。、|世界!`
	r = util.RmNonLetters(s)
	if r != `Hello世界` {
		c.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}
}

func (*sliceSuite) TestConv2Str(c *gc.C) {
	s := util.Conv2Str(nil)
	if s != "NULL" {
		c.Fatal("error, should be NULL")
	}
	s = util.Conv2Str(true)
	if s != "TRUE" {
		c.Fatal("error, should be TRUE")
	}
	s = util.Conv2Str(false)
	if s != "FALSE" {
		c.Fatal("error, should be FALSE")
	}
	s = util.Conv2Str(123)
	if s != "123" {
		c.Fatal("error, should be 123")
	}
	s = util.Conv2Str(123.45)
	if s != "123.45" {
		c.Fatal("error, should be 123.45")
	}
	s = util.Conv2Str("abc")
	if s != "abc" {
		c.Fatal("error, should be abc")
	}
}

func (*sliceSuite) TestSqueezeStrings(c *gc.C) {
	s := []string{"a", "ab", "ab", "abc", "dd", "dd", "defg"}
	s = util.SqueezeStrings(s)
	if len(s) != 5 ||
		s[0] != "a" || s[1] != "ab" || s[2] != "abc" || s[3] != "dd" ||
		s[4] != "defg" {
		c.Fatalf("got: %v, want: %v", s, []string{"a", "ab", "abc", "dd", "defg"})
	}
}

func (*sliceSuite) TestConvSlice2Str(c *gc.C) {
	si := []interface{}{"a", "ab", "abc", "dd"}
	ss := util.Conv2StrSlice(si)
	if ss[0] != "a" || ss[1] != "ab" || ss[2] != "abc" || ss[3] != "dd" {
		c.Fatalf("want: [a ab abc dd], got: %v", ss)
	}
}
