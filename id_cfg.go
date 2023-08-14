package xid

import (
	"log"
	"math"
	"strings"
	"time"
)

type idGeneratorsCfg struct {
	types         []IdType
	maxBatchNum   int
	nodeAlloc     NodeAllocator
	currentNodeId int

	snake                *snakeConfig
	redisXidNodesHashKey string
}

var defaultCfg = &idGeneratorsCfg{
	types:         []IdType{ID_Snake},
	currentNodeId: -1,
	maxBatchNum:   1000,
	nodeAlloc:     nil,
	snake: &snakeConfig{
		startTime: 0,
		timeUnit:  Second,

		nodeBits: 4,
		stepBits: 16,
	},
	redisXidNodesHashKey: "xid:nodes",
}

type Option func(*idGeneratorsCfg)

func Options(opts ...Option) {
	for _, opt := range opts {
		opt(defaultCfg)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Init() {
	if defaultCfg.nodeAlloc == nil {
		defaultCfg.nodeAlloc = NewNodeAllocationSingle()
	}

	// 同时支持多种 ID 类型时，需要获取最小的 nodeMax
	nodeMax := math.MaxInt
	for _, t := range defaultCfg.types {
		switch t {
		case ID_Snake:
			nodeMax = min(nodeMax, -1^(-1<<defaultCfg.snake.nodeBits))
		case ID_14:
			nodeMax = min(nodeMax, 9)
		}
	}
	defaultCfg.currentNodeId = defaultCfg.nodeAlloc.Node(nodeMax)
}

func GetIdTypes() []IdType {
	return defaultCfg.types
}

func MaxBatchNum(num int) Option {
	return func(o *idGeneratorsCfg) {
		o.maxBatchNum = num
	}
}

func NodeAlloc(nodeAlloc NodeAllocator) Option {
	return func(o *idGeneratorsCfg) {
		o.nodeAlloc = nodeAlloc
	}
}

func SnakeStartTime(startData string) Option {
	time, err := time.Parse("2006-01-02 15:04:05", startData)
	if err != nil {
		log.Fatal(err)
	}
	return func(o *idGeneratorsCfg) {
		o.snake.startTime = time.Unix()
	}
}

func SnakeTimeUnit(unitDesc string) Option {
	unit := Second
	switch unitDesc {
	case "s":
		unit = Second
	case "ms":
		unit = Millisecond
	case "10ms":
		unit = Millisecond10
	case "100ms":
		unit = Millisecond100
	default:
		log.Fatalf("时间单位错误：%s。只接受：s,ms,10ms,100ms", unitDesc)
	}
	return func(o *idGeneratorsCfg) {
		o.snake.timeUnit = unit
	}
}

func SnakeNodeBits(bits int) Option {
	return func(o *idGeneratorsCfg) {
		o.snake.nodeBits = uint(bits)
	}
}

func SnakeStepBits(bits int) Option {
	return func(o *idGeneratorsCfg) {
		o.snake.stepBits = uint(bits)
	}
}

func RedisXidNodesHashKey(key string) Option {
	return func(o *idGeneratorsCfg) {
		o.redisXidNodesHashKey = key
	}
}

func RunTypes(types string) Option {
	return func(o *idGeneratorsCfg) {
		o.types = []IdType{}
		types := strings.Split(types, "+")
		for _, t := range types {
			switch t {
			case "snake":
				o.types = append(o.types, ID_Snake)
			case "id14":
				o.types = append(o.types, ID_14)
			default:
				log.Fatalf("不支持的 ID 类型：%s", t)
			}
		}
	}
}
