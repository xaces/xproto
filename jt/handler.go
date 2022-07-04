package jt

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/xaces/xproto"
	"github.com/xaces/xproto/jt/util"
)

// onRecvResp 通用应答
func (p *Client) onRecvResp() error {
	p.jtBytes.toWord()
	msgId := p.jtBytes.toWord()
	retcode := p.jtBytes.toByte()
	if err := codeError(retcode); err != nil {
		p.Result.SetError(p.sequenceKey(msgId), err)
	} else {
		p.Result.Set(p.sequenceKey(msgId), msgId)
	}
	return nil
}

// onRecvQueryTime 通用应答
func (p *Client) onRecvQueryTime() error {
	return p.writeMessage(0x8004, nil)
}

// onRevRegister 注册
func (p *Client) onRevRegister() (err error) {
	p.terminal = p.jtBytes.toRegister()
	p.DeviceNo = p.phoneNumberString()
	p.Session = xproto.UUID()
	var res jtBytes
	res.addWord(p.MsgSerialNum)
	if err = xproto.SyncAdd(p, p.DeviceNo); err == nil {
		p.LinkType = xproto.Link_Signal
		if err = p.NotifyAccess(p.MsgBody); err == nil {
			res.addByte(0)
			res.addBytes([]byte(jtprotocolAuthority))
			return p.writeMessage(0x8100, res.Bytes)
		}
		res.addByte(2)
	} else {
		res.addByte(1)
	}
	p.writeMessage(0x8100, res.Bytes)
	return
}

// onRecvAuthority 授权
func (p *Client) onRecvAuthority() error {
	return p.writeResp(0)
}

// onRecvParameters 参数应答
func (p *Client) onRecvParameters() error {
	serialNum := p.jtBytes.toWord()
	paramNum := p.jtBytes.toByte()
	var i uint8 = 0
	params := jtParameter{}
	for ; i < paramNum; i++ {
		paramId := p.jtBytes.toDWord()
		length := int(p.jtBytes.toByte())
		fmt.Printf("0x%04x, length %d\n", paramId, length)
		m := jtParameterGet(paramId, p.jtBytes.toArrry(length))
		params.Data = append(params.Data, m)
	}
	p.Result.Set(p.sequenceKey(serialNum), params)
	return p.writeResp(0)
}

func (p *Client) OnRecvTerminalProperty() error {
	res := p.jtBytes.toTerminalProperty()
	p.Result.Set(p.sequenceKey(p.key.Property), res)
	return p.writeResp(0)
}

func (p *Client) OnRecvTerminalUpgrade() error {
	res := p.jtBytes.toTerminalUpgrade()
	p.Result.Set(p.sequenceKey(p.key.Upgrade), res)
	return p.writeResp(0)
}

// onRecvLoStatus 位置信息
func (p *Client) onRecvLoStatus() error {
	b := p.jtBytes
	st, alrflag := p.doStatusParse(b)
	if alrflag > 0 {
		p.NotifyAlarm(nil, p.toAlarm(alrflag, &st, b))
	} else {
		p.NotifyStatus("jGps", &st)
	}
	return p.writeResp(0)
}

// onRecvQueryLocation 查询位置
func (p *Client) onRecvQueryLocation() error {
	b := p.jtBytes
	sequence := b.toWord()
	st, _ := p.doStatusParse(b)
	p.Result.Set(p.sequenceKey(sequence), st)
	return p.writeResp(0)
}

// onRecvAnswer 问题应答
func (p *Client) onRecvAnswer() error {
	sequence := p.jtBytes.toWord()
	AnswerId := p.jtBytes.toByte()
	p.Result.Set(p.sequenceKey(sequence), map[string]interface{}{"answerId": AnswerId})
	return p.writeResp(0)
}

// 信息点播/取消
func (p *Client) onRecvOnDemand() error {
	b := p.jtBytes
	sequence := b.toWord()
	st, _ := p.doStatusParse(b)
	p.Result.Set(p.sequenceKey(sequence), st)
	return p.writeResp(0)
}

// onRecvVehiCtrl 车辆控制应答
func (p *Client) onRecvVehiCtrl() error {
	msgType := p.jtBytes.toByte()
	flag := p.jtBytes.toByte()
	p.Result.Set(p.sequenceKey(p.key.OnDemand), map[string]interface{}{"demandType": msgType, "demandFlag": flag})
	return p.writeResp(0)
}

// onRecvLoStatusBulk 批量上报
func (p *Client) onRecvLoStatusBulk() error {
	var jbs jtBulkStatus
	b := p.jtBytes
	jbs.DataNumber = b.toWord()
	jbs.BlindSpot = b.toByte()
	for i := 0; i < int(jbs.DataNumber); i++ {
		jbs.Length = b.toWord()
		st, alrflag := p.doStatusParse(b)
		if alrflag > 0 {
			p.NotifyAlarm(nil, p.toAlarm(alrflag, &st, b))
			continue
		}
		st.Flag = 1
		p.NotifyStatus("jGps", &st)
	}
	return p.writeResp(0)
}

// onRecvMediaEvent 多媒体事件信息
func (p *Client) onRecvMediaEvent() error {
	b := p.jtBytes
	id := b.toDWord()
	b.toByte()
	b.toByte() // 多媒体格式编码
	// switch  {
	// case 0:
	// 	"JPEG"
	// case 1:
	// 	"TIF"
	// case 2:
	// 	"MP3"
	// case 3:
	// 	"WAV"
	// case "4":
	// 	"WMV"
	// }
	b.toByte()
	b.toByte()
	p.doStatusParse(b)
	// write to file
	var res jtBytes
	res.addDWord(id)
	return p.writeMessage(0x8800, res.Bytes)
}

// onRecvQueryResult 查询结果
func (p *Client) onRecvQueryResult() error {
	b := p.jtBytes
	sequence := b.toWord()
	res := p.Result.Get(sequence)
	// 持续接收数据
	v := reflect.ValueOf(res.Data).Interface().(*[]xproto.File)
	count := int(b.toWord())
	for i := 0; i < count; i++ {
		fl := xproto.File{}
		fl.Channel = int(b.toByte())
		fl.StartTime = dtu(b.toBCD(6))
		fl.EndTime = dtu(b.toBCD(6))
		// 计算时长
		fl.FileDuration = int(util.ToUnixTime(fl.EndTime) - util.ToUnixTime(fl.StartTime))
		alrType := binary.BigEndian.Uint64(b.toArrry(8))
		fptype := b.toByte()
		if fptype == 2 {
			fl.FileType = xproto.File_NormalVideo
			if alrType > 0 {
				fl.FileType = xproto.File_AlarmVideo
			}
		} else if fptype == 1 {
			fl.FileType = xproto.File_NormalVoice
			if alrType > 0 {
				fl.FileType = xproto.File_AlarmVoice
			}
		}
		b.toByte() // 码流类型
		b.toByte()
		fl.FileSize = int(b.toDWord())
		if res == nil {
			p.NotityFile(&fl)
			continue
		}
		*v = append(*v, fl)
		if i+1 == count {
			res.Done <- true
		}
	}
	return nil
}

func (p *Client) onRecvAVProperty() error {
	return nil
}
