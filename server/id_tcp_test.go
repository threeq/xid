package main

import (
	"encoding/binary"
	"fmt"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/stretchr/testify/assert"
	"github.com/threeq/xid"
	"sync"
	"testing"
	"time"
)

func TestGetID(t *testing.T) {
	xid.Options(xid.RunTypes("snake+id14"))
	xid.Init()
	s := newIdTcp(8999)
	w := sync.WaitGroup{}
	go func() {
		s.Serve()
	}()
	client := znet.NewClient("127.0.0.1", 8999)

	client.SetOnConnStart(func(conn ziface.IConnection) {
		client.Conn().SendMsg(1, P2Bytes(&Params{
			Batch: 0,
			Gen:   "123",
		}))
		m := readMsg(conn)
		fmt.Printf("获取ID：%v\n", int64(binary.BigEndian.Uint64(m.GetData())))

		client.Conn().SendMsg(1, P2Bytes(&Params{
			Batch: 3,
			Gen:   "123",
		}))
		m = readMsg(conn)
		fmt.Printf("获取IDS：%v\n", string(m.GetData()))
		w.Done()
	})
	time.Sleep(time.Second)
	w.Add(1)
	client.Start()
	w.Wait()
	s.Stop()
}

func TestP2Bytes(t *testing.T) {
	p := &Params{
		0, "aaa",
	}
	bits := P2Bytes(p)
	fmt.Printf("%v\n", bits)
	bits2 := P2Bytes(&Params{
		12, "bbb",
	})
	fmt.Printf("%v\n", bits2)

	pp := Bytes2P(bits)
	assert.Equal(t, p.Batch, pp.Batch)
	assert.Equal(t, p.Gen, pp.Gen)
}
