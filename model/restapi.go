// Author: yangzq80@gmail.com
// Date: 2021-09-10
//
package model

type Data struct {
	Total int
	Items interface{}
}

type Response struct {
	Code string
	Data *Data
}

func NewResponse() *Response {
	return &Response{"20000", &Data{}}
}
