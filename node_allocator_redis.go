package xid

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
)

type UnlockFunc func()

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

func (alloc *redisNodeIdAllocation) Node(nodeMax int) int {
	defer func() {
		if alloc.nodeId != -1 {
			log.Println("当前节点 ID：", alloc.nodeId)
		}
	}()

	if alloc.nodeId >= 0 && alloc.nodeId <= nodeMax {
		return alloc.nodeId
	}

	nodeCount := nodeMax + 1
	// 为了尽量减少冲突，先随机获取一次 nodeId，不过失败在再逐个尝试
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	startNodeId := random.Intn(nodeCount)
	now := time.Now().UnixNano()
	r := alloc.rds.HSetNX(alloc.rdsXidNodesHash, nodeKey(startNodeId), now)
	exist, _ := alloc.rds.HGet(alloc.rdsXidNodesHash, nodeKey(startNodeId)).Int64()
	if r.Val() {
		alloc.nodeId = startNodeId
		alloc.capture()
		return alloc.nodeId
	} else if (now-exist)/1e9 > NodeIdRefreshTImeSecond*10 {
		// 有可能是上次的节点没有释放，这里再次尝试获取
		alloc.rds.HDel(alloc.rdsXidNodesHash, nodeKey(startNodeId))
		r = alloc.rds.HSetNX(alloc.rdsXidNodesHash, nodeKey(startNodeId), time.Now().UnixNano())
		if r.Val() {
			alloc.nodeId = startNodeId
			alloc.capture()
			return alloc.nodeId
		}
	}

	log.Println("逐个尝试获取未使用 nodeId")
	tryNodeId := (startNodeId + 1) % nodeCount
	for tryNodeId != startNodeId {
		exist, _ := alloc.rds.HGet(alloc.rdsXidNodesHash, nodeKey(tryNodeId)).Int64()
		now := time.Now().Unix()
		log.Println("check id ", now-exist, tryNodeId)
		if (now-exist)/1e9 > NodeIdRefreshTImeSecond*10 {
			alloc.rds.HDel(alloc.rdsXidNodesHash, nodeKey(tryNodeId))
			r := alloc.rds.HSetNX(alloc.rdsXidNodesHash, nodeKey(tryNodeId), time.Now().UnixNano())
			if r.Val() {
				alloc.nodeId = tryNodeId
				alloc.capture()
				return alloc.nodeId
			}
		}

		tryNodeId = (tryNodeId + 1) % nodeCount
	}
	panic("all nodes are in use")
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
		case <-t.C:
			now := time.Now().UnixNano()
			log.Printf("refresh id [%d] activity time to %d", nodeId, now)
			client.HSet(alloc.rdsXidNodesHash, nodeKey(nodeId), now)
		case <-ctx.Done():
			t.Stop()
			client.HDel(alloc.rdsXidNodesHash, nodeKey(nodeId))
			_ = client.Close()
			log.Printf("delete node id %d", nodeId)
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
		rds:             rds,
		rdsXidNodesHash: defaultCfg.redisXidNodesHashKey,
		nodeId:          -1,
		ctx:             ctx,
		canceler:        canceler,
		shutdown:        done,
	}
}

func nodeKey(id int) string {
	return fmt.Sprintf("node-%d", id)
}
