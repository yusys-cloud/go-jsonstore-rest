// Author: yangzq80@gmail.com
// Date: 2021-03-16
//
package rest

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/xujiajun/utils/filesystem"
	"github.com/yusys-cloud/go-jsonstore-rest/jsonstore"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Storage struct {
	buckets map[string]*jsonstore.JSONStore
	dir     string
	idNode  *snowflake.Node
}

type Data struct {
	K string      `json:"k"`
	V interface{} `json:"v"`
}
type Search struct {
	B       string `form:"b"`
	K       string `form:"k"`
	Node    string `form:"node"`
	Key     string `form:"key"`
	Value   string `form:"value"`
	ShortBy string `form:"shortBy"`
	Offset  int    `form:"offset"`
	Limit   int    `form:"limit"`
}

func NewStorage(dir string) *Storage {
	log.Println("Init JSON storage...", dir)
	//create dir
	mkdirIfNotExist(dir)

	node, _ := snowflake.NewNode(1)

	return &Storage{buckets: make(map[string]*jsonstore.JSONStore), dir: dir, idNode: node}
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

//查询bucket中 key 全部
func (s *Storage) ReadAll(bucket string, key string) []Data {

	rs := s.bucket(bucket).GetAll(regexp.MustCompile(key))

	return convertMapToArray(rs)
}
func (s *Storage) ReadAllSort(bucket string, key string) interface{} {

	rs := s.bucket(bucket).GetAll(regexp.MustCompile(key))
	b, _ := json.Marshal(convertMapToArray(rs))

	jq := gojsonq.New().FromString(string(b))
	jq.SortBy("k", "desc")
	return jq.Get()
}

//
func (s *Storage) Search(search Search) interface{} {
	all := s.ReadAll(search.B, search.K)
	b, _ := json.Marshal(all)

	jq := gojsonq.New().FromString(string(b))
	if search.Node != "" {
		jq.From(search.Node)
	}
	if search.Value != "" {
		jq.WhereContains(search.Key, search.Value)
	}

	if search.ShortBy != "" {
		sts := strings.Split(search.ShortBy, ",")
		jq.SortBy(sts[0], sts[1])
	} else {
		jq.SortBy("k", "desc")
	}
	// Offset and limit
	if search.Offset != 0 {
		jq.Offset(search.Offset)
	}
	if search.Limit != 0 {
		jq.Limit(search.Limit)
	}

	return jq.Get()
}

//查询单个
func (s *Storage) Read(bucket string, key string) Data {

	_, rs := s.bucket(bucket).GetRawMessage(key)

	var f interface{}

	json.Unmarshal(rs, &f)

	return Data{key, f}
}

//查询单个，返回 Struct 对象
func (s *Storage) ReadOneStruct(bucket string, key string, v interface{}) error {

	error := s.bucket(bucket).Get(key, v)

	return error
}

func (s *Storage) ReadOneRaw(bucket string, key string) []byte {

	_, rs := s.bucket(bucket).GetRawMessage(key)

	return rs
}

//保存key,value. bucket类似table
func (s *Storage) Create(bucket string, key string, value interface{}) string {

	//默认自增ID
	id := key + ":" + s.idNode.Generate().String()

	err := s.bucket(bucket).Set(id, value)
	if err != nil {
		panic(err)
	}

	s.savePersistent(bucket)

	return id
}

// 根据key更新
func (s *Storage) Update(bucket string, key string, value interface{}) error {

	err := s.bucket(bucket).Set(key, value)
	if err != nil {
		panic(err)
	}

	s.savePersistent(bucket)

	return err
}
func (s *Storage) UpdateWeight(bucket string, kid string) interface{} {

	d := s.Read(bucket, kid)

	i := d.V.(map[string]interface{})
	i["weight"] = strconv.FormatInt(time.Now().Unix(), 10)

	err := s.bucket(bucket).Set(kid, i)
	if err != nil {
		panic(err)
	}

	s.savePersistent(bucket)

	return i
}
func (s *Storage) UpdateMarshalValue(bucket string, key string, value []byte) error {

	err := s.bucket(bucket).SetMarshalValue(key, value)
	if err != nil {
		panic(err)
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
	for _, value := range rs {
		s.bucket(bucket).Delete(value.K)
	}
	s.savePersistent(bucket)
	return len(rs)
}

func (s *Storage) savePersistent(bucket string) {
	// Saving will automatically gzip if .gz is provided
	if err := jsonstore.Save(s.bucket(bucket), s.getFileName(bucket)); err != nil {
		log.Error(err)
		panic(err)
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

func convertMapToArray(raw map[string]json.RawMessage) []Data {
	datas := make([]Data, 0)
	for k, v := range raw {
		datas = append(datas, Data{k, v})
	}
	return datas
}
