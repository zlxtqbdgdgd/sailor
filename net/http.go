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

package net

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	client     *http.Client
	clientOnce sync.Once
)

func Post(url string, data []byte, a ...int) ([]byte, error) {
	return PostWithHeader(url, nil, data, a...)
}

func PostWithHeader(url string, header map[string]string, data []byte,
	a ...int) ([]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	if len(header) == 0 {
		header = make(map[string]string)
		header["Content-Type"] = "application/json"
	}
	return RetryDoRequest("POST", url, header, data, timeout, retryTimes, interval)
}

func Put(url string, data []byte, a ...int) ([]byte, error) {
	return PutWithHeader(url, nil, data, a...)
}

func PutWithHeader(url string, header map[string]string, data []byte,
	a ...int) ([]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	if len(header) == 0 {
		header = make(map[string]string)
		header["Content-Type"] = "application/json"
	}
	return RetryDoRequest("PUT", url, header, data, timeout, retryTimes, interval)
}

func Get(url string, a ...int) ([]byte, error) {
	return GetWithHeader(url, nil, a...)
}

func GetWithHeader(url string, header map[string]string, a ...int) (
	[]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	return RetryDoRequest("GET", url, header, nil, timeout, retryTimes, interval)
}

func Delete(url string, a ...int) ([]byte, error) {
	return DeleteWithHeader(url, nil, a...)
}

func DeleteWithHeader(url string, header map[string]string, a ...int) (
	[]byte, error) {
	timeout, retryTimes, interval, err := parseParameters(a...)
	if err != nil {
		return nil, err
	}
	return RetryDoRequest("DELETE", url, header, nil, timeout, retryTimes, interval)
}

func RetryDoRequest(reqType, URL string, headers map[string]string, data []byte,
	timeout, retryTimes, interval int) ([]byte, error) {
	var err1 error
	for i := 0; i < retryTimes+1; i++ {
		_, body, _, err := DoRequest(reqType, URL, headers, data, timeout)
		//if _, ok := err.(*url.Error); ok {
		//	return nil, err
		//}
		//if statusCode == 404 {
		//	return nil, err
		//}
		if err != nil {
			err1 = errors.New(fmt.Sprintf("%s[try %d times]", err, i+1))
			time.Sleep(time.Duration(interval) * time.Millisecond)
			continue
		}
		return body, nil
	}
	return nil, err1
}

// reqType is one of HTTP request strings (GET, POST, PUT, DELETE, etc.)
func DoRequest(reqType, url string, headers map[string]string, data []byte,
	timeout int) (int, []byte, map[string][]string, error) {
	var reader io.Reader
	if data != nil && len(data) > 0 {
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(reqType, url, reader)
	if err != nil {
		return 0, nil, nil, err
	}
	//req.Header.Set("User-Agent", "X11;Linux x86_64")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	to := time.Duration(time.Duration(timeout) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()
	req = req.WithContext(ctx)
	client := GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}
	statusCode := resp.StatusCode
	header := resp.Header
	if statusCode != 200 && statusCode != 201 {
		return statusCode, nil, header, errors.New(
			fmt.Sprintf("response status error: %d", statusCode))
	}
	return statusCode, body, header, nil
}

func GetClient() *http.Client {
	clientOnce.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				//Proxy: http.ProxyURL(proxyUrl),
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 90 * time.Minute, // default value may be 7200s
					DualStack: true,
				}).DialContext,
				MaxIdleConns:        200,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				//TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
			Timeout: 10 * time.Second,
		}
	})
	return client
}

// Post a json via http,
// a[0] - timeout, default is 1 ms
// a[1] - retry times, default is 0 times
// a[2] - retry period, default is 1000 ms
// time unit is Millisecond
func parseParameters(a ...int) (timeout, retryTimes, interval int, err error) {
	if len(a) == 2 || len(a) > 3 {
		err = errors.New("http Post retry parameters count error")
		return
	}
	timeout = 1000
	if len(a) > 0 {
		timeout = a[0]
	}
	retryTimes, interval = 0, 10
	if len(a) == 3 {
		retryTimes, interval = a[1], a[2]
	}
	return
}
