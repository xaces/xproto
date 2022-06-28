package xproto

import (
	"time"
)

type AdapterHandler func([]byte) IClient

type Handler struct {
	// Access register
	Access func([]byte, *Access) (interface{}, error)
	// Dropped
	Dropped func(interface{}, *Access, error)
	// status callback
	Status func(string, interface{}, *Status)
	// alarm callback
	Alarm func([]byte, interface{}, *Alarm)
	// event callback
	Event func([]byte, interface{}, *Event)
	// bin frame callback
	Frame func(interface{}, *Frame) error
	// File callback
	File func(interface{}, *File)
}

type Options struct {
	// server host
	Host string
	// server port
	Port uint16
	// http request timeout
	RequestTimeout time.Duration
	// 接收超时 s
	RecvTimeout float64
}

// InterProtocol 协议接口
type IClient interface {
	// callback init
	OnInit(ctx *Conn)
	// callback idel; for Proactively respond to tasks
	OnIdle(*time.Time, *time.Time) error
	// callback decodec; return valid data length and error
	OnDecodec([]byte) (int, error)
	// callback handle; if return error not nil, will close connection
	OnHandle() error
	// Request 发送指令
	Request(ReqType, interface{}, interface{}) error
}
