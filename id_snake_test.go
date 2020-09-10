package xid

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

	var IdNode, _ = NewIDSnakeGen(1)

	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want *SnakeID
	}{
		{"ok1", args{IdNode.Next()}, &SnakeID{node: 1}},
		{"ok2", args{IdNode.Next()}, &SnakeID{node: 1}},
		{"ok3", args{IdNode.Next()}, &SnakeID{node: 1}},
		{"ok4", args{IdNode.Next()}, &SnakeID{node: 1}},
	}

	step := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.node = 1
			tt.want.step = int64(step)
			tt.want.second = time.Now().UnixNano() / Second

			if got := Parse(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Error(got.time(0))
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
		step++
	}
}

func TestIDSnakeGenerator_Next(t *testing.T) {
	defaultEpoch = int64(time.Now().Nanosecond() / 1000000)
	idGen, _ := NewIDSnakeGen(0)

	id := idGen.Next()
	log.Println(id)
	log.Printf("%b\n", id)
	log.Println("1011011010110001100001011111000000101101100000000000")
	log.Println(len("1011011010110001100001011111000000101101100000000000"))

	log.Println(Parse(id).time(0))
}

func BenchmarkIDSnakeGenerator_Next(b *testing.B) {
	idGen, _ := NewIDSnakeGen(0)

	for i := 0; i < b.N; i++ {
		id := idGen.Next()
		if id<3 {
			b.Fatalf("error")
		}
	}
}

func TestParse2(t *testing.T) {
	idNum := Parse(100632443644096)
	println(fmt.Sprintf("%+v", idNum))
	println(idNum.time(1546272000))
}