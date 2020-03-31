package service

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v7"
	"heart/service/common"
	"log"
	"os"
	"strconv"
)

type ElasticSearchService interface {
	Query(keyword map[string]interface{}, index string, typ string, page *base.Page) ([][]byte, error)
}

func NewElasticSearchService(host string) (ElasticSearchService, error) {
	c, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		return nil, err
	}
	return &SimpleElasticSearchService{client: c}, nil
}

type SimpleElasticSearchService struct {
	client *elastic.Client
}

func (elasticSearchService *SimpleElasticSearchService) Query(keyword map[string]interface{}, index string, typ string, page *base.Page) ([][]byte, error) {
	s := elasticSearchService.client.Search()
	bq := elastic.NewBoolQuery()
	for k, v := range keyword {
		switch v.(type) {
		case string:
			if v.(string) != "" {
				bq.Must(elastic.NewMatchQuery(k, v))
			}
		default:
			bq.Must(elastic.NewTermQuery(k, v))
		}
	}
	s.Index(index).Query(bq).From(page.PageSize * (page.PageNo - 1)).Size(page.PageSize).SortBy(elastic.NewScoreSort())
	r, err := s.Do(context.Background())
	if err != nil {
		fmt.Sprint(err)
		return nil, err
	}
	strInt64 := strconv.FormatInt(r.Hits.TotalHits.Value, 10)
	page.Count, _ = strconv.Atoi(strInt64)
	if r.Hits.TotalHits.Value > 0 {
		b := make([][]byte, r.Hits.TotalHits.Value)
		for i, hit := range r.Hits.Hits {
			b[i] = hit.Source
		}
		return b, nil
	}
	return nil, nil
}
