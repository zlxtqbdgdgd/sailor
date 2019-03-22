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

package net

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

var starter sync.Once
var addr net.Addr

func testHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(200 * time.Millisecond)
	io.WriteString(w, "hello, world!\n")
}

func postHandler(w http.ResponseWriter, req *http.Request) {
	if raw, err := ioutil.ReadAll(req.Body); err != nil {
		log.Println("read request body error: ", err)
		fmt.Fprint(w, "error: ", err)
	} else {
		h := req.Header.Get("App-Id")
		if h != "" {
			h = ", header: " + h
		}
		time.Sleep(500 * time.Millisecond)
		io.WriteString(w, "OK: "+string(raw)+h)
	}
}

func closeHandler(w http.ResponseWriter, req *http.Request) {
	hj, _ := w.(http.Hijacker)
	conn, bufrw, _ := hj.Hijack()
	defer conn.Close()
	bufrw.WriteString("HTTP/1.1 200 OK\r\nConnection: close\r\n\r\n")
	bufrw.Flush()
}

func redirectHandler(w http.ResponseWriter, req *http.Request) {
	ioutil.ReadAll(req.Body)
	http.Redirect(w, req, "/post", 302)
}

func redirect2Handler(w http.ResponseWriter, req *http.Request) {
	ioutil.ReadAll(req.Body)
	http.Redirect(w, req, "/redirect", 302)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "START\n")
	f := w.(http.Flusher)
	f.Flush()
	time.Sleep(200 * time.Millisecond)
	io.WriteString(w, "WORKING\n")
	f.Flush()
	time.Sleep(200 * time.Millisecond)
	io.WriteString(w, "DONE\n")
	return
}

func setupMockServer(t *testing.T) {
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/redirect", redirectHandler)
	http.HandleFunc("/redirect2", redirect2Handler)
	http.HandleFunc("/close", closeHandler)
	http.HandleFunc("/slow", slowHandler)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to listen - %s", err.Error())
	}
	go func() {
		err = http.Serve(ln, nil)
		if err != nil {
			t.Fatalf("failed to start HTTP server - %s", err.Error())
		}
	}()
	addr = ln.Addr()
}

func NoTestHttpsConnection(t *testing.T) {
	resp, err := Get("https://httpbin.org/ip")
	if err != nil {
		t.Fatalf("Get request failed - %s", err.Error())
	}
	log.Println("TestHttpsConnection respnse: ", string(resp))

	_, err = Get("https://httpbin.org/delay/1", 900)
	if err == nil {
		t.Fatalf("HTTPS request should have timed out")
	}

	_, err = Get("https://httpbin.org/delay/1", 3000)
	if err != nil {
		t.Fatalf("Get request failed - %s", err.Error())
	}
}

func TestHttpClient(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })

	resp, err := Get("http://" + addr.String() + "/test")
	if err != nil {
		t.Fatalf("Get request failed - %s", err.Error())
	}
	log.Println("TestHttpClient respnse: ", string(resp))

	_, err = Get("http://"+addr.String()+"/test", 100)
	if err == nil {
		t.Fatalf("Get 2nd request should have timed out")
	}

	_, err = Get("http://"+addr.String()+"/test", 210)
	if err != nil {
		t.Fatal("Get 3rd request should not have timed out")
	}
}

func TestSlowServer(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })

	resp, err := Get("http://"+addr.String()+"/slow", 420)
	if err != nil {
		t.Fatalf("Get request failed - %s", err.Error())
	}
	log.Println("TestSlowServer respnse:\n", string(resp))
	resp, err = Get("http://"+addr.String()+"/slow", 390)
	if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("slow server request dind't return a context.DeadlineExceeded."+
			"Error - %s", err)
	}
}

func TestHttpPost(t *testing.T) {
	starter.Do(func() { setupMockServer(t) })

	url := "http://" + addr.String() + "/post"
	resp, err := Post(url, []byte("hello server"))
	if err != nil {
		t.Fatal("Post error: ", err)
	}
	if string(resp) != "OK: hello server" {
		t.Fatal("response error, want: OK: hello server, got:", string(resp))
	}
	// test post with header
	header := make(map[string]string)
	header["App-Id"] = "app_id_1"
	resp, err = PostWithHeader(url, header, []byte("hello server"))
	if err != nil {
		t.Fatal("Post error: ", err)
	}
	if string(resp) != "OK: hello server, header: app_id_1" {
		t.Fatal("response error, want: OK: hello server, header: app_id_1, got:", string(resp))
	}
	// test not timeout
	resp, err = Post(url, []byte("hello server"), 550)
	if err != nil {
		t.Fatal("Post error: ", err)
	}
	if string(resp) != "OK: hello server" {
		t.Fatal("response error, want: OK: hello server, got: ", string(resp))
	}
	// test timeout
	log.Println("start test timeout")
	resp, err = Post(url, []byte("hello server"), 100)
	if err == nil {
		t.Fatal("Post should timeout")
	}
	log.Println("timeout for post, error:", err)
	// test retry
	log.Println("start test retry")
	resp, err = Post(url, []byte("hello server"), 200, 3, 550)
	if err == nil {
		t.Fatal("Post should timeout")
	}
	log.Println("timeout for retry 3 times post, error:", err)
}

func TestPostTagdecode(t *testing.T) {
	url := "http://192.168.1.44:8086/tagKg/query"
	start := time.Now()
	_, err := Post(url, []byte(`{"query":"刘德华","userID":100}`))
	if err != nil {
		t.Logf("eclipse: %v", time.Since(start))
		t.Fatal("Post error: ", err)
	}
	t.Logf("eclipse: %v", time.Since(start))
}
