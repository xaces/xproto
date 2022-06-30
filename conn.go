package xproto

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

// Conn tcp Conn
type Conn struct {
	conn      net.Conn
	reader    io.Reader
	writer    io.Writer
	client    IClient
	bytesBuf  bytes.Buffer
	recvBytes []byte
	done      int32
	recvTime  time.Time
	sendTime  time.Time
	notify    *Handler
	usrArg    interface{}

	Srv    *Server
	Result *Result
	Access
}

// WriteBytes 直接写入
func (c *Conn) WriteBytes(data []byte) error {
	L, err := c.writer.Write(data)
	if err == nil {
		c.DownTraffic += int64(L)
		c.sendTime = time.Now()
	}
	return err
}

// Close 关闭socket
func (c *Conn) Close() {
	atomic.AddInt32(&c.done, 1)
}

// Access register
func (c *Conn) NotifyAccess(b []byte) (err error) {
	if c.notify.Access != nil {
		c.Access.Online = true
		c.usrArg, err = c.notify.Access(b, &c.Access)
		return
	}
	return nil
}

// Dropped
func (c *Conn) NotifyDropped(err error) {
	c.Access.Online = false
	if c.notify.Dropped != nil {
		c.notify.Dropped(c.usrArg, &c.Access, err)
	}
	if c.LinkType == Link_Signal {
		SyncRemove(c.DeviceNo)
		return
	}
	SyncRemove(c.DeviceNo, c.Session)
}

// alarm callback
func (c *Conn) NotifyStatus(tag string, v *Status) {
	if c.notify.Status != nil {
		v.DeviceID = c.DeviceID
		c.notify.Status(tag, c.usrArg, v)
	}
}

// alarm callback
func (c *Conn) NotifyAlarm(b []byte, v *Alarm) {
	if c.notify.Alarm != nil {
		v.DeviceID = c.DeviceID
		alarmUUID(v)
		c.notify.Alarm(b, c.usrArg, v)
	}
}

// event callback
func (c *Conn) NotifyEvent(b []byte, v *Event) {
	if c.notify.Event != nil {
		v.DeviceID = c.DeviceID
		c.notify.Event(b, c.usrArg, v)
	}
}

// av frame callback
func (c *Conn) NotifyFrame(f *Frame) error {
	if c.notify.Frame == nil {
		return ErrObjectNoExist
	}
	return c.notify.Frame(c.usrArg, f)
}

// file query callback
func (c *Conn) NotityFile(file *File) {
	if c.notify.File != nil {
		c.notify.File(c.usrArg, file)
	}
}

// 接收并返回消息
func (c *Conn) cacheRead(t time.Time) int {
	c.conn.SetReadDeadline(t.Add(time.Millisecond * 200)) // 读超时
	recvLen, err := c.reader.Read(c.recvBytes)
	if err == io.EOF {
		panic(err)
	}
	if recvLen > 0 {
		c.UpTraffic += int64(recvLen)
		c.bytesBuf.Write(c.recvBytes[:recvLen])
	}
	return recvLen
}

func (c *Conn) processMessage() {
	if nil == c.client {
		for _, v := range c.Srv.adapters {
			if c.client = v(c.bytesBuf.Bytes()); c.client != nil {
				break
			}
		}
		if nil == c.client {
			panic(errors.New("can't adapta protocol"))
		}
		c.client.OnInit(c)
	}
	for {
		l, err := c.client.OnDecodec(c.bytesBuf.Bytes())
		if err != nil {
			panic(err)
		}
		if l == 0 {
			break
		}
		err = c.client.OnHandle()
		if err != nil {
			panic(err)
		}
		c.bytesBuf.Next(l)
	}
}

func (c *Conn) start() {
	defer func() {
		c.conn.Close()
		e := recover()
		if nil != c.client && c.LinkType != Link_Unknow {
			c.NotifyDropped(fmt.Errorf("%s", e))
		}
		log.Printf("%s closed. <%s> %s | %s\n", c.RemoteAddress, e, c.DeviceNo, c.Session)
	}()
	c.RemoteAddress = c.conn.RemoteAddr().String()
	log.Println(c.RemoteAddress + " connected")
	c.recvTime = time.Now()
	c.sendTime = c.recvTime
	lstDispatchTime := c.recvTime
	for {
		if atomic.LoadInt32(&c.done) > 0 {
			panic(errors.New("closed by user"))
		}
		now := time.Now()
		if c.cacheRead(now) > 0 {
			c.processMessage()
			c.recvTime = now
		}
		sec := now.Sub(lstDispatchTime).Seconds()
		if sec < 5 {
			continue
		}
		lstDispatchTime = now
		sec = now.Sub(c.recvTime).Seconds()
		if sec > c.Srv.RecvTimeout {
			panic("io/timeout by options")
		}
		if c.client == nil {
			continue
		}
		if err := c.client.OnIdle(&c.recvTime, &c.sendTime); err != nil {
			panic(err)
		}
	}
}

// newConn 创建Conn对象
func newConn(conn net.Conn, s *Server) *Conn {
	ctx := &Conn{
		conn:      conn,
		reader:    io.Reader(conn),
		writer:    io.Writer(conn),
		recvBytes: make([]byte, 2048),
		Srv:       s,
		notify:    &s.Handler,
		Result:    s.Result,
	}
	return ctx
}
