package xproto

import "errors"

var (
	ErrNo            = errors.New("success")               //
	ErrTimeout       = errors.New("timeout")               // 超时
	ErrParam         = errors.New("parameter error")       // 参数错误
	ErrNoResult      = errors.New("no result")             // 结果错误
	ErrObjectExist   = errors.New("object already exist")  // 链接已存在
	ErrObjectNoExist = errors.New("object does not exist") // 链接不存在
	ErrInvalidDevice = errors.New("invalid device")
	ErrUnSupport     = errors.New("unsupport")
	ErrRecvOver      = errors.New("data over")
)
