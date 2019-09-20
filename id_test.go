package xid

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {

	var IdNode, _ = NewIDGen(1)

	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want *ID
	}{
		{"ok1", args{IdNode.Next()}, &ID{node:1}},
		{"ok2", args{IdNode.Next()}, &ID{node:1}},
		{"ok3", args{IdNode.Next()}, &ID{node:1}},
		{"ok4", args{IdNode.Next()}, &ID{node:1}},
	}

	step := 0
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.node = 1
			tt.want.step = int64(step)
			tt.want.second = time.Now().UnixNano() / SECOND

			if got := Parse(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Error(got.time())
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
		step++
	}
}
