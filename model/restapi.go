// Author: yangzq80@gmail.com
// Date: 2021-09-10
//
package model

import "net/http"

type RespData struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type Response struct {
	Code int       `json:"code"`
	Data *RespData `json:"data"`
}

func NewResponse() *Response {
	return &Response{20000, &RespData{}}
}

func NewResponseData(data interface{}) *Response {
	return &Response{20000, &RespData{1, data}}
}

func ResponseError(err string) *Response {
	resp := NewResponse()
	resp.Code = http.StatusBadRequest
	resp.Data.Items = err
	return resp
}

func ResponseOne(item interface{}) *Response {
	resp := NewResponse()
	resp.Data.Items = item
	resp.Data.Total = 1
	return resp
}

type Data struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
}
