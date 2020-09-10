package xid

type IDGen interface {
	Next() int64
}
