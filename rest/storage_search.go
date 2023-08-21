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
}

func (s *Storage) Search(search Search) *model.Response {
	resp := model.NewResponse()

	all := s.ReadAll(search.B, search.K).Items
	b, _ := json.Marshal(all)

	jq := gojsonq.New().FromString(string(b))
	if search.Node != "" {
		jq.From(search.Node)
	}
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
				if search.Relation == "isNotEq" {
					jq.WhereNotEqual(search.Key, search.Value)
				} else if search.Relation == "like" {
					jq.WhereContains(search.Key, search.Value)
				} else {
					jq.WhereEqual(search.Key, search.Value)
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
