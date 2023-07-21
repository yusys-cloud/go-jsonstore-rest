// Author: yangzq80@gmail.com
// Date: 2021-03-16
package rest

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/xujiajun/utils/filesystem"
	"github.com/yusys-cloud/go-jsonstore-rest/jsonstore"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Storage struct {
	buckets map[string]*jsonstore.JSONStore
	dir     string
	IdNode  *snowflake.Node
}

type Search struct {
	B        string `form:"b"`
	K        string `form:"k"`
	Node     string `form:"node"`
	Key      string `form:"key"`      // Search conditions key
	Value    string `form:"value"`    // Search conditions value
	Relation string `form:"relation"` // Search relation,default equal; equal,like
	ShortBy  string `form:"shortBy"`
	Page     int    `form:"page"`
	Size     int    `form:"size"`
}

const (
	CACHE_BUCKET string = "meta"
)

func NewStorage(dir string) *Storage {
	log.Println("Init JSON storage...", dir)
	//create dir
	mkdirIfNotExist(dir)

	node, _ := snowflake.NewNode(1)

	return &Storage{buckets: make(map[string]*jsonstore.JSONStore), dir: dir, IdNode: node}
}

func (s *Storage) bucket(bucket string) *jsonstore.JSONStore {
	// From memory
	if s.buckets[bucket] != nil {
		return s.buckets[bucket]
	}
	// From local
	if ss, err := jsonstore.Open(s.getFileName(bucket)); err == nil {
		s.buckets[bucket] = ss
		return s.buckets[bucket]
	}
	// New json storage
	s.buckets[bucket] = new(jsonstore.JSONStore)
	return s.buckets[bucket]
}

// 查询bucket中 key 全部
func (s *Storage) ReadAll(bucket string, key string) *model.Response {

	resp := model.NewResponse()

	rs := s.bucket(bucket).GetAll(regexp.MustCompile(key))

	resp.Total = len(rs)
	resp.Items = convertMapToArray(rs)

	return resp
}
func (s *Storage) ReadAllSort(bucket string, key string) *model.Response {

	resp := model.NewResponse()

	rs := s.bucket(bucket).GetAll(regexp.MustCompile(key))

	resp.Total = len(rs)
	b, _ := json.Marshal(convertMapToArray(rs))

	jq := gojsonq.New().FromString(string(b))
	jq.SortBy("k", "desc")

	resp.Items = jq.Get()
	return resp
}

// 查询单个
func (s *Storage) Read(bucket string, key string) *model.Data {

	_, rs := s.bucket(bucket).GetRawMessage(key)

	if rs == nil {
		return nil
	}

	var f interface{}

	json.Unmarshal(rs, &f)

	return &model.Data{key, f}
}

// 查询单个，返回 Struct 对象
func (s *Storage) ReadOneStruct(bucket string, key string, v interface{}) error {

	error := s.bucket(bucket).Get(key, v)

	return error
}

func (s *Storage) ReadOneRaw(bucket string, key string) []byte {

	_, rs := s.bucket(bucket).GetRawMessage(key)

	return rs
}

// 保存key,value. bucket类似table
func (s *Storage) Create(bucket string, key string, value interface{}) *model.Data {

	//默认自增ID
	id := key + ":" + s.IdNode.Generate().String()

	err := s.bucket(bucket).Set(id, value)
	if err != nil {
		log.Error(err.Error())
	}

	s.savePersistent(bucket)

	return s.Read(bucket, id)
}

// 根据key更新
func (s *Storage) Update(bucket string, key string, value interface{}) *model.Data {

	old := s.Read(bucket, key)
	if old == nil {
		return nil
	}

	err := s.bucket(bucket).Set(key, value)
	if err != nil {
		log.Error(err.Error())
	}

	s.savePersistent(bucket)

	return s.Read(bucket, key)
}
func (s *Storage) UpdateWeight(bucket string, kid string) interface{} {

	d := s.Read(bucket, kid)

	i := d.V.(map[string]interface{})
	i["weight"] = strconv.FormatInt(time.Now().Unix(), 10)

	err := s.bucket(bucket).Set(kid, i)
	if err != nil {
		log.Error(err.Error())
	}

	s.savePersistent(bucket)

	return i
}
func (s *Storage) UpdateMarshalValue(bucket string, key string, value []byte) error {

	err := s.bucket(bucket).SetMarshalValue(key, value)
	if err != nil {
		log.Error(err.Error())
	}

	s.savePersistent(bucket)

	return err
}

// 根据key删除
func (s *Storage) Delete(bucket string, key string) {

	s.bucket(bucket).Delete(key)

	s.savePersistent(bucket)
}
func (s *Storage) DeleteAll(bucket string, key string) int {
	rs := s.ReadAll(bucket, key)
	return s.DeleteList(bucket, rs.Items, true)
}

func (s *Storage) DeleteList(bucket string, items interface{}, isData bool) int {
	n := 0
	if isData {
		for _, value := range items.([]model.Data) {
			s.bucket(bucket).Delete(value.K)
			n++
		}
	} else {
		for _, value := range items.([]interface{}) {
			s.bucket(bucket).Delete(value.(map[string]interface{})["k"].(string))
			n++
		}
	}
	s.savePersistent(bucket)
	return n
}

func (s *Storage) savePersistent(bucket string) {
	// Saving will automatically gzip if .gz is provided
	if err := jsonstore.Save(s.bucket(bucket), s.getFileName(bucket)); err != nil {
		log.Error(err)
	}
}

func (s *Storage) getFileName(bucket string) string {
	return s.dir + "/" + bucket + ".json.gz"
}

func mkdirIfNotExist(rootDir string) error {
	if ok := filesystem.PathIsExist(rootDir); !ok {
		if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func convertMapToArray(raw map[string]json.RawMessage) []model.Data {
	datas := make([]model.Data, 0)
	for k, v := range raw {
		datas = append(datas, model.Data{k, v})
	}
	return datas
}
