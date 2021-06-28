package xid

import (
	"testing"
)

func TestGetIDS(t *testing.T) {
	type args struct {
		gen string
		num int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"-1", args{"", -1}, 1},
		{"0", args{"", 0}, 1},
		{"1", args{"", 1}, 1},
		{"10", args{"", 10}, 10},
		{"1000", args{"", 1000}, 1000},
		{"10000", args{"", 10000}, 1000},
	}

	Config("snake", NewNodeAllocationSingle())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIDS(tt.args.gen, tt.args.num)
			if len(got) != tt.want && got[0] > 0 {
				t.Errorf("GetIDS() = %v, want %v", got, tt.want)
			}

			for i := 0; i < len(got)-1; i++ {
				if got[i] >= got[i+1] {
					t.Errorf("GetIDS() = \n%v\n%v", got[i], got[i+1])
				}
			}
		})
	}
}
