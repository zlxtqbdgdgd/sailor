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

package chart

import (
	"math/rand"
	"os"
	"testing"
)

func TestPolyline(t *testing.T) {
	//ms := []float64{4.1, 2.2, 3.3, 4.4, 2.4, 2.5, 3.6, 5.7, 8.8, 9.1, 12.4, 8.0}
	ms := make([]float64, 3000)
	for i := 0; i < 3000; i++ {
		ms[i] = rand.Float64() * 10
	}
	svg := Polyline(ms, "red")
	f, err := os.OpenFile("line.svg", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal("create file failed!")
	}
	if _, err := f.WriteString(svg); err != nil {
		t.Fatal("write string failed, ", err)
	}
	if err := f.Close(); err != nil {
		t.Fatal("close file error:", err)
	}
}
