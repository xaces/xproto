package xproto

import (
	"testing"
	"time"
)

// Client t协议结构体
type Client struct {
	*Conn
}

// callback init
func (p *Client) OnInit(ctx *Conn) {
}

// callback idel; for Proactively respond to tasks
func (p *Client) OnIdle(rect, sendt *time.Time) error {
	return nil
}

// callback decodec; return valid data length and error
func (p *Client) OnDecodec(b []byte) (int, error) {
	return 0, nil
}

// callback handle; if return error not nil, will close connection
func (p *Client) OnHandle() error {
	return nil
}

// Request 发送指令
func (p *Client) Request(cmd ReqCode, s interface{}, r interface{}) error {
	return nil
}

func TestServer(t *testing.T) {
	s, err := NewServer(&Options{Port: 22000}, func(b []byte) IClient {
		return &Client{}
	})
	if err != nil {
		return
	}
	defer s.Release()
	s.Handler.Status = LogStatus
	s.ListenTCPAndServe()
}

func TestRequestLivestream(t *testing.T) {
	live := LiveStream{
		Channel:    1,
		StreamType: 0,
		On:         1,
		Frames:     "1;2;3",
		Server:     "192.168.1.203:22000",
	}
	// 不等待返回结果
	SyncSend(Req_LiveStream, live, nil, "deviceNo")
	// 等待结果
	var reslut interface{}
	live.On = 0
	SyncSend(Req_LiveStream, live, &reslut, "deviceNo")
}
