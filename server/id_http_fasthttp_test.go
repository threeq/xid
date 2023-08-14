package main_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
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

func testRequest(t *testing.T) {
	n := 20
	size := 10000
	w := sync.WaitGroup{}
	w.Add(n)
	q := make(chan string, 10000)
	for i := 0; i < n; i++ {
		go func() {
			defer w.Done()
			j := 0
			for j < size {
				res, err := http.Get("http://127.0.0.1:8888/xid/snake")
				if err != nil {
					continue
				}
				data, err := io.ReadAll(res.Body)
				if err != nil {
					continue
				}
				q <- string(data)
				j++
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		f, err := os.OpenFile("/Users/three3q/workspaces/myself/xid/out/ids.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		for d := range q {
			_, e := io.WriteString(f, d+"\n")
			if e != nil {
				fmt.Println(e)
			}
		}
		e := f.Close()
		if e != nil {
			fmt.Println(e)
		}
	}()

	w.Wait()
	for {
		if len(q) == 0 {
			close(q)
			break
		}
	}
	<-done
}
