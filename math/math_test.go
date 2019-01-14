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

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) { gc.TestingT(t) }

type mathSuite struct {
}

var _ = gc.Suite(&mathSuite{})

func (*mathSuite) TestEqualF64(c *gc.C) {
	var a, b float64 = 1.001, 1.0012
	if EqualF64(a, b, 0) {
		c.Fatal("EqualF64 error, precision is 0, want: false, got: true")
	}
	if !EqualF64(a, b, 1) || !EqualF64(a, b, 2) || !EqualF64(a, b, 3) {
		c.Fatal("EqualF64 error, precision is 1/2/3, want true, got false")
	}
	if EqualF64(a, b, 4) || EqualF64(a, b, 5) || EqualF64(a, b, 10) {
		c.Fatal("EqualF64 error, precision is 4/5/10, want false, got true")
	}
}

func (*mathSuite) TestMaxMinAveF64(c *gc.C) {
	// t.Fatal("not implemented")
	var a, b, c1, d float64 = 1.1, 2.2, 1.2, 0.9
	max, i := MaxF64(a, b, c1, d)
	if max != 2.2 || i != 1 {
		c.Fatalf("test max failed, want: 2.2, 1, got: %f, %d", max, i)
	}
	min, i := MinF64(a, b, c1, d)
	if min != 0.9 || i != 3 {
		c.Fatalf("test min failed, want: 0.9, 3, got: %f, %d", max, i)
	}
	ave := AveF64(a, b, c1, d)
	if min != 0.9 || i != 3 {
		c.Fatalf("test min failed, want: 0.9, 3, got: %f, %d", max, i)
	}
	if !EqualF64(ave, 1.35, 6) {
		c.Fatalf("test AveF64 failed, want: 1.35, got: %f", ave)
	}
}
