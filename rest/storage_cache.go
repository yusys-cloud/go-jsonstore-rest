// Author: yangzq80@gmail.com
// Date: 2023/2/28
package rest

import "encoding/json"

// 通用 key：value 对象存储
func (s *Storage) CachePut(key string, val interface{}) {
	s.bucket(CACHE_BUCKET).Set("b-c-"+key, val)
	s.savePersistent(CACHE_BUCKET)
}
func (s *Storage) CacheGet(key string, val *interface{}) {
	var rs interface{}
	s.bucket(CACHE_BUCKET).Get("b-c-"+key, &rs)

	if rs != nil {
		body, _ := json.Marshal(rs)
		json.Unmarshal(body, val)
	}
}

// 通用 key：value 字符串类型kv存储
func (s *Storage) CachePutString(category string, kvKey string, kvValue string) {

	var vObj interface{}
	s.CacheGet(category, &vObj)
	val := make(map[string]interface{})
	if vObj != nil {
		val = vObj.(map[string]interface{})
	}
	val[kvKey] = kvValue
	s.CachePut(category, val)
}
func (s *Storage) CacheGetString(category string, kvKey string) string {
	var vObj interface{}
	s.CacheGet(category, &vObj)
	if vObj == nil {
		return ""
	}
	kv := vObj.(map[string]interface{})
	if kv == nil || kv[kvKey] == nil {
		return ""
	}
	return kv[kvKey].(string)
}

func (s *Storage) FIFO(key string, val interface{}, size int) {
	resp := s.ReadAllSort(CACHE_BUCKET, key)
	if resp.Total >= size {
		s.Delete(CACHE_BUCKET, resp.Items.([]interface{})[size-1].(map[string]interface{})["k"].(string))
	}
	s.Create(CACHE_BUCKET, key, val)
}
