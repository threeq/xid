package xid

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"time"
)

type UnlockFunc func()

const LockKey = "Rds_Xid_Hash_Key_Lock"
const LockKeyTimeoutMs = 3 * 1000 // 3 秒
const LockKeyRetryTimes = 100     // 100 次
const RdsXidNodesHash = "Rds_Xid_Node_Key_"
const NodeIdRefreshTImeSecond = 6 // 3 秒

type redisNodeIdAllocation struct {
	rds             *redis.Client
	rdsXidNodesHash string
	nodeId          int
	//salt            string
	shutdown chan interface{}
	ctx      context.Context
	canceler context.CancelFunc
}

func (alloc *redisNodeIdAllocation) Node(mode string, nodeMax int) int {
	if alloc.nodeId >= 0 && alloc.nodeId <= nodeMax {
		return alloc.nodeId
	}

	//alloc.salt = uuid.NewV4().String()
	alloc.rdsXidNodesHash = RdsXidNodesHash + mode
	nodeCount := nodeMax + 1
	// 为了尽量减少冲突，先随机获取一次 nodeId，不过失败在再逐个尝试
	rand.Seed(time.Now().UnixNano())
	startNodeId := rand.Intn(nodeCount)
	r := alloc.rds.HSetNX(alloc.rdsXidNodesHash, nodeKey(startNodeId), time.Now().Unix())
	if r.Val() {
		alloc.nodeId = startNodeId
		alloc.capture()
		return alloc.nodeId
	}

	log.Println("逐个尝试获取未使用 nodeId")
	tryNodeId := (startNodeId + 1) % nodeCount
	for ; tryNodeId != startNodeId; {

		// 多个节点并发控制
		// 这里采用 redis setnx 方式实现简单分布式锁，
		// 这里是在服务器启动时才会获取一次，并发量不会很大
		unlock, err := lock(alloc.rds, LockKey, LockKeyTimeoutMs, LockKeyRetryTimes)
		if err != nil {
			panic(err)
		}
		exist, _ := alloc.rds.HGet(alloc.rdsXidNodesHash, nodeKey(tryNodeId)).Int64()
		now := time.Now().Unix()
		log.Println("check id ", now-exist, tryNodeId)
		if (now - exist) > NodeIdRefreshTImeSecond*10 {
			alloc.rds.HDel(alloc.rdsXidNodesHash, nodeKey(tryNodeId))
			r := alloc.rds.HSetNX(alloc.rdsXidNodesHash, nodeKey(tryNodeId), time.Now().Unix())
			if r.Val() {
				unlock()
				alloc.nodeId = tryNodeId
				alloc.capture()
				return alloc.nodeId
			}
		}
		unlock()

		tryNodeId = (tryNodeId + 1) % nodeCount
	}
	panic(errors.New(fmt.Sprintf("all %d nodes are in use", nodeCount)))
}

func (alloc *redisNodeIdAllocation) DestroyNode(timeoutCtx context.Context) {
	alloc.canceler()
	select {
	case <-timeoutCtx.Done():
	case <-alloc.shutdown:
	}
}

func (alloc *redisNodeIdAllocation) capture() {
	go refreshNodeStatus(alloc.ctx, alloc, alloc.nodeId, alloc.shutdown)
}

func refreshNodeStatus(ctx context.Context, alloc *redisNodeIdAllocation, nodeId int, done chan interface{}) {
	client := alloc.rds
	d := time.Second * NodeIdRefreshTImeSecond
	t := time.NewTicker(d)

	for {
		select {
		case _, _ = <-t.C:
			now := time.Now().Unix()
			log.Printf("refresh id [%d] activity time to %d", nodeId, now)
			client.HSet(alloc.rdsXidNodesHash, nodeKey(nodeId), now)
		case <-ctx.Done():
			t.Stop()
			client.HDel(alloc.rdsXidNodesHash, nodeKey(nodeId))
			_ = client.Close()
			log.Println(fmt.Sprintf("delete node id %d", nodeId))
			done <- 1
			return
		}
	}
}

func NewNodeAllocationRedis(redisAddr, redisPwd string) *redisNodeIdAllocation {
	log.Println("初始化 redis node 分配器")
	rds := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPwd, // no password set
		DB:       15,       // use default DB
		PoolSize: 2,
	})

	_, err := rds.Ping().Result()
	if err != nil {
		panic(err)
	}

	ctx, canceler := context.WithCancel(context.Background())
	done := make(chan interface{})

	return &redisNodeIdAllocation{
		rds:      rds,
		nodeId:   -1,
		ctx:      ctx,
		canceler: canceler,
		shutdown: done,
	}
}

func nodeKey(id int) string {
	return "node-" + string(id)
}
