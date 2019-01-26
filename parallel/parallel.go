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

package parallel

import (
	"log"
	"time"
)

// Any job must implement the Runnable interface
type Runnable interface {
	Run() (ret map[string]interface{}, err error)
}

// Jobs is the job set that need to run parallel
type Jobs struct {
	jobs          []Runnable //jobs
	timeoutInMs   int        //timeout ms for all jobs
	maxRunJobsNum int        //max runnig jobs number
	chans         chan map[string]interface{}
	rets          map[string]interface{}
}

// Jobs's constuctor
func NewJobs(jobs []Runnable, timeoutInMs int, maxRunJobsNum int) *Jobs {
	return &Jobs{jobs: jobs, timeoutInMs: timeoutInMs, maxRunJobsNum: maxRunJobsNum}
}

func (p *Jobs) callOne(one Runnable) {
	ret, _ := one.Run()
	p.chans <- ret
}

// Run jobs parallel
func (p *Jobs) Run() (err error) {
	size := len(p.jobs)
	p.chans = make(chan map[string]interface{}, size)
	p.rets = make(map[string]interface{})
	i, j := 0, 0
	for i < p.maxRunJobsNum && i < size {
		go p.callOne(p.jobs[i])
		i++
	}
DONE:
	for {
		select {
		case ret := <-p.chans:
			for key, val := range ret {
				if nil != val {
					p.rets[key] = val
				}
			}
			j++
			if i < size {
				go p.callOne(p.jobs[i])
				i++
			}
			if j >= size {
				break DONE
			}
		case <-time.After(time.Duration(p.timeoutInMs) * time.Millisecond):
			log.Printf("Warning] job time over %d", p.timeoutInMs)
			break DONE
		}
	}
	
	return
}

// Get Jobs result
func (p *Jobs) Result() map[string]interface{} {
	return p.rets
}
