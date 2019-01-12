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

import "runtime"

type MemCheck struct {
	mallocs uint64
	frees   uint64
}

func MemMark() MemCheck {
	var stat runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&stat)
	return MemCheck{mallocs: stat.Mallocs, frees: stat.Frees}
}

func MemLeakCheck(m MemCheck) int {
	var stat runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&stat)
	return int((stat.Mallocs - m.mallocs) - (stat.Frees - m.frees))
}
