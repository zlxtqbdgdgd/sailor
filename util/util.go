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

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func LoadJson(fileName string, result interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return nil
}

func GetJsonSegment(raw []byte, v interface{}, keys ...string) error {
	var message map[string]*json.RawMessage
	err := json.Unmarshal(raw, &message)
	if err != nil {
		return err
	}
	var m interface{}
	m = message
	for i, k := range keys {
		if m, ok := m.(map[string]*json.RawMessage); ok {
			if d, ok := m[k]; ok {
				if i == len(keys)-1 {
					if err1 := json.Unmarshal(*d, &v); err1 != nil {
						return errors.New(fmt.Sprintf("unmarshal segment error, key: %v,"+
							" error: %s", keys, err1))
					}
					return nil
				}
				if err1 := json.Unmarshal(*d, &m); err1 != nil {
					return errors.New(fmt.Sprintf("unmarshal segment error, key: %v,"+
						" error: %s", keys, err1))
				}
			}
		}
	}
	return errors.New(fmt.Sprintf("invalid Key: %v", keys))
}
