package xid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nodeKey1(t *testing.T) {
	s := nodeKey(1)
	println(s)
	assert.Equal(t, "node-1", s)
}
