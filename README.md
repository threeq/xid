# xid

分布式 ID 生成工具

# 依赖

* redis


# 运行 web 服务

## 单节点 

```bash
docker run -p8080:8080 threewq/xid:1.0.0 -model=single 
```

所有支持参数

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
