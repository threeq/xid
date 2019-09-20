package main

import (
	"context"
	"flag"
	"github.com/threeq/xid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var webapp http.Handler

func main() {
	webAddr := flag.String("web-addr", ":8080", "web 监听地址和端口")
	model := flag.String("model", "single", "运行模式：single 单机 id 生成；redis 使用redis 分布式 id生成")
	redisAddr := flag.String("redis-addr", "localhost:6379", "redis 地址和端口")
	redisPwd := flag.String("redis-pwd", "", "redis 密码")
	nodeBits := flag.Uint("node-bits", 5, "机器长度")
	stepBits := flag.Uint("step-bits", 6, "")
	flag.Parse()

	log.Println("启动服务 ...")
	var nodeAllocation xid.NodeAllocation
	if *model == "single" {
		nodeAllocation = xid.NewNodeAllocationSingle()
	} else if *model == "redis" {
		nodeAllocation = xid.NewNodeAllocationRedis(*redisAddr, *redisPwd)
	}
	xid.ConfigBits(nodeAllocation, *nodeBits, *stepBits)
	clean := func(ctx context.Context) {
		nodeAllocation.DestroyNode(ctx)
	}
	graceShutdownServe(*webAddr, webapp, clean)

}

func graceShutdownServe(addr string, router http.Handler, clean func(ctx context.Context)) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clean(ctx)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}
