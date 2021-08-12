package service

import (
	"fmt"
	"gopkg.in/olivere/elastic.v7"
	base "heart/service/common"
	"log"
	"os"
	"testing"
)

func TestSimpleElasticSearchService_Query(t *testing.T) {
	c, err := elastic.NewClient(elastic.SetURL("http://149.28.77.66:9200"), elastic.SetSniff(false), elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		fmt.Println(err)
	}
	s := &SimpleElasticSearchService{c}
	m := make(map[string]interface{})
	m["content"] = ""
	m["enable"] = true
	m["user_id"] = 1
	p := &base.Page{PageNo: 1, PageSize: 10}
	s.Query(m, "3dheart_posts", "posts", p)
}
