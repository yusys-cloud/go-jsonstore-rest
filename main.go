// Author: yangzq80@gmail.com
// Date: 2021-02-02
//
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/yusys-cloud/go-jsonstore-rest/internal"
)

func main() {

	path := flag.String("path", "./json-db", "-path=./json-db")

	port := flag.String("port", "9999", "-port=9999")

	flag.Parse()

	r := gin.Default()

	internal.NewJsonStoreRest(*path, r)

	r.Run(":" + *port)
}
