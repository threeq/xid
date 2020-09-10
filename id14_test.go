package xid

import (
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"strconv"
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
	for i := 0; i < b.N; i++ {
		fmt.Println(g.Next())
	}
}

func TestAlgEqual(t *testing.T) {
	num := fmt.Sprintf("%010d%d%03d", 1599749011, 1, 49)

	n0, _ := strconv.ParseInt(num, 10, 64)
	n1 := int64(1599749011*10000 + 1*1000 + 49)
	fmt.Println(n0)
	fmt.Println(n1)
	assert.Equal(t, n0, n1)
}

func BenchmarkID14Generator_AlgStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		num := fmt.Sprintf("%010d%d%03d", 1599749011, 1, 49)

		_, _ = strconv.ParseInt(num, 10, 64)
	}
}

func BenchmarkID14Generator_AlgM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = int64(1599749011)*10000 + int64(1)*1000 + int64(49)
	}
}
