package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/threeq/xid"
	"github.com/valyala/fasthttp"
)

var (
	types        *string
	protocols    *string
	basePath     *string
	webAddr      *string
	tcpAddr      *int
	node         *string
	redisAddr    *string
	redisPwd     *string
	start        *string
	timeUnitDesc *string
	nodeBits     *uint
	stepBits     *uint
)

func main() {
	types = flag.String("types", "snake", "id 生成模式【snake/id14】，同时支持多个用+分割")
	protocols = flag.String("protocols", "http", "服务协议【http/tcp】，同时支持多个用+分割")
	basePath = flag.String("path", "/xid", "访问路径")
	webAddr = flag.String("web-addr", ":8888", "web 监听地址和端口")
	tcpAddr = flag.Int("tcp-addr", 8999, "tcp 监听地址和端口")
	node = flag.String("node", "single", "node 分配方式：single 单机 id 生成；redis 使用redis 分布式 id生成")

	// redis 特定参数
	redisAddr = flag.String("redis-addr", "localhost:6379", "redis 地址和端口")
	redisPwd = flag.String("redis-pwd", "", "redis 密码")

	// snake 特定参数
	start = flag.String("start", "1970-01-01 08:00:00", "开始时间")
	timeUnitDesc = flag.String("time-unit", "s", "时间单位: s,ms,10ms,100ms")
	nodeBits = flag.Uint("node-bits", 4, "机器长度")
	stepBits = flag.Uint("step-bits", 16, "计数器长度")
	flag.Parse()

	log.Println("支持的 ID 类型：", *types)
	log.Println("支持的服务类型：", *protocols)
	var nodeAllocation xid.NodeAllocator
	if *node == "single" {
		nodeAllocation = xid.NewNodeAllocationSingle()
	} else if *node == "redis" {
		nodeAllocation = xid.NewNodeAllocationRedis(*redisAddr, *redisPwd)
	}

	xid.Options(
		xid.RunTypes(*types),
		xid.NodeAlloc(nodeAllocation),
		xid.SnakeStartTime(*start),
		xid.SnakeTimeUnit(*timeUnitDesc),
		xid.SnakeNodeBits(int(*nodeBits)),
		xid.SnakeStepBits(int(*stepBits)))
	xid.Init()

	serves, cleans := createProtocolServers(protocols, basePath, webAddr, *tcpAddr, nodeAllocation)
	log.Println("启动服务 ...")

	graceShutdownServes(serves, cleans)
}

// createProtocolServers 实现更多的协议服务类型
func createProtocolServers(protocols *string, basePath *string, webAddr *string, tcpPort int, nodeAllocation xid.NodeAllocator) ([]runServe, []func(ctx context.Context)) {
	var serves []runServe
	var cleans []func(ctx context.Context)
	if strings.Contains(*protocols, "http") {
		srvFastHttp := &fasthttp.Server{
			Concurrency:  1000,
			Handler:      newIDFastHttp(*basePath),
			LogAllErrors: true,
		}

		serves = append(serves, func() error {
			log.Println("开始 http 协议服务", *webAddr)
			return srvFastHttp.ListenAndServe(*webAddr)
		})
		cleans = append(cleans, func(ctx context.Context) {
			srvFastHttp.ShutdownWithContext(ctx)
		})
	}

	if strings.Contains(*protocols, "tcp") {
		srvTcp := newIdTcp(tcpPort)
		serves = append(serves, func() error {
			log.Println("开始 tcp 协议服务:", tcpPort)
			srvTcp.Serve()
			return nil
		})
		cleans = append(cleans, func(ctx context.Context) {
			srvTcp.Stop()
		})
	}

	cleans = append(cleans, func(ctx context.Context) {
		nodeAllocation.DestroyNode(ctx)
	})

	return serves, cleans
}

type runServe = func() error

func graceShutdownServes(serves []runServe, cleans []func(ctx context.Context)) {
	waiting := &sync.WaitGroup{}
	for _, serve := range serves {
		waiting.Add(1)
		go func(serve runServe) {
			defer waiting.Done()
			if err := serve(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}(serve)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, clean := range cleans {
		clean(ctx)
	}

	waiting.Wait()
	log.Println("Server exiting")
}
