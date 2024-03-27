// Author: yangzq80@gmail.com
// Date: 2021-03-16
package rest

import (
	"encoding/json"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"strings"
)

type Search struct {
	B        string `form:"b"`
	K        string `form:"k"`
	Node     string `form:"node"`
	Key      string `form:"key"`      // Search conditions key
	Value    string `form:"value"`    // Search conditions value
	Relation string `form:"relation"` // Search relation,default equal; equal,like,isNotEq
	ShortBy  string `form:"shortBy"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
	Fields   string `form:"fields"`
}

func (s *Storage) Search(search Search) *model.Response {
	resp := model.NewResponse()

	all := s.ReadAll(search.B, search.K).Items
	b, _ := json.Marshal(all)

	jq := gojsonq.New().FromString(string(b))
	if search.Node != "" {
		jq.From(search.Node)
	}
	isQueryNil := false
	if search.Value != "" && search.Key != "" {
		//Multiple conditions dynamic search
		if strings.Contains(search.Key, ",") {
			ks := strings.Split(search.Key, ",")
			vs := strings.Split(search.Value, ",")
			for i := 0; i < len(ks); i++ {
				if ks[i] != "" && vs[i] != "" {
					if strings.Contains(vs[i], "|") {
						jq.WhereIn(ks[i], strings.Split(vs[i], "|"))
					} else {
						if search.Relation == "like" {
							jq.WhereContains(ks[i], vs[i])
						} else {
							jq.WhereEqual(ks[i], vs[i])
						}
					}
				}
			}
		} else {
			if strings.Contains(search.Value, "|") {
				jq.WhereIn(search.Key, strings.Split(search.Value, "|"))
			} else {
				switch search.Relation {
				case "isNotEq":
					if search.Value == "nil" {
						jq.WhereNotNil(search.Key)
					} else {
						jq.WhereNotEqual(search.Key, search.Value)
					}
				case "like":
					jq.WhereContains(search.Key, search.Value)
				default:
					// 如果查询 nil = 全量 - 非nil
					if search.Value == "nil" {
						jq.WhereNotNil(search.Key)
						isQueryNil = true
					} else {
						jq.WhereEqual(search.Key, search.Value)
					}
				}
			}
		}
	}

	if search.ShortBy != "" {
		sts := strings.Split(search.ShortBy, ",")
		jq.SortBy(sts[0], sts[1])
	} else {
		jq.SortBy("k", "desc")
	}
	resp.Total = len(jq.Get().([]interface{}))
	// Offset and limit
	if search.Page != 0 {
		jq.Offset((search.Page - 1) * search.Size)
	}
	// limit
	if search.Size != 0 {
		jq.Limit(search.Size)
	}

	resp.Items = jq.Get().([]interface{})

	// 如果查询字段为nil，则通过[全量-非nil]来实现
	if isQueryNil {
		resp = substrResult(resp, jq)
	}
	var fs []string
	if search.Fields != "" {
		fs = strings.Split(search.Fields, ",")
	}
	return resp.FormatFields(fs)
}

// 从全量中排除Response中目标
func substrResult(resp *model.Response, jq *gojsonq.JSONQ) *model.Response {
	var result []interface{}
	if notNilArr, ok := resp.Items.([]interface{}); ok {
		if allArr, ok := jq.Reset().Get().([]interface{}); ok {
			for _, a := range allArr {
				isNil := true
				// 根据非nil与全量对比出 nil 的
				for _, n := range notNilArr {
					if a.(map[string]interface{})["k"].(string) == n.(map[string]interface{})["k"].(string) {
						isNil = false
						break
					}
				}
				if isNil {
					result = append(result, a)
				}
			}
		}
		resp.Items = result
		resp.Total = len(result)
	}
	return resp
}

func (s *Storage) SearchStruct(search Search, obj interface{}) *model.Response {

	rs := s.Search(search)

	for _, item := range rs.Items.([]interface{}) {
		in := item.(map[string]interface{})["v"].(map[string]interface{})
		jsonbody, _ := json.Marshal(in)
		json.Unmarshal(jsonbody, &obj)
		item.(map[string]interface{})["v"] = obj
	}

	return rs
}
