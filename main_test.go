// Author: yangzq80@gmail.com
// Date: 2021-09-02
package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
)

var url = "http://localhost:9999/kv/meta/node"

func BenchmarkCreate100(b *testing.B) {

	var jsonStr = []byte(`{"ip": "192.168.49.69","name":"redis-n1","dc":"default","lable":"Redis"}`)

	for i := 1; i < 100; i++ {
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		//req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}
		defer resp.Body.Close()
	}
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
func BenchmarkQuery100(b *testing.B) {
	for i := 1; i < 100; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}
