package xid

import (
	"context"
	"log"
	"sync"
)

type NodeAllocation interface {
	Node(mode string, nodeMax int) int
	DestroyNode(timeoutCtx context.Context)
}

// MaxBatchNum 最大批量生成数
const MaxBatchNum = 1000

var mu = &sync.Mutex{}
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
		curNodeId = nodeAlloc.Node(mode, nodeMax)
		genFactory = NewIDSnakeGen
	} else {
		nodeMax := 9
		curNodeId = nodeAlloc.Node(mode, nodeMax)
		genFactory = NewID14Gen
	}

	log.Printf("current node id is: %d", curNodeId)
}

func MultiIdGenerator(gen string) IDGen {
	if idGen, ok := idGenerators[gen]; ok {
		return idGen
	}
	mu.Lock()
	defer mu.Unlock()

	if idGen, ok := idGenerators[gen]; ok {
		return idGen
	}
	idGen, _ := genFactory(curNodeId)

	idGenerators[gen] = idGen
	return idGen
}

// GetIDS 获取多个ID
func GetIDS(gen string, num int) []int64 {
	if num > MaxBatchNum {
		num = MaxBatchNum
	}
	if num < 1 {
		num = 1
	}
	idGen := MultiIdGenerator(gen)
	var ids = make([]int64, num)
	for i:=0;i<num; i++ {
		ids[i] = idGen.Next()
	}
	return ids
}
