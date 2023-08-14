# xid 
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/threewq/xid)](https://cloud.docker.com/repository/docker/threewq/xid/builds)
[![Docker Pulls](https://img.shields.io/docker/pulls/threewq/xid)](https://cloud.docker.com/repository/docker/threewq/xid/tags)

分布式 ID 生成工具

# 依赖

* redis

# 运行 web 服务

## 单节点 

```bash
docker run -p8888:8888 threewq/xid
```

### Http 获取 id: 
* 方式 1：`GET http://127.0.0.1:8080/xid/snake`
* 方式 2：`GET http://127.0.0.1:8080/xid/snake?gen=type1`
* 方式 3：`GET http://127.0.0.1:8080/gen/snake?num=100`
* 方式 3：`GET http://127.0.0.1:8080/gen/snake?gen=type1&num=100`

### Tcp 获取参数


查看 所有支持参数

```text
Usage of ./bin/xid_mac:
  -c string
        The configuration file defaults to <exeDir>/conf/zinx.json if it is not set. (default "/Users/three3q/workspaces/myself/xid/conf/zinx.json")
  -node string
        node 分配方式：single 单机 id 生成；redis 使用redis 分布式 id生成 (default "single")
  -node-bits uint
        机器长度 (default 4)
  -path string
        访问路径 (default "/xid")
  -protocols string
        服务协议【http/tcp】，同时支持多个用+分割 (default "http")
  -redis-addr string
        redis 地址和端口 (default "localhost:6379")
  -redis-pwd string
        redis 密码
  -start string
        开始时间 (default "1970-01-01 08:00:00")
  -step-bits uint
        计数器长度 (default 16)
  -tcp-addr int
        tcp 监听地址和端口 (default 8999)
  -time-unit string
        时间单位: s,ms,10ms,100ms (default "s")
  -types string
        id 生成模式【snake/id14】，同时支持多个用+分割 (default "snake")
  -web-addr string
        web 监听地址和端口 (default ":8888")
```

## 查看所有运行参数

```bash
docker run --rm threewq/xid -h
```

# 性能测试

服务器配置：`4C8G Linux`

启动参数
```sh
bin/xid_linux -node-bits=4 -step-bits=7
```

压测配置
```sh
wrk --latency http://127.0.0.1:8888/xid/snake 
```

结果
```text
Running 10s test @ http://127.0.0.1:8888/xid/snake
  2 threads and 10 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   138.71us  154.01us   8.64ms   96.95%
    Req/Sec    31.16k     4.76k   52.12k    87.62%
  Latency Distribution
     50%  121.00us
     75%  164.00us
     90%  213.00us
     99%  428.00us
  626218 requests in 10.10s, 90.18MB read
Requests/sec:  61999.03
Transfer/sec:      8.93MB
```