# xid 
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/threewq/xid)](https://cloud.docker.com/repository/docker/threewq/xid/builds)
[![Docker Pulls](https://img.shields.io/docker/pulls/threewq/xid)](https://cloud.docker.com/repository/docker/threewq/xid/tags)

分布式 ID 生成工具

# 依赖

* redis

# 运行 web 服务

## 单节点 

```bash
docker run -p8080:8080 threewq/xid:1.1.0
```

获取 id: 
* 方式 1：`GET http://127.0.0.1:8080`
* 方式 2：`GET http://127.0.0.1:8080?gen=type1`
* 方式 2：`GET http://127.0.0.1:8080/gen/type2`

查看所有支持参数

```text
Usage of /app/xid:
  -model string
        运行模式：single 单机 id 生成；redis 使用redis 分布式 id生成 (default "single")
  -node-bits uint
        机器长度 (default 5)
  -redis-addr string
        redis 地址和端口 (default "localhost:6379")
  -redis-pwd string
        redis 密码
  -step-bits uint
         (default 6)
  -web-addr string
        web 监听地址和端口 (default ":8080")

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
./wrk/wrk -t32 -c300 --latency http://10.105.55.218:8080
```

结果
```text
Running 10s test @ http://10.105.55.218:8080
  32 threads and 300 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.76ms    2.07ms  34.96ms   70.47%
    Req/Sec     3.50k   243.75     7.20k    78.05%
  Latency Distribution
     50%    2.52ms
     75%    3.91ms
     90%    5.26ms
     99%    9.22ms
  1123959 requests in 10.10s, 142.56MB read
Requests/sec: 111284.29
Transfer/sec:     14.12MB
```