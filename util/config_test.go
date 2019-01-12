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

func TestGetCfgVal(t *testing.T) {
	if err := InitConf("../testdata/app.json"); err != nil {
		t.Fatal("init config file error")
	}
	// test int
	testInt, err := GetCfgVal(1000, "test", "test1")
	if err != nil {
		t.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testInt != 1234 {
		t.Fatalf("want: 1234, got: %d", testInt)
	}
	if Int, _ := GetIntCfgVal(1000, "test", "test1"); Int != 1234 {
		t.Fatalf("want: 1234, got: %d", Int)
	}
	// test string
	testStr, err := GetCfgVal("", "test", "test2", "test21")
	if err != nil || testStr != "test" {
		t.Fatalf("want: test, got: %s, error: %s", testStr, err.Error())
	}
	if Str, _ := GetStringCfgVal("", "test", "test2", "test21"); Str != "test" {
		t.Fatalf("want: test, got: %s", Str)
	}
	// test bool
	testbool, err := GetCfgVal(false, "test", "test5")
	if err != nil {
		t.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testbool.(bool) != true {
		t.Fatalf("want: true, got: %v", testbool)
	}
	if Bool, _ := GetBoolCfgVal(false, "test", "test5"); !Bool {
		t.Fatalf("want: true, got: %t", Bool)
	}
	// test float
	testfloat, err := GetCfgVal(0.0, "test", "test2", "test22")
	if err != nil || testfloat != 2.2 {
		t.Fatalf("want: abc, got: %f, error: %s", testStr, err.Error())
	}
	if Float, _ := GetFloatCfgVal("", "test", "test2", "test22"); Float != 2.2 {
		t.Fatalf("want: 2.2, got: %f", Float)
	}
	// test string
	testStr, err = GetCfgVal("", "test", "test3")
	if err != nil || testStr != "abc" {
		t.Fatalf("want: abc, got: %s, error: %s", testStr, err.Error())
	}
	// test array
	array, err := GetCfgVal([]string{}, "test", "test4")
	if err != nil {
		t.Fatal(err)
	}
	ss := Conv2StrSlice(array.([]interface{}))
	if ss[0] != "abc" || ss[1] != "abd" || ss[2] != "hij" {
		t.Fatal("want: [abc abd hij], got: %v", ss)
	}
	// test wrong config
	testStr, err = GetCfgVal("", "test", "test100")
	if err == nil {
		t.Fatalf("test|test4 test wrong")
	}
	testStr, err = GetCfgVal("", "test", "test2", "test5")
	if err == nil {
		t.Fatalf("test|test2|test5 test wrong, got testStr: %v", testStr)
	}
}

func TestGetSpecCfgVal(t *testing.T) {
	var cfgData CfgData
	cfgData, err := InitSpecConf("../testdata/app.json")
	if err != nil {
		t.Fatal("init config file error")
	}
	testInt, err := GetSpecCfgVal(cfgData, 1000, "test", "test1")
	if err != nil {
		t.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testInt != 1234 {
		t.Fatalf("want: 1234, got: %d", testInt)
	}
	testStr, err := GetSpecCfgVal(cfgData, "", "test", "test2", "test21")
	if err != nil || testStr != "test" {
		t.Fatalf("want: test, got: %s, error: %s", testStr, err.Error())
	}
	testfloat, err := GetSpecCfgVal(cfgData, 0.0, "test", "test2", "test22")
	if err != nil || testfloat != 2.2 {
		t.Fatalf("want: abc, got: %f, error: %s", testStr, err.Error())
	}
	testStr, err = GetSpecCfgVal(cfgData, "", "test", "test3")
	if err != nil || testStr != "abc" {
		t.Fatalf("want: abc, got: %s, error: %s", testStr, err.Error())
	}
	// test wrong config
	testStr, err = GetSpecCfgVal(cfgData, "", "test", "test100")
	if err == nil {
		t.Fatalf("test|test4 test wrong")
	}
	testStr, err = GetSpecCfgVal(cfgData, "", "test", "test2", "test5")
	if err == nil {
		t.Fatalf("test|test2|test5 test wrong")
	}
	// test int array
	ints, err := GetIntArraySpecCfgVal(cfgData, nil, "test", "test6")
	if err != nil {
		t.Fatal(err)
	}
	if ints[0] != 1 || ints[1] != 3 || ints[2] != 5 {
		t.Fatalf("got: %+v, want: [1, 3, 5]", ints)
	}
	// test float array
	ff, err := GetFloatArraySpecCfgVal(cfgData, nil, "test", "test7")
	if err != nil {
		t.Fatal(err)
	}
	if ff[0] != 1.1 || ff[1] != 3.3 || ff[2] != 5.5 {
		t.Fatalf("got: %+v, want: [1.1, 3.3, 5.5]", ff)
	}
}
