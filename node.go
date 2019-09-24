package xid

import (
	"context"
	"log"
	"sync"
)

type NodeAllocation interface {
	Node(nodeMax int) int64
	DestroyNode(timeoutCtx context.Context)
}

var mu sync.Mutex
var curNodeId int64 = -1
var idGenerators = map[string]*IDGenerator{}

func Config(nodeAlloc NodeAllocation) {
	ConfigCustom(nodeAlloc, defaultEpoch, defaultTimeUnit, defaultNodeBits, defaultStepBits)
}

func ConfigBits(nodeAlloc NodeAllocation, nodeBits, stepBits uint) {
	ConfigCustom(nodeAlloc, defaultEpoch, defaultTimeUnit, nodeBits, stepBits)
}

func ConfigCustom(nodeAlloc NodeAllocation, epoch int64, timeUnit Units, nodeBits, stepBits uint) {
	if nodeBits+stepBits > 22 {
		log.Fatal("nodeBits + stepBits 超过最大值【22】")
	}

	if epoch < 0 {
		log.Fatal("epoch 不能为负数")
	}

	defaultEpoch = epoch
	defaultTimeUnit = timeUnit
	defaultTimeScale = 1000000000 / timeUnit
	defaultNodeBits = nodeBits
	defaultStepBits = stepBits

	nodeMax := -1 ^ (-1 << nodeBits)
	curNodeId = nodeAlloc.Node(nodeMax)
	log.Printf("current node id is: %d", curNodeId)
}

func MultiIdGenerator(gen string) *IDGenerator {
	if idGen, ok := idGenerators[gen]; ok {
		return idGen
	}
	mu.Lock()
	if idGen, ok := idGenerators[gen]; ok {
		mu.Unlock()
		return idGen
	}
	idGen, _ := NewIDGen(curNodeId)
	idGenerators[gen] = idGen
	mu.Unlock()
	return idGen
}
