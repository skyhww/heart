package extend

import (
	"github.com/huichen/sego"
	"os"
)

type SegmentsFilter struct {
	Tree
	Seg *sego.Segmenter
}

func (segmentsFilter *SegmentsFilter) Match(word string) (bool, string) {
	keys := sego.SegmentsToSlice(segmentsFilter.Seg.Segment([]byte(word)), true)
	for i := 0; i < len(keys); i++ {
		if segmentsFilter.Tree.match(keys[i]) {
			return true, keys[i]
		}
	}
	return segmentsFilter.Tree.match(word), word
}

func NewSegmentsFilter() *SegmentsFilter {
	s := &SegmentsFilter{Seg: &sego.Segmenter{}}
	f, err := os.Open("extend/fd.txt")
	if err != nil {
		panic(err)
	}
	s.LoadFile(f)
	f, err = os.Open("extend/tf.txt")
	if err != nil {
		panic(err)
	}
	s.LoadFile(f)
	s.Seg.LoadDictionary("extend/dictionary.txt")
	return s
}
