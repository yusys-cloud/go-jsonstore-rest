// Author: yangzq80@gmail.com
// Date: 2021-09-10
package model

//type RespData struct {
//	Total int         `json:"total,omitempty"`
//	Items interface{} `json:"items"`
//}

type Response struct {
	//Code int       `json:"code"`
	//Data *RespData `json:"data"`
	Total int         `json:"total,omitempty"`
	Items interface{} `json:"items"`
}

// 将key的值存放到value中id字段，规范前端使用
func (r *Response) FormatKV() *Response {
	switch t := r.Items.(type) {
	case []interface{}:
		items := r.Items.([]interface{})
		for i, o := range items {
			m := o.(map[string]interface{})
			mv := m["v"].(map[string]interface{})
			mv["id"] = m["k"]
			items[i] = mv
		}
		r.Items = items
	case interface{}:
		m := t.(*Data)
		if m == nil {
			r.Items = "无记录"
			return r
		}
		mv := m.V.(map[string]interface{})
		mv["id"] = m.K
		r.Items = mv
	}
	return r
}

func NewResponse() *Response {
	return &Response{}
}

func NewResponseData(data interface{}) *Response {
	return &Response{0, data}
}

func ResponseError(err string) *Response {
	resp := NewResponse()
	//resp.Code = http.StatusBadRequest
	resp.Items = err
	return resp
}

func ResponseOne(item interface{}) *Response {
	resp := NewResponse()
	resp.Items = item
	return resp
}

type Data struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
}
