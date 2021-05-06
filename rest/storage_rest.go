// Author: yangzq80@gmail.com
// Date: 2021-03-25
//
package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type JsonStoreRest struct {
	db *Storage
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
}

func (s *Storage) create(c *gin.Context) {

	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := s.Create(c.Param("b"), c.Param("k"), data)

	c.JSON(http.StatusOK, id)
}

func (s *Storage) update(c *gin.Context) {
	var data interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.Update(c.Param("b"), c.Param("kid"), data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
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

func (s *Storage) updateWeight(c *gin.Context) {
	c.JSON(http.StatusOK, s.UpdateWeight(c.Param("b"), c.Param("kid")))
}

func (s *Storage) read(c *gin.Context) {

	kv := s.Read(c.Param("b"), c.Param("kid"))

	c.JSON(http.StatusOK, kv)
}

func (s *Storage) delete(c *gin.Context) {

	s.Delete(c.Param("b"), c.Param("kid"))

	c.JSON(http.StatusOK, "ok")
}
func (s *Storage) deleteAll(c *gin.Context) {

	i := s.DeleteAll(c.Param("b"), c.Param("k"))

	c.JSON(http.StatusOK, "Delete nums:"+strconv.Itoa(i))
}
