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
	gc "gopkg.in/check.v1"

	"io/ioutil"

	"github.com/zlxtqbdgdgd/sailor/util"
)

const json = `{
  "test": {
    "triple_data":"/roobo/data/pzg/singer.raw",
    "test1":1234,
    "test2": {
      "test21":"test",
      "test22":2.2
    },
    "test3":"abc",
    "test4": [
      "abc",
      "abd",
      "hij"
    ],
    "test5": true,
    "test6": [1, 3, 5],
    "test7": [1.1, 3.3, 5.5]
  }
}`

type configSuite struct {
	filepath string
}

var _ = gc.Suite(&configSuite{})

// Setupsuite 准备测试用的临时文件
func (s *configSuite) SetUpSuite(c *gc.C) {
	dir := c.MkDir() // Suite结束后会自动销毁c.MkDir()创建的目录

	tmpfile, err := ioutil.TempFile(dir, "")
	if err != nil {
		c.Errorf("Fail to create test file: %v\n", tmpfile.Name(), err)
	}

	_, err = tmpfile.Write([]byte(json))
	if err != nil {
		c.Errorf("Fail to prepare test file.%v\n", tmpfile.Name(), err)
	}
	if err := tmpfile.Sync(); err != nil {
		c.Errorf("Fail to prepare test file.%v\n", tmpfile.Name(), err)
	}
	if err := tmpfile.Close(); err != nil {
		c.Errorf("Fail to prepare test file.%v\n", tmpfile.Name(), err)
	}

	s.filepath = tmpfile.Name()
}

func (s *configSuite) TestGetCfgVal(c *gc.C) {
	if err := util.InitConf(s.filepath); err != nil {
		c.Fatal("init config file error")
	}
	// test int
	testInt, err := util.GetCfgVal(1000, "test", "test1")
	if err != nil {
		c.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testInt != 1234 {
		c.Fatalf("want: 1234, got: %d", testInt)
	}
	if Int, _ := util.GetIntCfgVal(1000, "test", "test1"); Int != 1234 {
		c.Fatalf("want: 1234, got: %d", Int)
	}
	// test string
	testStr, err := util.GetCfgVal("", "test", "test2", "test21")
	if err != nil || testStr != "test" {
		c.Fatalf("want: test, got: %s, error: %s", testStr, err.Error())
	}
	if Str, _ := util.GetStringCfgVal("", "test", "test2", "test21"); Str != "test" {
		c.Fatalf("want: test, got: %s", Str)
	}
	// test bool
	testbool, err := util.GetCfgVal(false, "test", "test5")
	if err != nil {
		c.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testbool.(bool) != true {
		c.Fatalf("want: true, got: %v", testbool)
	}
	if Bool, _ := util.GetBoolCfgVal(false, "test", "test5"); !Bool {
		c.Fatalf("want: true, got: %t", Bool)
	}
	// test float
	testfloat, err := util.GetCfgVal(0.0, "test", "test2", "test22")
	if err != nil || testfloat != 2.2 {
		c.Fatalf("want: abc, got: %f, error: %s", testStr, err.Error())
	}
	if Float, _ := util.GetFloatCfgVal("", "test", "test2", "test22"); Float != 2.2 {
		c.Fatalf("want: 2.2, got: %f", Float)
	}
	// test string
	testStr, err = util.GetCfgVal("", "test", "test3")
	if err != nil || testStr != "abc" {
		c.Fatalf("want: abc, got: %s, error: %s", testStr, err.Error())
	}
	// test array
	array, err := util.GetCfgVal([]string{}, "test", "test4")
	if err != nil {
		c.Fatal(err)
	}
	ss := util.Conv2StrSlice(array.([]interface{}))
	if ss[0] != "abc" || ss[1] != "abd" || ss[2] != "hij" {
		c.Fatal("want: [abc abd hij], got: %v", ss)
	}
	// test wrong config
	testStr, err = util.GetCfgVal("", "test", "test100")
	if err == nil {
		c.Fatalf("test|test4 test wrong")
	}
	testStr, err = util.GetCfgVal("", "test", "test2", "test5")
	if err == nil {
		c.Fatalf("test|test2|test5 test wrong, got testStr: %v", testStr)
	}
}

func (s *configSuite) TestGetSpecCfgVal(c *gc.C) {
	var cfgData util.CfgData
	cfgData, err := util.InitSpecConf(s.filepath)
	if err != nil {
		c.Fatal("init config file error")
	}
	testInt, err := util.GetSpecCfgVal(cfgData, 1000, "test", "test1")
	if err != nil {
		c.Fatalf("test test1 failed, error: %s", err.Error())
	}
	if testInt != 1234 {
		c.Fatalf("want: 1234, got: %d", testInt)
	}
	testStr, err := util.GetSpecCfgVal(cfgData, "", "test", "test2", "test21")
	if err != nil || testStr != "test" {
		c.Fatalf("want: test, got: %s, error: %s", testStr, err.Error())
	}
	testfloat, err := util.GetSpecCfgVal(cfgData, 0.0, "test", "test2", "test22")
	if err != nil || testfloat != 2.2 {
		c.Fatalf("want: abc, got: %f, error: %s", testStr, err.Error())
	}
	testStr, err = util.GetSpecCfgVal(cfgData, "", "test", "test3")
	if err != nil || testStr != "abc" {
		c.Fatalf("want: abc, got: %s, error: %s", testStr, err.Error())
	}
	// test wrong config
	testStr, err = util.GetSpecCfgVal(cfgData, "", "test", "test100")
	if err == nil {
		c.Fatalf("test|test4 test wrong")
	}
	testStr, err = util.GetSpecCfgVal(cfgData, "", "test", "test2", "test5")
	if err == nil {
		c.Fatalf("test|test2|test5 test wrong")
	}
	// test int array
	ints, err := util.GetIntArraySpecCfgVal(cfgData, nil, "test", "test6")
	if err != nil {
		c.Fatal(err)
	}
	if ints[0] != 1 || ints[1] != 3 || ints[2] != 5 {
		c.Fatalf("got: %+v, want: [1, 3, 5]", ints)
	}
	// test float array
	ff, err := util.GetFloatArraySpecCfgVal(cfgData, nil, "test", "test7")
	if err != nil {
		c.Fatal(err)
	}
	if ff[0] != 1.1 || ff[1] != 3.3 || ff[2] != 5.5 {
		c.Fatalf("got: %+v, want: [1.1, 3.3, 5.5]", ff)
	}
}
