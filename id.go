package xid

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

/**
 浏览器 JS 支持最大的数是 2^53 -1，所以这里取低位 52 位
 设计如下：
	时间取秒
		|... 12 ...|... 32 ...|... 5 ...| ... 15 ...|
		|  unused  | time(s)  |  curNodeId   |    step   |

		time： 		时间戳秒数，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2105 年
		step：		机器 id，默认取 5 位。决定了支持 32 个服务器或独立进程
		counter：	计数器，默认取 15 位。决定了每个服务器或进程每秒最多能产生 2^15 个 id
	时间取毫秒
		|... 12 ...|... 41 ...|... 5 ...| ... 6  ...|
		|  unused  | time(ms) |  curNodeId   |    step   |

		time： 		时间戳毫秒数，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2039 年
		step：		机器 id，默认取 5 位。决定了支持 32 个服务器或独立进程
		counter：	计数器，默认取 6 位。决定了每个服务器或进程每毫秒秒最多能产生 2^6 个 id
                    理论每秒最多产生 2^16 *1000 = 6.4w 个 id
 注：machine 和 counter 的位数可一个根据计算机的算力调整

*/

type Units = int64

const (
	SECOND        Units = 1000000000
	MillCROSECOND Units = 1000000
)

var (
	// Epoch is set to the twitter snowflake epoch of 1970-01-01 08:00:00 UTC in seconds
	// You may customize this to set a different epoch for your application.
	defaultEpoch     int64 = 0
	defaultTimeUnit        = MillCROSECOND
	defaultTimeScale       = 1000000000 / MillCROSECOND

	//
	defaultNodeBits uint = 5
	defaultStepBits uint = 6
)

type ID struct {
	second int64
	node   int64
	step   int64
}

// 输出时间，格式：2006-01-02 15:04:05
func (id *ID) time() string {
	tm := time.Unix(id.second, 0)
	return tm.Format("2006-01-02 15:04:05")
}

type IDGenerator struct {
	mu       sync.Mutex
	epoch    time.Time
	nodeBits uint
	stepBits uint

	time int64
	node int64
	step int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint
	nodeShift uint
}

func (n *IDGenerator) Next() int64 {
	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / defaultTimeUnit

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / defaultTimeUnit
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	id := (now << n.timeShift) |
		(n.node << n.nodeShift) |
		(n.step)

	n.mu.Unlock()
	return id
}

func (n *IDGenerator) Parse(id int64) *ID {
	return ParseByBits(id, n.nodeBits, n.stepBits)
}

func NewIDGen(node int64) (*IDGenerator, error) {
	return NewIDGenBits(node, defaultNodeBits, defaultStepBits)
}

func NewIDGenBits(node int64, nodeBits, stepBits uint) (*IDGenerator, error) {

	if node < 0 {
		panic(errors.New("current node id unset, please set current node id"))
	}

	n := IDGenerator{}
	n.node = node
	n.nodeBits = nodeBits
	n.stepBits = stepBits
	n.nodeMax = -1 ^ (-1 << nodeBits)
	n.nodeMask = n.nodeMax << stepBits
	n.stepMask = -1 ^ (-1 << stepBits)
	n.timeShift = nodeBits + stepBits
	n.nodeShift = stepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("IDGenerator number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(defaultEpoch/defaultTimeScale, defaultEpoch*defaultTimeUnit).Sub(curTime))

	return &n, nil
}

func Parse(id int64) *ID {
	return ParseByBits(id, defaultNodeBits, defaultStepBits)
}

func ParseByBits(id int64, nodeBits, stepBits uint) *ID {

	var nodeMask int64 = -1 ^ (-1 << nodeBits)
	var stepMask int64 = -1 ^ (-1 << stepBits)

	return &ID{
		second: (id >> (nodeBits + stepBits)) / defaultTimeScale,
		node:   (id >> stepBits) & nodeMask,
		step:   id & stepMask,
	}
}
