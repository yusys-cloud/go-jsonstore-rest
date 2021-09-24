// Author: yangzq80@gmail.com
// Date: 2021-03-25
//
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"net/http"
)

type JsonStoreRest struct {
	D *Storage
	//gin *gin.Context
}

func NewJsonStoreRest(dir string, r *gin.Engine) *JsonStoreRest {

	s := NewStorage(dir)

	s.ConfigHandles(r)

	return &JsonStoreRest{s}
}

func (s *Storage) ConfigHandles(r *gin.Engine) {
	rg := r.Group("/api/kv")
	rg.POST("/:b/:k", s.create)
	rg.GET("/:b/:k", s.readAll)
	rg.GET("/:b/:k/:kid", s.read)
	rg.PUT("/:b/:k/:kid", s.update)
	rg.PUT("/:b/:k/:kid/weight", s.updateWeight)
	rg.DELETE("/:b/:k/:kid", s.delete)
	rg.DELETE("/:b/:k", s.deleteAll)
	//Search
	r.GET("/api/search", s.search)
	r.DELETE("/api/search", s.deleteSearch)
}

func (s *Storage) create(c *gin.Context) {

	var data interface{}
	resp := model.NewResponse()

	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Data.Items = "BindError:" + err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data.Total = 1
	resp.Data.Items = s.Create(c.Param("b"), c.Param("k"), data)
	c.JSON(http.StatusOK, resp)
}

func (s *Storage) update(c *gin.Context) {
	var data interface{}
	resp := model.NewResponse()

	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		resp.Code = http.StatusBadRequest
		resp.Data.Items = err.Error()
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Data.Total = 1
	resp.Data.Items = s.Update(c.Param("b"), c.Param("kid"), data)
	c.JSON(http.StatusOK, resp)
}

func (s *Storage) readAll(c *gin.Context) {
	b := s.ReadAllSort(c.Param("b"), c.Param("k"))
	c.JSON(http.StatusOK, b)
}

func (s *Storage) search(c *gin.Context) {
	var search Search
	c.ShouldBind(&search)
	c.JSON(http.StatusOK, s.Search(search))
}

//根据搜索内容删除
func (s *Storage) deleteSearch(c *gin.Context) {
	var search Search
	c.ShouldBind(&search)
	c.JSON(http.StatusOK, model.NewResponseData(s.DeleteList(search.B, s.Search(search).Data.Items, false)))
}

func (s *Storage) updateWeight(c *gin.Context) {
	c.JSON(http.StatusOK, model.NewResponseData(s.UpdateWeight(c.Param("b"), c.Param("kid"))))
}

func (s *Storage) read(c *gin.Context) {

	c.JSON(http.StatusOK, model.NewResponseData(s.Read(c.Param("b"), c.Param("kid"))))
}

func (s *Storage) delete(c *gin.Context) {

	s.Delete(c.Param("b"), c.Param("kid"))

	c.JSON(http.StatusOK, model.NewResponseData("success"))
}
func (s *Storage) deleteAll(c *gin.Context) {

	i := s.DeleteAll(c.Param("b"), c.Param("k"))

	resp := model.NewResponseData("success")
	resp.Data.Total = i

	c.JSON(http.StatusOK, resp)
}
