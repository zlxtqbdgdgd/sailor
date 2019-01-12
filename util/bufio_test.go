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
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestReadlnTrimmed(t *testing.T) {
	data := "hello world\nlet's start test TestReadlnTrimmed\ncome on!"
	b := bufio.NewReader(strings.NewReader(data))
	results := ""
	if result, err := ReadlnTrimmed(b, 100); err == nil {
		results += result + "\n"
	}
	if result, err := ReadlnTrimmed(b, 100); err == nil {
		results += result + "\n"
	}
	if result, err := ReadlnTrimmed(b, 100); err == nil {
		results += result
	}
	if results != data {
		t.Fatal("ReadlnTrimmed test failed, want:", data, "\ngot:", results)
	}
	// test trim
	data = "ABCDEFGHIJKLMN\n1234567890123"
	b = bufio.NewReader(strings.NewReader(data))
	result, err := ReadlnTrimmed(b, 10)
	if err != nil {
		t.Fatal("ReadlnTrimmed test failed, err:", err)
	}
	if result != "ABCDEFGHIJ" {
		t.Fatal("ReadlnTrimmed test failed, want: ABCDEFGHIJ", "\ngot:", result)
	}
	result, err = ReadlnTrimmed(b, 10)
	if err != nil {
		t.Fatal("ReadlnTrimmed test failed, err:", err)
	}
	if result != "1234567890" {
		t.Fatal("ReadlnTrimmed test failed, want: 1234567890", "\ngot:", result)
	}
	result, err = ReadlnTrimmed(b, 10)
	if err != io.EOF {
		t.Fatal("ReadlnTrimmed test failed, err:", err)
	}
	// test buffer more than 64kb(bufio.MaxScanTokenSize)
	data = ""
	for i := 0; i < 66*1024; i++ {
		data += "A"
	}
	b = bufio.NewReader(strings.NewReader(data))
	result, err = ReadlnTrimmed(b, 65*1024)
	if err != nil {
		t.Fatal("ReadlnTrimmed test failed, err:", err)
	}
	if len(result) != 65*1024 {
		t.Fatal("ReadlnTrimmed test failed, exceed max buf size!")
	}
}
