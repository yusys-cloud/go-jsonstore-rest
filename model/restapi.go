// Author: yangzq80@gmail.com
// Date: 2021-09-10
//
package model

type Data struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type Response struct {
	Code int   `json:"code"`
	Data *Data `json:"data"`
}

func NewResponse() *Response {
	return &Response{20000, &Data{}}
}
