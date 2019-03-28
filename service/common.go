package service

import JSON "encoding/json"

type Json interface {
	toJsonString() string
}

type Page struct {
	PageSize  int
	PageNo    int
	PageCount int
}

type Info struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var Success = &Info{"000000", ""}

func (e *Info) toJsonString() string {
	b, _ := JSON.Marshal(e)
	return string(b)
}

func (e *Info) IsSuccess() bool {
	return e == Success
}
