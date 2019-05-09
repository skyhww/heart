package extend

import (
	"testing"
	"os"
	"fmt"
)

func TestAddWord(t *testing.T) {

/*	  var segmenter sego.Segmenter
	segmenter.LoadDictionary("dictionary.txt")

	// 分词
	text := []byte("中共退党")
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	strs:=sego.SegmentsToSlice(segments,true)*/
	f,err:=os.Open("fd.txt")
	if err!=nil{
		t.Error(err)
	}
	tree := &Tree{}
	tree.LoadFile(f)
	f,err=os.Open("tf.txt")
	if err!=nil{
		t.Error(err)
	}
	tree.LoadFile(f)
	fmt.Println( tree.match("廖伯年"))
}
