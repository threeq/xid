package main

import (
	"github.com/gin-gonic/gin"
	"github.com/threeq/xid"
	"net/http"
	"strconv"
)

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func newIDApp(basePath string) http.Handler {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	router := engine.Group(basePath)
	router.GET("/", func(c *gin.Context) {
		gen := c.Query("gen")
		id := xid.MultiIdGenerator(gen).Next()
		c.String(200, strconv.FormatInt(id, 10))
	})

	router.GET("/batch", func(c *gin.Context) {
		gen := c.Query("gen")
		numStr := c.Query("num")

		var res Response
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			res.Code = 422
			res.Message = err.Error()
			c.JSON(http.StatusOK, res)

			return
		}

		ids := xid.GetIDS(gen, num)

		res.Code = 200
		res.Message = "success"
		res.Data = map[string]interface{}{
			"ids": ids,
		}

		c.JSON(http.StatusOK, res)
	})

	router.GET("/gen/:gen", func(c *gin.Context) {
		gen := c.Param("gen")
		id := xid.MultiIdGenerator(gen).Next()
		c.String(200, strconv.FormatInt(id, 10))
	})

	router.GET("/test/empty", func(c *gin.Context) {
		c.String(200, "empty")
	})

	return engine

}
