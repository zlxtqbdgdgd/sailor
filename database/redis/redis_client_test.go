package redis

import (
	"encoding/json"
	"testing"
	"time"
)

type Mystruct struct {
	One   int
	Two   string
	Three float64
}

func init() {
	InitConf("192.168.1.33:6379", "", "5")
}

func TestRedisClient(t *testing.T) {
	// test string
	key := "test.string"
	if err := SetValue("test.string", "abc"); err != nil {
		t.Fatal(err)
	}
	v, err := GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "abc" {
		t.Fatalf("test GetValue string, want: abc, get %v", v)
	}
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	// test byte
	key = "test.bytes"
	if err := SetValue(key, []byte("abcdef")); err != nil {
		t.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "abcdef" {
		t.Fatalf("test GetValue bytes wrong, want: abcdef, get %v", v)
	}
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	// test int
	key = "test.int"
	if err := SetValue(key, 100); err != nil {
		t.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "100" {
		t.Fatalf("test GetValue int wrong, want: 100, get %v", v)
	}
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	// test float
	key = "test.float"
	if err := SetValue(key, 1.23456); err != nil {
		t.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	if vv, ok := v.([]byte); !ok || string(vv) != "1.23456" {
		t.Fatalf("test GetValue float wrong, want: 1.23456, get %v", v)
	}
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	// test struct
	stru := Mystruct{1, "second", 3.3333}
	key = "test.struct"
	bytes, _ := json.Marshal(stru)
	if err := SetValue(key, bytes); err != nil {
		t.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	rStru := Mystruct{}
	json.Unmarshal(v.([]byte), &rStru)
	t.Logf("Get struct result: %#v", rStru)
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	// test map
	m := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": 1.23,
	}
	key = "test.map"
	bytes, _ = json.Marshal(m)
	if err := SetValue(key, bytes); err != nil {
		t.Fatal(err)
	}
	v, err = GetValue(key)
	if err != nil {
		t.Fatal(err)
	}
	var rM map[string]interface{}
	json.Unmarshal(v.([]byte), &rM)
	t.Logf("Get map result: %#v", rM)
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	//test set nx
	key = "testnx"
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		t.Fatal(err)
	}
	t.Logf("result of setnx is %#v", v)
	if v == nil {
		t.Fatal("set nx key failed")
	}
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Fatal("set nx expired failed")
	}
	time.Sleep(5 * time.Second)
	if v, err = SetNxKeyAndExpire(key, "abc", "px", 3000); err != nil {
		t.Fatal(err)
	}
	if v == nil {
		t.Fatal("set nx key failed")
	}
}
