package acam

import (
	"container/list"
	"fmt"
	"time"
	"unicode/utf8"

	. "graduate/platform"

	"gorm.io/gorm"
)

type trieNode struct {
	val      int64
	fail     *trieNode
	children map[rune]*trieNode
	isEndPos bool
	word     string
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
		isEndPos: false,
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
	now.isEndPos = false
}

func (matcher *Matcher) insert(word WordVal) {
	curNode := matcher.root
	for _, v := range word.WordName {
		if curNode.children[v] == nil {
			curNode.children[v] = newTrieNode()
		}
		curNode = curNode.children[v]
	}
	curNode.val = int64(word.WordVal)
	matcher.size += 1
	curNode.isEndPos = true
	curNode.word = word.WordName
}

func (matcher *Matcher) batchInsert(words []WordVal) {
	for _, val := range words {
		matcher.insert(val)
	}
}

func (matcher *Matcher) build() {
	//BFS构建AC自动机fail树
	que := list.New()
	que.PushBack(matcher.root)

	for que.Len() > 0 {
		tmp := que.Remove(que.Front()).(*trieNode)
		var p *trieNode = nil
		for i, v := range tmp.children {
			if tmp == matcher.root {
				v.fail = tmp
			} else {
				p = tmp.fail
				cnt := 0
				//常数k次 k < 5
				for p != nil {
					if p.children[i] != nil {
						v.fail = p.children[i]
						break
					}
					p = p.fail
					cnt++
				}
				if p == nil {
					v.fail = matcher.root
				}
				if cnt > 10 {
					fmt.Println(cnt)
				}
			}
			que.PushBack(v)
		}
	}
}

func (matcher *Matcher) ReBuild(db *gorm.DB) {

	//删除原数据
	dfsRefresh(matcher.root)
	var id int64 = 0
	sum := 0
	word := make([]WordVal, 0, 1000)

	//拉取新数据
	fmt.Println("start rebuild: ", time.Now())
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
		sum += len(word)
		fmt.Println("get words : ", len(word))
		matcher.batchInsert(word)
	}

	//构建AC自动机

	matcher.build()
	fmt.Println("sum is : ", sum)
	fmt.Println("rebuild finish ,", time.Now())
}

func (matcher *Matcher) FuncDynamicProgrammingAndDivide(text string) {
	dp := make([]uint64, len(text)+10)
	pre := make([]int, len(text)+10)
	curNode := matcher.root
	for idx, _ := range dp {
		dp[idx] = 0
	}
	for idx, val := range text {
		tmp := curNode
		for tmp != nil {
			if tmp.children[val] != nil {
				curNode = tmp.children[val]
				break
			} else {
				tmp = tmp.fail
			}
		}
		if tmp == nil {
			curNode = matcher.root
		}
		cal := curNode
		for cal.isEndPos {
			lenth := utf8.RuneCountInString(curNode.word)
			last := idx - lenth
			if last < 0 {
				if dp[idx] < uint64(cal.val) {
					dp[idx] = uint64(cal.val)
					pre[idx] = -1
				}
			} else {
				if dp[idx] < uint64(cal.val)+dp[last] {
					dp[idx] = uint64(cal.val) + dp[last]
					pre[idx] = last
				}
			}
			cal = cal.fail
		}
	}
}
