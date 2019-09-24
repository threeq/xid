package main

import (
	"github.com/gin-gonic/gin"
	"github.com/threeq/xid"
	"log"
	"net/http"
	"strconv"
)

func init() {
	log.Println("初始化 web 服务 ...")
	webapp = newApp()
}

func newApp() http.Handler {
	//http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	gen := request.FormValue("gen")
	//	id := xid.MultiIdGenerator(gen)
	//
	//	writer.WriteHeader(200)
	//	_, _ = fmt.Fprint(writer, strconv.FormatInt(id, 10))
	//})
	//
	//graceShutdownServe(addr, http.DefaultServeMux)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/", func(c *gin.Context) {
		gen := c.Query("gen")
		id := xid.MultiIdGenerator(gen).Next()
		c.String(200, strconv.FormatInt(id, 10))
	})

	engine.GET("/empty", func(c *gin.Context) {
		c.String(200, "empty")
	})

	return engine

}
