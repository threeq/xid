package xid

import (
	"sync"
	"time"
)

const MAX = 1000

// ID14Generator
type ID14Generator struct {
	mu    sync.Mutex
	time  int64
	start int64
	step  int64
	node  int64
}

// ---10--- | ---1--- | ---3---  // 260年，单个发号器速度 1000/s【最大时间 9999999999 => 2286-11-21 01:46:39】
// 时间戳(s) | 机器 id  |  计数器
func (g *ID14Generator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now().Unix()
	if now == g.time {
		g.step = (g.step + 1) % MAX
		if g.step == g.start {
			for now <= g.time {
				now = time.Now().Unix()
			}
			g.start = 0
			g.step = g.start
		}
	} else {
		g.start = 0
		g.step = g.start
	}
	g.time = now
	return now*10000 + g.node*1000 + g.step

}

func NewID14Gen(node int) (IDGen, error) {
	return &ID14Generator{node: int64(node)}, nil
}
