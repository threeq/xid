package xid

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type singleNodeIdAllocation struct {
}

func (alloc *singleNodeIdAllocation) Node(nodeMax int) int {
	return rand.Intn(nodeMax)
}

func (alloc *singleNodeIdAllocation) DestroyNode(timeoutCtx context.Context) {
}

func NewNodeAllocationSingle() NodeAllocation {
	log.Println("初始化 single node 分配器")
	return &singleNodeIdAllocation{}
}
