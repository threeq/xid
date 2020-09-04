package xid

import (
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"testing"
	"time"
)

func TestID14Generator_Next(t *testing.T) {
	g := ID14Generator{}
	n1 := g.Next()
	n2 := g.Next()
	fmt.Println(n1, fmt.Sprintf("len=%d", len(fmt.Sprintf("%d", n1))))
	fmt.Println(n2, fmt.Sprintf("len=%d", len(fmt.Sprintf("%d", n2))))

	assert.Equal(t, n2-n1, int64(1))
	s1 := g.start
	time.Sleep(time.Second)
	g.Next()
	s2 := g.start
	assert.NotEqual(t, s1, s2)
}

func BenchmarkID14Generator_Next(b *testing.B) {
	g := ID14Generator{}
	for i := 0; i <b.N; i++ {
		fmt.Println(g.Next())
	}
}