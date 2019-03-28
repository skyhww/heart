package service

type Json interface {
	toJsonString() string
}

type Page struct {
	pageSize  int
	pageNo    int
	pageCount int
}
