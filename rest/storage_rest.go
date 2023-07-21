// Author: yangzq80@gmail.com
// Date: 2021-03-25
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yusys-cloud/go-jsonstore-rest/model"
	"net/http"
	"strconv"
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
	rg := r.Group("/kv")
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
	//通用缓存
	r.POST("/api/cache", s.cache)
	r.GET("/api/cache/:key", s.cacheGet)
	r.POST("/api/fifo", s.fifo)
}

func (s *Storage) create(c *gin.Context) {

	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, model.ResponseError("BindError:"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.ResponseOne(s.Create(c.Param("b"), c.Param("k"), data)).FormatKV())
}

func (s *Storage) update(c *gin.Context) {
	var data interface{}
	resp := model.NewResponse()

	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		//resp.Code = http.StatusBadRequest
		resp.Items = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	resp.Items = s.Update(c.Param("b"), c.Param("kid"), data)
	c.JSON(http.StatusOK, resp.FormatKV())
}

func (s *Storage) readAll(c *gin.Context) {
	b := s.ReadAllSort(c.Param("b"), c.Param("k"))
	c.JSON(http.StatusOK, b.FormatKV())
}

func (s *Storage) search(c *gin.Context) {
	var search Search
	c.ShouldBind(&search)
	c.JSON(http.StatusOK, s.Search(search).FormatKV())
}

// 根据搜索内容删除
func (s *Storage) deleteSearch(c *gin.Context) {
	var search Search
	c.ShouldBind(&search)
	c.JSON(http.StatusOK, model.NewResponseData(s.DeleteList(search.B, s.Search(search).Items, false)))
}

func (s *Storage) updateWeight(c *gin.Context) {
	c.JSON(http.StatusOK, model.NewResponseData(s.UpdateWeight(c.Param("b"), c.Param("kid"))))
}

func (s *Storage) read(c *gin.Context) {

	c.JSON(http.StatusOK, model.NewResponseData(s.Read(c.Param("b"), c.Param("kid"))).FormatKV())
}

func (s *Storage) delete(c *gin.Context) {

	s.Delete(c.Param("b"), c.Param("kid"))

	c.JSON(http.StatusOK, model.NewResponseData("success"))
}
func (s *Storage) deleteAll(c *gin.Context) {

	i := s.DeleteAll(c.Param("b"), c.Param("k"))

	resp := model.NewResponseData("success")
	resp.Total = i

	c.JSON(http.StatusOK, resp)
}

// 通用全局key value 缓存到bucket=meta中
func (s *Storage) cache(c *gin.Context) {
	var data model.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, model.ResponseError("BindError:"+err.Error()))
		return
	}

	s.CachePut(data.K, data.V)

	c.JSON(http.StatusOK, model.ResponseOne(data))
}
func (s *Storage) cacheGet(c *gin.Context) {
	var b interface{}
	s.CacheGet(c.Param("key"), &b)
	c.JSON(http.StatusOK, model.NewResponseData(b))
}

func (s *Storage) fifo(c *gin.Context) {
	var data model.Data
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, model.ResponseError("BindError:"+err.Error()))
		return
	}

	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	s.FIFO(data.K, data.V, size)

	c.JSON(http.StatusOK, model.ResponseOne(data))
}
