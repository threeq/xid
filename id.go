package xid

import (
	"context"
	"log"
	"sync"
)

type IdType int

const (
	ID_Snake IdType = iota
	ID_14
)

type IDGen interface {
	Next() int64
}

type NodeAllocator interface {
	Node(nodeMax int) int
	DestroyNode(timeoutCtx context.Context)
}

var mu = &sync.Mutex{}
var curNodeId = -1
var idGenerators = []map[string]IDGen{
	make(map[string]IDGen), make(map[string]IDGen),
}

func multiIdGenerator(idType IdType, gen string) IDGen {
	if idGen, ok := idGenerators[idType][gen]; ok {
		return idGen
	}
	mu.Lock()
	defer mu.Unlock()

	if idGen, ok := idGenerators[idType][gen]; ok {
		return idGen
	}

	var idGen IDGen
	switch idType {
	case ID_Snake:
		idGen, _ = NewIDSnakeGen(
			defaultCfg.currentNodeId,
			defaultCfg.snake.startTime,
			defaultCfg.snake.nodeBits,
			defaultCfg.snake.stepBits,
			defaultCfg.snake.timeUnit,
		)
	case ID_14:
		idGen, _ = NewID14Gen(defaultCfg.currentNodeId)
	default:
		log.Panicf("不支持类型 %v", idType)
	}

	idGenerators[idType][gen] = idGen
	return idGen
}

func GetID(idType IdType, gen string) int64 {
	idGen := multiIdGenerator(idType, gen)
	return idGen.Next()
}

// GetIDS 获取多个ID
func GetIDS(idType IdType, gen string, num int) []int64 {
	if num > defaultCfg.maxBatchNum {
		num = defaultCfg.maxBatchNum
	}
	if num < 1 {
		num = 1
	}
	idGen := multiIdGenerator(idType, gen)
	var ids = make([]int64, num)
	for i := 0; i < num; i++ {
		ids[i] = idGen.Next()
	}
	return ids
}
