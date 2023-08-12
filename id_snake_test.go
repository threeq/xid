package xid

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

	var idGen, _ = NewIDSnakeGen(1, 0, 4, 16, Second)

	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want *SnakeID
	}{
		{"ok1", args{idGen.Next()}, &SnakeID{node: 1}},
		{"ok2", args{idGen.Next()}, &SnakeID{node: 1}},
		{"ok3", args{idGen.Next()}, &SnakeID{node: 1}},
		{"ok4", args{idGen.Next()}, &SnakeID{node: 1}},
	}

	step := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.node = 1
			tt.want.step = int64(step)
			tt.want.second = time.Now().UnixNano() / Second

			if got := idGen.Parse(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Error(got.time(0))
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
		step++
	}
}

func TestIDSnakeGenerator_Next(t *testing.T) {
	startTime := int64(time.Now().Nanosecond() / int(Millisecond))
	idGen, _ := NewIDSnakeGen(1, startTime, 5, 6, Millisecond)

	id := idGen.Next()
	log.Println(id)
	log.Printf("%b\n", id)
	log.Println("1011011010110001100001011111000000101101100000000000")
	log.Println(len("1011011010110001100001011111000000101101100000000000"))

	log.Println(idGen.Parse(id).time(0))
}

func BenchmarkIDSnakeGenerator_Next(b *testing.B) {
	idGen, _ := NewIDSnakeGen(0, 0, 5, 6, Millisecond)

	for i := 0; i < b.N; i++ {
		id := idGen.Next()
		if id < 3 {
			b.Fatalf("error")
		}
	}
}
