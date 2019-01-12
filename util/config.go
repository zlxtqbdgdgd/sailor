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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type CfgData map[string]*json.RawMessage

var cfgData CfgData

// NOTE: not thread-safe
func InitSpecConf(cfgFile string) (CfgData, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read json file: %s, error: %s",
			cfgFile, err))
	}
	raw := make(CfgData)
	err = json.Unmarshal(cfgBytes, &raw)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to Unmarshal file: %s, error: %s",
			cfgFile, err))
	}
	return raw, nil
}

// TODO: load configure file dynamic by start a timer and goroutine
// NOTE: not thread-safe
func InitConf(cfgFile string) error {
	raw, err := InitSpecConf(cfgFile)
	if err != nil {
		return err
	}
	cfgData = raw
	return nil
}

func GetSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (
	v interface{}, err error) {
	v = def
	var m interface{}
	m = cfgData
	for i, k := range keys {
		if mm, ok := m.(CfgData); ok {
			if d, ok := mm[k]; ok {
				if i == len(keys)-1 {
					if err1 := json.Unmarshal(*d, &v); err1 != nil {
						err = errors.New(fmt.Sprintf("failed to Unmarshal config, key: %v,"+
							" error: %s", keys, err1))
					}
					if _, ok := def.(int); ok {
						v = int(v.(float64))
					}
					return
				}
				mm = map[string]*json.RawMessage{}
				if err1 := json.Unmarshal(*d, &mm); err1 != nil {
					err = errors.New(fmt.Sprintf("failed to Unmarshal config, key: %v,"+
						" error: %s", keys, err1))
					return
				}
				m = mm
			} else {
				err = errors.New(fmt.Sprintf("cfgData has no key[%s]", k))
				return
			}
		} else {
			err = errors.New(fmt.Sprintf("cfgData wrong value for key[%s]", k))
			return
		}
	}
	err = errors.New(fmt.Sprintf("GetCfgVal error: invalid Key: %v", keys))
	return
}

func GetCfgVal(def interface{}, keys ...string) (interface{}, error) {
	return GetSpecCfgVal(cfgData, def, keys...)
}

func GetStringSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (
	string, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", errors.New(fmt.Sprintf("get string cfg value[%v] assert failed,"+
			" origin: %v", keys, v))
	}
	return s, nil
}

func GetStringCfgVal(def interface{}, keys ...string) (string, error) {
	return GetStringSpecCfgVal(cfgData, def, keys...)
}

func GetIntSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (int, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return 0, err
	}
	i, ok := v.(int)
	if !ok {
		return 0, errors.New(fmt.Sprintf("get int cfg value[%v] assert failed,"+
			" origin: %v", keys, v))
	}
	return i, nil
}

func GetIntCfgVal(def interface{}, keys ...string) (int, error) {
	return GetIntSpecCfgVal(cfgData, def, keys...)
}

func GetBoolSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (bool, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, errors.New(fmt.Sprintf("get bool cfg value[%v] assert failed,"+
			" origin: %v", keys, v))
	}
	return b, nil
}

func GetBoolCfgVal(def interface{}, keys ...string) (bool, error) {
	return GetBoolSpecCfgVal(cfgData, def, keys...)
}

func GetFloatSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (float64, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return 0.0, err
	}
	f, ok := v.(float64)
	if !ok {
		return 0.0, errors.New(fmt.Sprintf("get float cfg value[%v] assert failed,"+
			" origin: %v", keys, v))
	}
	return f, nil
}

func GetFloatCfgVal(def interface{}, keys ...string) (float64, error) {
	return GetFloatSpecCfgVal(cfgData, def, keys...)
}

func GetStrArraySpecCfgVal(cfgData CfgData, def interface{}, keys ...string) ([]string, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return nil, err
	}
	ss, ok := v.([]interface{})
	if !ok {
		return nil, NewErrf("value[%+v] type assertion failed", v)
	}
	return Conv2StrSlice(ss), nil
}

func GetStrArrayCfgVal(def interface{}, keys ...string) ([]string, error) {
	return GetStrArraySpecCfgVal(cfgData, def, keys...)
}

func GetFloatArraySpecCfgVal(cfgData CfgData, def interface{}, keys ...string) ([]float64, error) {
	v, err := GetSpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return nil, err
	}
	ss, ok := v.([]interface{})
	if !ok {
		return nil, NewErrf("value[%+v] type assertion failed", v)
	}
	ff := Conv2FloatSlice(ss)
	if ff == nil {
		return nil, NewErrf("convert value for keys[%v] to []float64 error", keys)
	}
	return ff, nil
}

func GetFloatArrayCfgVal(def interface{}, keys ...string) ([]float64, error) {
	return GetFloatArraySpecCfgVal(cfgData, def, keys...)
}

func GetIntArraySpecCfgVal(cfgData CfgData, def interface{}, keys ...string) ([]int, error) {
	ff, err := GetFloatArraySpecCfgVal(cfgData, def, keys...)
	if err != nil {
		return nil, err
	}
	ints := make([]int, len(ff))
	for i, v := range ff {
		ints[i] = int(v)
	}
	return ints, nil
}

func GetIntArrayCfgVal(def interface{}, keys ...string) ([]int, error) {
	return GetIntArraySpecCfgVal(cfgData, def, keys...)
}
