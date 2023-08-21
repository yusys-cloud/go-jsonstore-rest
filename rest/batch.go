// Author: yangzq80@gmail.com
// Date: 2021/3/21
package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// /api/batch?b=b&k=k
//
// 批量保存[json]数组，如果有 id 则修改，否则新增
func (s *Storage) batchSave(c *gin.Context) {
	b := c.Query("b")
	k := c.Query("k")

	resp := make(map[string]interface{})

	var data []interface{}
	if err := c.ShouldBindJSON(&data); err == nil {
		for _, o := range data {
			s.Create(b, k, o)
		}
	} else {
		resp["json数组格式错误"] = err.Error()
	}
	resp["total"] = len(data)

	c.JSON(http.StatusOK, resp)
}
