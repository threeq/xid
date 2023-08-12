package xid

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type singleNodeIdAllocation struct {
	currentNodeId int
}

func (alloc *singleNodeIdAllocation) Node(nodeMax int) int {
	if curNodeId >= 0 && curNodeId <= nodeMax {
		return curNodeId
	}

	nodeCount := nodeMax + 1
	alloc.currentNodeId = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(nodeCount)
	log.Println("当前节点 ID：", alloc.currentNodeId, "。最大数量：", nodeMax)
	return alloc.currentNodeId
}

func (alloc *singleNodeIdAllocation) DestroyNode(timeoutCtx context.Context) {
}

func NewNodeAllocationSingle() NodeAllocator {
	log.Println("初始化 single node 分配器")
	return &singleNodeIdAllocation{
		currentNodeId: -1,
	}
}
