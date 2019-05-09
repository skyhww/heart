package extend

import (
	"os"
	"bufio"
	"io"
)

/**
* 字典树
 */
type Tree struct {
	Bucket  map[rune]*Tree
	Include bool
}

func (tree *Tree) AddWord(text string) {
	r := []rune(text)
	t := tree
	for i := 0; i < len(r); i++ {
		if t.Bucket == nil {
			t.Bucket = make(map[rune]*Tree)
		}
		tmp := t.Bucket[r[i]]
		if tmp == nil {
			bucket := make(map[rune]*Tree)
			t.Bucket[r[i]] = &Tree{Bucket: bucket}
		} else if tmp.Include {
			return
		}
		t = t.Bucket[r[i]]
	}
	t.Include = true
}

func (tree *Tree) LoadFile(file *os.File) {
	rd := bufio.NewReader(file)
	defer file.Close()
	text, _, err := rd.ReadLine()
	for err == nil || io.EOF != err {
		tree.AddWord(string(text))
		text, _, err = rd.ReadLine()
	}
}
func (tree *Tree) Match(word string) bool {
	r := []rune(word)
	for i := 0; i < len(r); i++ {
		if tree.match(string(r[i:])) {
			return true
		}
	}
	return false
}

func (tree *Tree) match(word string) bool {
	r := []rune(word)
	t := tree
	for i := 0; i < len(r) && len(t.Bucket) > 0; i++ {
		t = t.Bucket[r[i]]
		if t != nil && t.Include {
			return true
		} else if t == nil {
			return false
		}
	}
	return false
}
