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

package math

import "math"

func MaxF64(values ...float64) (float64, int) {
	j, m := 0, values[0]
	for i, v := range values[1:] {
		if m < v {
			m = v
			j = i + 1
		}
	}
	return m, j
}

func MinF64(values ...float64) (float64, int) {
	j, m := 0, values[0]
	for i, v := range values[1:] {
		if v < m {
			m = v
			j = i + 1
		}
	}
	return m, j
}

func AveF64(values ...float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// prec short for precison
func EqualF64(a, b float64, prec int) bool {
	if prec <= 0 {
		prec = 6
	}
	var delta float64 = 1.0
	for prec > 0 {
		delta *= 0.1
		prec--
	}
	if math.Abs(a-b) < delta {
		return true
	}
	return false
}
