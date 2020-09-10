package main

import (
	"github.com/gin-gonic/gin"
	"github.com/threeq/xid"
	"net/http"
	"strconv"
)

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
