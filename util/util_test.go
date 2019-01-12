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

package util_test

import (
	"github.com/zlxtqbdgdgd/sailor/util"
	gc "gopkg.in/check.v1"
)

const (
	esRespBody = `{
  "took": 6,
  "timed_out": false,
  "_shards": {
    "total": 12,
    "successful": 12,
    "failed": 0
  },
  "hits": {
    "total": 58,
    "max_score": null,
    "hits": [
      {
        "_index": "baike_all_1",
        "_type": "baike",
        "_id": "165a9a0e9f031f86dfe71541375916a8",
        "_score": null,
        "_source": {
          "content": "玫瑰：原产地中国。属蔷薇目，蔷薇科落叶灌木。玫瑰是英国的国花。",
          "pv": 4024101,
          "custom_tag": "default",
          "title": "玫瑰"
        },
        "sort": [
          4024101
        ]
      },
      {
        "_index": "baike_all_1",
        "_type": "baike",
        "_id": "0ca045924e03c50028334c4b185992d2",
        "_score": null,
        "_source": {
          "content": "《牡丹》是唐代女诗人薛涛创作的一首七言律诗。",
          "pv": 2145127,
          "custom_tag": "文学作品",
          "title": "牡丹"
        },
        "sort": [
          2145127
        ]
      }
    ]
  }
}
`
)

type utilSuite struct {
}

var _ = gc.Suite(&utilSuite{})

func (*utilSuite) TestGetJsonSegment(c *gc.C) {
	type Source struct {
		Label   string `json:"custom_tag"`
		Word    string `json:"title"`
		Content string `json:"content"`
		PV      int64  `json:"pv"`
	}
	type S struct {
		Index  string `json:"_index"`
		Type   string `json:"_type"`
		Id     string `json:"_id"`
		Score  string `json:"_score"`
		Source `json:"_source"`
		Sort   []int64 `json:"sort"`
	}
	var ss []S
	err := util.GetJsonSegment([]byte(esRespBody), &ss, "hits", "hits")
	if err != nil {
		c.Fatal("error:", err)
	}
	c.Logf("%+v", ss)
	if len(ss) != 2 ||
		ss[0].Id != "165a9a0e9f031f86dfe71541375916a8" ||
		ss[1].Id != "0ca045924e03c50028334c4b185992d2" ||
		ss[0].PV != 4024101 || ss[0].Word != "玫瑰" || ss[0].Label != "default" ||
		ss[0].Content != "玫瑰：原产地中国。属蔷薇目，蔷薇科落叶灌木。玫瑰是英国的国花。" ||
		ss[1].PV != 2145127 || ss[1].Word != "牡丹" || ss[1].Label != "文学作品" ||
		ss[1].Content != "《牡丹》是唐代女诗人薛涛创作的一首七言律诗。" {
		c.Fatal("Pickup Result error, got: ", ss, "\n====want: 玫瑰,4024101,default,"+
			"牡丹,2145127,文学作品")
	}
}
