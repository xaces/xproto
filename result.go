package xproto

import (
	"reflect"
	"sync"
	"time"
)

// rValue 设备响应数据
type rValue struct {
	Data  interface{}
	Done  chan bool
	Error error
}

type Result struct {
	Vals       sync.Map
	reqTimeout time.Duration
}

var DefaultResult = &Result{
	reqTimeout: 30,
}

func NewResult(timeout time.Duration) *Result {
	return &Result{
		reqTimeout: timeout,
	}
}

// Get 获取
func (o *Result) Get(key interface{}) *rValue {
	if res, ok := o.Vals.Load(key); ok {
		return res.(*rValue)
	}
	return nil
}

// Set 设置
func (o *Result) Set(key string, result interface{}) error {
	if r := o.Get(key); r != nil {
		v := reflect.ValueOf(r.Data).Interface().(*interface{})
		*v = result
		r.Done <- true
	}
	return nil
}

// SetError 设置错误
func (o *Result) SetError(key string, err error) error {
	if r := o.Get(key); r != nil {
		r.Error = err
		r.Done <- true
	}
	return nil
}

// WaitVal 等待结果
func (o *Result) WaitVal(k string, v interface{}) error {
	if k == "" {
		return ErrParam
	}
	if v == nil || k == "ok" {
		return nil
	}
	r := &rValue{
		Done: make(chan bool),
		Data: v,
	}
	o.Vals.Store(k, r)
	defer o.Vals.Delete(k)
	select {
	case <-time.After(time.Second * o.reqTimeout):
		return ErrTimeout
	case <-r.Done:
		return r.Error
	}
}
