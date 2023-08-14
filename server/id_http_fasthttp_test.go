package main_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/threeq/xid"
)

func init() {
	xid.Options(xid.RunTypes("snake+id14"))
	xid.Init()
}

func idHttpFunc(writer http.ResponseWriter, request *http.Request) {

	gen := request.FormValue("gen")
	id := xid.GetID(xid.ID_Snake, gen)

	writer.WriteHeader(200)
	_, _ = fmt.Fprint(writer, strconv.FormatInt(id, 10))
}

func Test_routers(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	idHttpFunc(w, r)
	id, _ := strconv.ParseUint(w.Body.String(), 10, 64)

	log.Println(id)

}

func Benchmark_routers(b *testing.B) {

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
