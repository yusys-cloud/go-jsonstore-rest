// Author: yangzq80@gmail.com
// Date: 2021-09-02
package rest

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"io"
	"net/http"
	"testing"
)

var url = "http://localhost:9999/kv/meta/node"

func TestSaveOrUpdate(t *testing.T) {

	var jsonStr = `{"id":"node1693552235629252608","ip": "192.168.49.69","name":"redis-n1","dc":"default","lable":"Redis"}`

	// 有id 更新
	resp := doSave(jsonStr)
	if resp["id"] != "node1693552235629252608" {
		t.Error("not updated")
	}
	// 无id 新增
	jsonStr = `{"ip": "192.168.49.69","name":"redis-n1","dc":"default","lable":"Redis"}`
	doSave(jsonStr)
}

func doSave(jsonStr string) map[string]interface{} {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var respm model.Response
	json.Unmarshal(body, &respm)

	return respm.Items.(map[string]interface{})
}
