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

import "sync"

type trieNode struct {
	flag  bool        // 该节点是不是一个词的终点
	value interface{} // 存储任意值
	word  string
	child map[rune]*trieNode // 孩子节点
}

type Val struct {
	Word  string
	Value interface{}
}

func newTrieNode() *trieNode {
	return &trieNode{
		flag:  false,
		value: nil,
		child: make(map[rune]*trieNode),
	}
}

type Trie struct {
	Root *trieNode
	size int
	sync.RWMutex
}

func New() *Trie {
	return &Trie{
		Root: newTrieNode(),
		size: 0,
	}
}

// 得到Trie上的词的总数
func (this *Trie) Size() int {
	return this.size
}

// 插入一个词(key, value)，会覆盖相同的key
func (this *Trie) Insert(key string, value interface{}) {
	this.Lock()
	defer this.Unlock()

	curNode := this.Root
	for _, v := range key {
		if curNode.child[v] == nil {
			curNode.child[v] = newTrieNode()
		}
		curNode = curNode.child[v]
	}

	if !curNode.flag {
		this.size++
		curNode.flag = true
	}
	curNode.value = value
	curNode.word = key
}

// 删除一个词(key)，删除成功返回true
func (this *Trie) Delete(key string) bool {
	this.Lock()
	defer this.Unlock()

	curNode := this.Root
	preNode := this.Root
	var ru rune
	for _, v := range key {
		if curNode.child[v] == nil {
			return false
		}
		preNode = curNode
		curNode = curNode.child[v]
		ru = v
	}

	// 若是叶子节点，则真正删除，否则懒惰删除
	if len(curNode.child) == 0 {
		delete(preNode.child, ru)
	} else {
		curNode.flag = false
		curNode.value = nil
	}

	this.size--

	return true
}

// 查找一个词(key)
// flag为true时表示存在该词(key)， 否则表示不存在该词
// value为该词(key)所在节点的value值
// index表示该词(key)所在路径上最长的一个词的最后一个rune的末尾位置，
// 比如：词典中有一个词："你好"，key为"你好"，那么index为6
// 比如：词典中有一个词："世界"，key为"世界你好"，那么index为6
// 比如：词典中有一个词："hello"，key为"helloworld"，那么index为5
// 比如：词典中有一个词："helloworld"，key为"hello"，那么index为0
func (this *Trie) Find(key string) (flag bool, value interface{}, index int) {
	node, i := this.findNode(key)
	if node == nil {
		return false, nil, i
	} else {
		return node.flag, node.value, i
	}
}

// 功能跟Find完全一样，只是参数在外面就已经将string拆分成了[]rune
// 返回的index为[]rune的下标：
// 比如：词典中有一个词："世界"，key为"世界你好"，那么index为2
func (this *Trie) FindByRunes(runes []rune) (flag bool, value interface{}, index int) {
	node, i := this.findNodeByRunes(runes)
	if node == nil {
		return false, nil, i
	} else {
		return node.flag, node.value, i
	}
}

// 匹配出所有前缀为key的词所在节点的value值
func (this *Trie) PrefixMatch(key string) []*Val {
	this.RLock()
	defer this.RUnlock()

	node, _ := this.findNode(key)
	if node != nil {
		return this.Walk(node)
	}
	return []*Val{}
}

// 功能跟PrefixMatch完全一样，只是参数在外面就已经将string拆分成了[]rune
func (this *Trie) PrefixMatchByRunes(runes []rune) []*Val {
	node, _ := this.findNodeByRunes(runes)
	if node != nil {
		return this.Walk(node)
	}
	return []*Val{}
}

// 遍历
/*func (this *Trie) walk(node *trieNode) (ret []interface{}) {
	if node.flag {
		ret = append(ret, node.value)
	}
	for _, v := range node.child {
		ret = append(ret, this.walk(v)...)
	}
	return
}*/

func (this *Trie) Walk(node *trieNode) (ret []*Val) {
	if node.flag {
		var val = Val{
			Word:  node.word,
			Value: node.value,
		}
		ret = append(ret, &val)
	}
	for _, v := range node.child {
		ret = append(ret, this.Walk(v)...)
	}
	return
}

// 查找一个词(key)所在的节点
// node不为空且node.flag为true时表示存在该词(key)， 否则表示不存在该词
// index表示该词(key)所在路径上最长的一个词的最后一个rune的末尾位置，
// 比如：词典中有一个词："你好"，key为"你好"，那么index为6
// 比如：词典中有一个词："世界"，key为"世界你好"，那么index为6
// 比如：词典中有一个词："hello"，key为"helloworld"，那么index为5
// 比如：词典中有一个词："helloworld"，key为"hello"，那么index为0
func (this *Trie) findNode(key string) (node *trieNode, index int) {
	curNode := this.Root
	ff := false
	for k, v := range key {
		if ff {
			index = k
			ff = false
		}
		if curNode.child[v] == nil {
			return nil, index
		}
		curNode = curNode.child[v]
		if curNode.flag {
			ff = true
		}
	}

	if curNode.flag {
		index = len(key)
	}

	return curNode, index
}

// 功能跟findNode完全一样，只是参数在外面就已经将string拆分成了[]rune
// 返回的index为[]rune的下标：
// 比如：词典中有一个词："世界"，key为"世界你好"，那么index为2
func (this *Trie) findNodeByRunes(runes []rune) (node *trieNode, index int) {
	curNode := this.Root
	ff := false
	for k, v := range runes {
		if ff {
			index = k
			ff = false
		}
		if curNode.child[v] == nil {
			return nil, index
		}
		curNode = curNode.child[v]
		if curNode.flag {
			ff = true
		}
	}

	if curNode.flag {
		index = len(runes)
	}

	return curNode, index
}
