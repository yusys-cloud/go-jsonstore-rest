// Author: yangzq80@gmail.com
// Date: 2021-02-02
//
package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/yusys-cloud/go-jsonstore-rest/internal"
	"net/http"
)

func main() {

	path := flag.String("path", "./json-db", "-path=./json-db")

	port := flag.String("port", "9999", "-port=9999")

	flag.Parse()

	r := gin.Default()

	internal.NewJsonStoreRest(*path, r)

	r.Use(DisableCors())

	r.Run(":" + *port)
}

//Needed in order to disable CORS for local development
func DisableCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
