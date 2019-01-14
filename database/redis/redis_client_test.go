package redis

import (
	"encoding/json"
	"testing"
	"time"

	gc "gopkg.in/check.v1"
)

type Mystruct struct {
	One   int
	Two   string
	Three float64
}

func Test(t *testing.T) { gc.TestingT(t) }

type redisSuite struct {
}

var _ = gc.Suite(&redisSuite{})

func (s *redisSuite) SetUpSuite(c *gc.C) {
	ConnectInit("192.168.1.33:6379", "", "5")
}

func (s *redisSuite) TestRedisClient(c *gc.C) {
	// test string
	key := "tesc.string"
	if err := SetValue("tesc.string", "abc"); err != nil {
		c.Fatal(err)
	}
	v, err := GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "abc" {
		c.Fatalf("test GetValue string, want: abc, get %v", v)
	}
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	// test byte
	key = "tesc.bytes"
	if err := SetValue(key, []byte("abcdef")); err != nil {
		c.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "abcdef" {
		c.Fatalf("test GetValue bytes wrong, want: abcdef, get %v", v)
	}
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	// test int
	key = "tesc.int"
	if err := SetValue(key, 100); err != nil {
		c.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "100" {
		c.Fatalf("test GetValue int wrong, want: 100, get %v", v)
	}
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	// test float
	key = "tesc.float"
	if err := SetValue(key, 1.23456); err != nil {
		c.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "1.23456" {
		c.Fatalf("test GetValue float wrong, want: 1.23456, get %v", v)
	}
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	// test struct
	stru := Mystruct{1, "second", 3.3333}
	key = "tesc.struct"
	bytes, _ := json.Marshal(stru)
	if err := SetValue(key, bytes); err != nil {
		c.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	rStru := Mystruct{}
	json.Unmarshal(v.([]byte), &rStru)
	c.Logf("Get struct result: %#v", rStru)
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	// test map
	m := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": 1.23,
	}
	key = "tesc.map"
	bytes, _ = json.Marshal(m)
	if err := SetValue(key, bytes); err != nil {
		c.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		c.Fatal(err)
	}
	var rM map[string]interface{}
	json.Unmarshal(v.([]byte), &rM)
	c.Logf("Get map result: %#v", rM)
	if err := Delete(key); err != nil {
		c.Fatal(err)
	}
	//test set nx
	key = "testnx"
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		c.Fatal(err)
	}
	c.Logf("result of setnx is %#v", v)
	if v == nil {
		c.Fatal("set nx key failed")
	}
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		c.Fatal(err)
	}
	if v != nil {
		c.Fatal("set nx expired failed")
	}
	time.Sleep(5 * time.Second)
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		c.Fatal(err)
	}
	if v == nil {
		c.Fatal("set nx key failed")
	}
}
