package xid

import (
	"context"
	"log"
	"sync"
)

type NodeAllocation interface {
	Node(nodeMax int) int
	DestroyNode(timeoutCtx context.Context)
}

var mu sync.Mutex
var curNodeId = -1
var idGenerators = map[string]IDGen{}
var genFactory = NewIDSnakeGen

func Config(mode string, nodeAlloc NodeAllocation) {
	ConfigCustom(mode, nodeAlloc, defaultEpoch, defaultTimeUnit, defaultNodeBits, defaultStepBits)
}

func ConfigBits(mode string, nodeAlloc NodeAllocation, nodeBits, stepBits uint) {
	ConfigCustom(mode, nodeAlloc, defaultEpoch, defaultTimeUnit, nodeBits, stepBits)
}

func ConfigCustom(mode string, nodeAlloc NodeAllocation, epoch int64, timeUnit Units, nodeBits, stepBits uint) {
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

	if mode == "snake" {
		nodeMax := -1 ^ (-1 << nodeBits)
		curNodeId = nodeAlloc.Node(nodeMax)
		genFactory = NewIDSnakeGen
	} else {
		nodeMax := 10
		curNodeId = nodeAlloc.Node(nodeMax)
		genFactory = NewID14Gen
	}

	log.Printf("current node id is: %d", curNodeId)
}

func MultiIdGenerator(gen string) IDGen {
	if idGen, ok := idGenerators[gen]; ok {
		return idGen
	}
	mu.Lock()
	if idGen, ok := idGenerators[gen]; ok {
		mu.Unlock()
		return idGen
	}
	idGen, _ := genFactory(curNodeId)

	idGenerators[gen] = idGen
	mu.Unlock()
	return idGen
}
