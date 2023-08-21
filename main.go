// Author: yangzq80@gmail.com
// Date: 2021-02-02
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/yusys-cloud/go-jsonstore-rest/rest"
	"strings"
)

func main() {

	path := flag.String("path", "./data", "--path=./data")

	port := flag.String("port", "9999", "--port=9999")

	basicAuth := flag.String("basicAuth", "", "--basicAuth=admin:admin")

	flag.Parse()

	r := gin.Default()

	//r.Use(DisableCors())

	s := rest.NewJsonStoreRest(*path)
	s.DisableCors = true
	if *basicAuth != "" {
		s.BasicAuth = parseStringToMap(*basicAuth)
	}
	s.ConfigHandles(r)

	r.Run(":" + *port)
}

func parseStringToMap(inputStr string) map[string]string {
	result := make(map[string]string)

	// 通过 ":" 符号分割字符串
	parts := strings.Split(inputStr, ":")

	if len(parts) == 2 {
		key := parts[0]
		value := parts[1]
		// 将键值对添加到结果的 map 中
		result[key] = value
	}

	return result
}
