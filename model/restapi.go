// Author: yangzq80@gmail.com
// Date: 2021-09-10
//
package model

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

type Data struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
}
