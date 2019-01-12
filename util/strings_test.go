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
	"testing"
)

func TestEscapeChars(t *testing.T) {
	chars := `"',`
	s := `AAAAAA`
	r := EscapeChars(s, chars)
	if r != `AAAAAA` {
		t.Fatalf(`EscapeChars failed, want: "AAAAAA", got: %s`, r)
	}

	s = `'AA"A'AA,A`
	r = EscapeChars(s, chars)
	if r != `\'AA\"A\'AA\,A` {
		t.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A", got: %s`, r)
	}

	s = `AA"A'AA,A"`
	chars = `"',`
	r = EscapeChars(s, chars)
	if r != `AA\"A\'AA\,A\"` {
		t.Fatalf(`EscapeChars failed, want: "AA\"A\'AA\,A\"", got: %s`, r)
	}

	s = `'AA"A'AA,A"`
	r = EscapeChars(s, chars)
	if r != `\'AA\"A\'AA\,A\"` {
		t.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A\"", got: %s`, r)
	}

	s = `'AA"刘德华A'AA,A哈哈"`
	r = EscapeChars(s, chars)
	if r != `\'AA\"刘德华A\'AA\,A哈哈\"` {
		t.Fatalf(`EscapeChars failed, want: "\'AA\"A\'AA\,A\"", got: %s`, r)
	}
}

func TestRmNonLetters(t *testing.T) {
	s := `\Hello 世界`
	r := RmNonLetters(s)
	if r != `Hello世界` {
		t.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello  世界`
	r = RmNonLetters(s)
	if r != `Hello世界` {
		t.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello12`
	r = RmNonLetters(s)
	if r != `Hello12` {
		t.Fatalf(`RmNonLetters failed, want: "Hello12", got: %s`, r)
	}

	s = `Hello，。、|世界!`
	r = RmNonLetters(s)
	if r != `Hello世界` {
		t.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}

	s = `Hello~!@#$%^&*()-+=<>?，。、|世界!`
	r = RmNonLetters(s)
	if r != `Hello世界` {
		t.Fatalf(`RmNonLetters failed, want: "Hello世界", got: %s`, r)
	}
}

func TestConv2Str(t *testing.T) {
	s := Conv2Str(nil)
	if s != "NULL" {
		t.Fatal("error, should be NULL")
	}
	s = Conv2Str(true)
	if s != "TRUE" {
		t.Fatal("error, should be TRUE")
	}
	s = Conv2Str(false)
	if s != "FALSE" {
		t.Fatal("error, should be FALSE")
	}
	s = Conv2Str(123)
	if s != "123" {
		t.Fatal("error, should be 123")
	}
	s = Conv2Str(123.45)
	if s != "123.45" {
		t.Fatal("error, should be 123.45")
	}
	s = Conv2Str("abc")
	if s != "abc" {
		t.Fatal("error, should be abc")
	}
}

func TestSqueezeStrings(t *testing.T) {
	s := []string{"a", "ab", "ab", "abc", "dd", "dd", "defg"}
	s = SqueezeStrings(s)
	if len(s) != 5 ||
		s[0] != "a" || s[1] != "ab" || s[2] != "abc" || s[3] != "dd" ||
		s[4] != "defg" {
		t.Fatalf("got: %v, want: %v", s, []string{"a", "ab", "abc", "dd", "defg"})
	}
}

func TestConvSlice2Str(t *testing.T) {
	si := []interface{}{"a", "ab", "abc", "dd"}
	ss := Conv2StrSlice(si)
	if ss[0] != "a" || ss[1] != "ab" || ss[2] != "abc" || ss[3] != "dd" {
		t.Fatalf("want: [a ab abc dd], got: %v", ss)
	}
}
