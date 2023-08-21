// Author: yangzq80@gmail.com
// Date: 2021-11-24
package rest

import (
	"fmt"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"log"
	"testing"
)

var storage = NewStorage("../json-db")
var bucket = "test"
var key = "test"

func TestStorage_SearchStruct(t *testing.T) {
	s := Search{}
	s.B = bucket
	s.K = key
	s.Key = "v.Name"
	s.Value = "joy"

	//t:=storage.SearchStruct(s, Test{}).Data.Items.(map[string]interface{})[0].(Test)

}

func TestStorage_Create(t *testing.T) {
	storage.Create(bucket, key, Test{"1", "joy"})
}

func TestStorage_ReadAll(t *testing.T) {
	list := storage.ReadAll(bucket, key).Items.([]model.Data)
	for _, v := range list {
		fmt.Println(v.K, v.V)
	}
}

func TestStorage_DeleteAll(t *testing.T) {
	storage.DeleteAll(bucket, key)
}

type Test struct {
	Id   string
	Name string
}

func TestStorage_CacheGetString(t *testing.T) {

	storage.CachePutString("c-s", "s-path2", "93")

	if storage.CacheGetString("c-s", "s-path2") != "93" {
		t.Error("CacheGetString error,put 91 ")
	}
	log.Println(storage.CacheGetString("c-s", "s-path"))
}
