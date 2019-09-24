# xid

分布式 ID 生成工具

# 依赖

* redis


# 运行 web 服务

## 单节点 

```bash
docker run -p8080:8080 threewq/xid:1.0.0
```

获取 id: GET `http://127.0.0.1:8080`

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
docker run --rm threewq/xid:1.0.0 -h
```

# 性能测试

服务器配置：`4C8G Linux`

启动参数
```sh
bin/xid_linux -node-bits=4 -step-bits=7
```

压测配置
```sh
./wrk/wrk -t1000 -c3000 -d2m --latency http://10.105.55.218:8080
```

结果
```text
Running 2m test @ http://10.105.55.218:8080
  1000 threads and 3000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    29.66ms   27.34ms 300.99ms   44.86%
    Req/Sec   112.96     43.77     2.31k    73.34%
  Latency Distribution
     50%   30.01ms
     75%   54.76ms
     90%   65.99ms
     99%   95.61ms
  13530842 requests in 2.00m, 1.68GB read
Requests/sec: 112630.42
Transfer/sec:     14.29MB
```