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
	"bytes"
	"errors"
)

func ReadlnTrimmed(r *bufio.Reader, lineLimit int) (string, error) {
	isPrefix := true
	var err error
	var ln []byte
	count := 0
	if lineLimit <= 0 {
		return "", errors.New("read file count not more than 0")
	}
	line := make([]byte, 0, lineLimit)
	buf := bytes.NewBuffer(line)
	for isPrefix && err == nil {
		ln, isPrefix, err = r.ReadLine()
		if count+len(ln) <= lineLimit {
			n, _ := buf.Write(ln)
			count += n
		} else if count < lineLimit {
			n, _ := buf.Write(ln[:lineLimit-count])
			count += n
		}
	}
	return buf.String(), err
}
