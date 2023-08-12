package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/aceld/zinx/zpack"
	"github.com/threeq/xid"
	"io"
)

var dp = zpack.NewDataPack()

func newIdTcp(port int) ziface.IServer {
	//1 创建一个server服务
	zconf.UserConfToGlobal(&zconf.Config{
		TCPPort: port,
	})
	s := znet.NewServer()

	//2 配置路由
	for _, t := range xid.GetIdTypes() {
		switch t {
		case xid.ID_Snake:
			s.AddRouter(uint32(xid.ID_Snake), &IdGenHandler{idtype: xid.ID_Snake})
		case xid.ID_14:
			s.AddRouter(uint32(xid.ID_14), &IdGenHandler{idtype: xid.ID_14})
		}
	}

	return s
}

// IdGenHandler MsgId=1的路由
type IdGenHandler struct {
	znet.BaseRouter
	idtype xid.IdType
}

// Handle 消息处理方法
func (h *IdGenHandler) Handle(request ziface.IRequest) {
	//读取客户端的数据
	msg := request.GetMessage()
	//fmt.Printf("%+v\n", msg)
	//fmt.Printf("%v\n", request.GetData())
	//fmt.Println(h.idtype, "recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	p := Bytes2P(msg.GetData())

	if p.Batch > 0 {
		ids := xid.GetIDS(h.idtype, p.Gen, int(p.Batch))
		s, _ := json.Marshal(ids)
		_ = request.GetConnection().SendMsg(uint32(h.idtype), s)
	} else {
		id := xid.GetID(h.idtype, p.Gen)
		_ = request.GetConnection().SendMsg(uint32(h.idtype), msgBytes(id))
	}

}

type Params struct {
	Batch uint32
	Gen   string
}

func P2Bytes(p *Params) []byte {
	return msgBytes(p.Batch, p.Gen)
}

func Bytes2P(data []byte) *Params {
	p := new(Params)
	p.Batch = binary.BigEndian.Uint32(data[0:4])
	p.Gen = string(data[4:])
	return p
}

func msgBytes(data ...interface{}) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	for _, d := range data {
		if d == nil {
			continue
		}

		switch d.(type) {
		case string:
			bytesBuffer.WriteString(d.(string))
		case *string:
			bytesBuffer.WriteString(*(d.(*string)))
		default:
			_ = binary.Write(bytesBuffer, binary.BigEndian, d)
		}

	}
	return bytesBuffer.Bytes()
}

func readMsg(conn ziface.IConnection) ziface.IMessage {
	//服务器就应该给我们回复一个message数据
	binaryDataHead := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(conn.GetConnection(), binaryDataHead); err != nil {
		fmt.Println("read msg head error", err)
		return nil
	}

	fmt.Println(binaryDataHead)

	//先读取数据的head部分，得到id和datalen
	msg, err := dp.Unpack(binaryDataHead)
	if err != nil {
		fmt.Println("unpack err", err)
		return nil
	}

	//然后根据datalen进行第二次读取，得到data
	data := make([]byte, msg.GetDataLen())
	if msg.GetDataLen() > 0 {
		if _, err := io.ReadFull(conn.GetConnection(), data); err != nil {
			fmt.Println("read msg data error", err)
			return nil
		}
	}
	msg.SetData(data)
	return msg
}
