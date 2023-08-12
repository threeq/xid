package xid

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

/**
 浏览器 JS 支持最大的数是 2^53 -1，所以这里取低位 52 位
 设计如下：
 	时间力度-毫秒
		|... 12 ...|... 41 ...|... 5 ...| ... 6  ...|
		|  unused  | time(ms) |  node   |    step   |

		time： 		时间戳毫秒数，可用约69年，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2039 年
		step：		机器 id，默认取 5 位。决定了支持 32 个服务器或独立进程
		counter：	计数器，默认取 6 位。决定了每个服务器或进程每毫秒最多能产生 2^6 个 id
                    单台服务器理论每秒最多产生 2^6 *1000 = 6.4w 个 id
	时间粒度-秒
		|... 12 ...|... 32 ...|... 4 ...| ... 16 ...|
		|  unused  | time(s)  |  node   |    step   |

		time： 		时间戳秒数，可用约135年，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2105 年
		step：		机器 id，默认取 4 位。决定了支持 16 个服务器或独立进程
		counter：	计数器，默认取 16 位。决定了每个服务器或进程每秒最多能产生 2^16 个 id
		    		单台服务器理论每秒最多产生 2^16 = 65535 个 id
	时间力度-10毫秒
		|... 12 ...|... 39 ...|... 4 ...| ... 9  ...|
		|  unused  | time(ms) |  node   |    step   |

		time： 		时间戳毫秒数，可用约173年，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2143 年
		step：		机器 id，默认取 4 位。决定了支持 16 个服务器或独立进程
		counter：	计数器，默认取 9 位。决定了每个服务器或进程每10毫秒最多能产生 2^9 个 id
                    单台服务器理论每秒最多产生 2^9 *100 = 5.12w 个 id
	时间力度-100毫秒
		|... 12 ...|... 36 ...|... 4 ...| ... 12  ...|
		|  unused  | time(ms) |  node   |    step   |

		time： 		时间戳毫秒数，可用约217年，默认以 1970-01-01 08:00:00 开始计算，可以支持到 2187 年
		step：		机器 id，默认取 4 位。决定了支持 16 个服务器或独立进程
		counter：	计数器，默认取 12 位。决定了每个服务器或进程每100毫秒最多能产生 2^12 个 id
                    单台服务器理论每秒最多产生 2^12*100 = 4w 个 id
 注：machine 和 counter 的位数可一个根据计算机的算力调整

*/

type Units = int64

const (
	Millisecond    Units = 1000000
	Millisecond10  Units = 10000000
	Millisecond100 Units = 100000000
	Second         Units = 1000000000

	/**
	 * 最大容忍时间, 单位毫秒, 即如果时钟只是回拨了该变量指定的时间, 那么等待相应的时间即可;
	 * 考虑到服务的高性能, 这个值不易过大
	 */
	MaxBackwardMs = 30 * Millisecond
)

type snakeConfig struct {
	// Epoch is set to the twitter snowflake startTime of 1970-01-01 08:00:00 UTC in seconds
	// You may customize this to set a different startTime for your application.
	startTime int64
	timeUnit  Units

	nodeBits uint
	stepBits uint
}

type SnakeID struct {
	second int64
	node   int64
	step   int64
}

// 输出时间，格式：2006-01-02 15:04:05
func (id *SnakeID) time(epoch int64) string {
	tm := time.Unix(id.second+epoch, 0)
	return tm.Format("2006-01-02 15:04:05")
}

type IDSnakeGenerator struct {
	mu       sync.Mutex
	epoch    time.Time
	nodeBits uint
	stepBits uint
	timeUnit Units

	time int64
	node int64
	step int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint
	nodeShift uint
}

func (n *IDSnakeGenerator) Next() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Nanoseconds() / n.timeUnit

	//时间回拨处理
	if now < n.time {
		// 如果时钟回拨在可接受范围内, 等待即可
		backwards := (n.time - now) * n.timeUnit
		if backwards < MaxBackwardMs {
			time.Sleep(time.Duration(backwards))
		} else {
			msg := fmt.Sprintf("clock is moving backwards. Rejecting requests until %d.", n.time)
			log.Println(msg)
			panic(msg)
		}
	}

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / n.timeUnit
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	id := (now << n.timeShift) |
		(n.node << n.nodeShift) |
		(n.step)

	return id
}

func (n *IDSnakeGenerator) Parse(id int64) *SnakeID {
	var nodeMask int64 = -1 ^ (-1 << n.nodeBits)
	var stepMask int64 = -1 ^ (-1 << n.stepBits)

	timeScale := Second / n.timeUnit

	return &SnakeID{
		second: (id >> (n.nodeBits + n.stepBits)) / timeScale,
		node:   (id >> n.stepBits) & nodeMask,
		step:   id & stepMask,
	}

}

func NewIDSnakeGen(node int, startTime int64, nodeBits, stepBits uint, timeUnit Units) (*IDSnakeGenerator, error) {

	n := IDSnakeGenerator{
		mu: sync.Mutex{},
	}
	n.time = 0
	n.node = int64(node)
	n.nodeBits = nodeBits
	n.stepBits = stepBits
	n.nodeMax = -1 ^ (-1 << nodeBits)
	n.nodeMask = n.nodeMax << stepBits
	n.stepMask = -1 ^ (-1 << stepBits)
	n.timeShift = nodeBits + stepBits
	n.nodeShift = stepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("IDSnakeGenerator number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(startTime, 0).Sub(curTime))
	n.timeUnit = timeUnit

	return &n, nil
}
