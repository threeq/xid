package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/threeq/xid"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func idHttpFunc(writer http.ResponseWriter, request *http.Request) {

	gen := request.FormValue("gen")
	id := xid.MultiIdGenerator(gen).Next()

	writer.WriteHeader(200)
	_, _ = fmt.Fprint(writer, strconv.FormatInt(id, 10))
}

func Test_routers(t *testing.T) {
	xid.Config("id",xid.NewNodeAllocationSingle())

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	idHttpFunc(w, r)
	id, _ := strconv.ParseUint(w.Body.String(), 10, 64)

	log.Println(id)

}

func Benchmark_routers(b *testing.B) {
	xid.Config("id",xid.NewNodeAllocationSingle())

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		idHttpFunc(w, r)
		id, _ := strconv.ParseUint(w.Body.String(), 10, 64)

		if id < 1000 {
			b.Fatalf("id 生产错误: %d", id)
		}
	}
}

func Test_gin(t *testing.T)  {
	xid.Config("id14",xid.NewNodeAllocationSingle())

	r := httptest.NewRequest(http.MethodGet, "/id14", nil)
	w := httptest.NewRecorder()

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/id14", func(c *gin.Context) {
		gen := c.Query("gen")
		id := xid.MultiIdGenerator(gen).Next()
		c.String(200, strconv.FormatInt(id, 10))
	})

	engine.ServeHTTP(w, r)

	id, _ := strconv.ParseUint(w.Body.String(), 10, 64)

	log.Println(id, fmt.Sprintf("len=%d", len(fmt.Sprintf("%d", id))))
	assert.Equal(t, len(fmt.Sprintf("%d", id)), 14)
}

func Benchmark_gin(b *testing.B) {
	xid.Config("id",xid.NewNodeAllocationSingle())

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/", func(c *gin.Context) {
		gen := c.Query("gen")
		id := xid.MultiIdGenerator(gen).Next()
		c.String(200, strconv.FormatInt(id, 10))
	})

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, r)
		id, _ := strconv.ParseUint(w.Body.String(), 10, 64)

		if id < 1000 {
			b.Fatalf("id 生产错误: %d", id)
		}
	}
}