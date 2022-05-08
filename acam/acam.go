package acam

import (
	"fmt"

	. "graduate/platform"

	"gorm.io/gorm"
)

type trieNode struct {
	val      int64
	fail     *trieNode
	children map[rune]*trieNode
}

type Matcher struct {
	root *trieNode
	size int64
}

func NewMatcher() *Matcher {
	return &Matcher{
		root: newTrieNode(),
		size: 0,
	}
}

func newTrieNode() *trieNode {
	return &trieNode{
		val:      0,
		fail:     nil,
		children: make(map[rune]*trieNode),
	}
}

func dfsRefresh(now *trieNode) {
	if now == nil {
		return
	}
	for _, val := range now.children {
		dfsRefresh(val)
	}
	now.children = make(map[rune]*trieNode)
	now.fail = nil
	now.val = 0
}

func (matcher *Matcher) ReBuild(db *gorm.DB) {
	dfsRefresh(matcher.root)
	var id int64 = 0
	sum := 0
	word := make([]WordVal, 1000)
	for {
		err := db.Model(&WordVal{}).Limit(1000).Find(&word, "word_id > ?", id).Error
		if err != nil {
			panic(err)
		}
		if len(word) != 0 {
			id = word[len(word)-1].WordId
		} else {
			break
		}
		fmt.Println(len(word))
		sum += len(word)
	}
	fmt.Println("sum is : ", sum)
}
