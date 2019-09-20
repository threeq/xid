package xid

import (
	"context"
	"log"
)


type singleNodeIdAllocation struct {

}

func (alloc *singleNodeIdAllocation) Node(nodeMax int) int64 {
	return 0
}

func (alloc *singleNodeIdAllocation) DestroyNode(timeoutCtx context.Context)  {
}


func NewNodeAllocationSingle() NodeAllocation {
	log.Println("初始化 single node 分配器")
	return &singleNodeIdAllocation{}
}

