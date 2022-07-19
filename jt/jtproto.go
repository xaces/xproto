package jt

// jt808 2013
import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/xaces/xproto"
	"github.com/xaces/xproto/jt/util"
)

func init() {
	xproto.Use(NewClient)
}

type terminalKey struct {
	Property  uint16
	Upgrade   uint16
	Localtion uint16
	OnDemand  uint16
}

// Client jt8080
type Client struct {
	*xproto.Conn

	jtMsgPacket
	jtBytes  *jtBytes
	terminal *jtRegister
	sequence uint16
	key      terminalKey
}

// isValid 协议判断
func isValid(data []byte) (int, error) {
	if len(data) < jtprotocolHeaderLength {
		return -1, nil
	}
	if data[0] != jtprotocolFlag {
		return -1, errors.New("invalid jt protocol")
	}
	return bytes.IndexByte(data[1:], jtprotocolFlag), nil
}

func NewClient(b []byte) xproto.IClient {
	if l, _ := isValid(b); l < 0 {
		return nil
	}
	return &Client{}
}

func (p *Client) sequenceKey(sequence uint16) string {
	if sequence == 0 {
		return "ok"
	}
	return fmt.Sprintf("%s-%d", p.DeviceNo, sequence)
}

// writeMessage 发送消息
func (p *Client) writeMessage(msgId uint16, b []byte) error {
	rdata := p.jtPacketBytes(msgId, b)
	return p.Conn.WriteBytes(rdata)
}

// writeResp 通用应答
// 0 success; 1 failed; 2 message error;3 unsupport;4 alarm
func (p *Client) writeResp(result uint8) error {
	var data jtBytes
	data.addWord(p.MsgSerialNum)
	data.addWord(p.MsgId)
	data.addByte(result)
	return p.writeMessage(0x8001, data.Bytes)
}

// phoneNumberString 设备号
func (p *Client) phoneNumberString() string {
	phoneNumber := util.Bcd2Number(p.PhoneNumberBcd[0:])
	return strings.TrimLeft(phoneNumber, "0")
}

// OnInit 初始化回调
func (p *Client) OnInit(ctx *xproto.Conn) {
	p.Conn = ctx
	p.LinkType = xproto.Link_Unknow
}

// OnDecodec 协议解析回调
func (p *Client) OnDecodec(data []byte) (int, error) {
	endPos, err := isValid(data)
	if err != nil || endPos <= 0 {
		return 0, err
	}
	p.jtBytes = toJtBytes(data[1:endPos])
	p.jtMsgPacket = *p.jtBytes.toMsgPacket()
	return endPos + 2, nil
}

func (p *Client) OnIdle(rect, sendt *time.Time) error {
	return nil
}

// OnHandle 协议处理回调
func (p *Client) OnHandle() error {
	log.Printf("code 0x%04x, body length %d\n", p.MsgId, len(p.MsgBody))
	switch p.MsgId {
	case 0x0001:
		p.onRecvResp()
	case 0x0002: // heartbeat
	case 0x0003: // 终端注销
	case 0x0004:
		p.onRecvQueryTime()
	case 0x0100:
		return p.onRevRegister()
	case 0x0102:
		p.onRecvAuthority()
	case 0x0104:
		p.onRecvParameters()
	case 0x0107:
		p.OnRecvTerminalProperty()
	case 0x0108:
		p.OnRecvTerminalUpgrade()
	case 0x0200:
		p.onRecvLoStatus()
	case 0x0201:
		p.onRecvQueryLocation()
	case 0x0302:
		p.onRecvAnswer()
	case 0x0303:
		p.onRecvOnDemand()
	case 0x0500:
		p.onRecvVehiCtrl()
	case 0x0704:
		p.onRecvLoStatusBulk()
	case 0x0800:
		p.onRecvMediaEvent()
	case 0x1205:
		p.onRecvQueryResult()
	case 0x1003:
		p.onRecvAVProperty()
	default:
		log.Printf("code 0x%04x, body length %d\n", p.MsgId, len(p.MsgBody))
	}
	return nil
}

// Request 发送请求
func (p *Client) Request(cmd xproto.ReqCode, s interface{}, r interface{}) error {
	var (
		key uint16 = 0
		err error  = nil
	)
	switch cmd {
	case xproto.Req_Query:
		key, err = p.doQuery(s)
	case xproto.Req_Parameters:
		key, err = p.doParameters(s)
	case xproto.Req_LiveStream:
		key, err = p.doLiveAVStream(s)
	case xproto.Req_Playback:
		key, err = p.doPlayback(s)
	case xproto.Req_FileTransfer:
		key, err = p.doFileTransfer(s)
	case xproto.Req_Control:
		key, err = p.doControl(s)
	case xproto.Req_Voice:
	case xproto.Req_User:
		key, err = p.doJt808(s)
	case xproto.Req_Close:
		p.Conn.Close()
		return nil
	}
	if err != nil {
		return err
	}
	// log.Printf("Request sequence %d\n", key)
	return p.Result.WaitVal(p.sequenceKey(key), r)
}
