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
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Success = &Info{Code: "000000", Message: ""}

func NewSuccess(data interface{}) *Info {
	return &Info{Success.Code, Success.Message, data}
}

func (e *Info) toJsonString() string {
	b, _ := JSON.Marshal(e)
	return string(b)
}

func (e *Info) IsSuccess() bool {
	return e.Code == Success.Code
}
