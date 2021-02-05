package xid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_nodeKey(t *testing.T) {
	s := nodeKey(1)
	println(s)
	assert.Equal(t, "node-1", s)
}
