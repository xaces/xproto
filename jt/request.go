package jt

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/xaces/xproto"
	"github.com/xaces/xproto/jt/util"
)

func (p *Client) doParameters(s interface{}) (uint16, error) {
	data, _ := jsoniter.Marshal(s)
	var params jtParameter
	jsoniter.Get(data, "jtparameter").ToVal(&params.Data)
	paramLen := len(params.Data)
	if paramLen == 0 {
		return p.sequence, p.writeMessage(0x8104, nil)
	}
	var (
		res   jtBytes
		isGet bool = false
	)
	if params.Data[0].Value == nil {
		isGet = true
	}
	res.addByte(uint8(paramLen))
	for _, v := range params.Data {
		id := uint32(util.HexStringToUint16(v.Id))
		res.addDWord(id)
		if isGet {
			continue
		}
		jtParameterSet(id, &res, v.Value)
	}
	if isGet {
		return p.sequence, p.writeMessage(0x8106, res.Bytes)
	}
	return 0, p.writeMessage(0x8103, res.Bytes)
}

// doFileTransfer 文件传输
func (p *Client) doFileTransfer(s interface{}) (uint16, error) {
	v, ok := s.(*xproto.FileTransfer)
	if !ok {
		return 0, xproto.ErrParam
	}
	var res jtBytes
	if v.Action == xproto.Act_Upload && v.FileType == xproto.File_Upgrade {
		fp, err := os.Open(v.FileName)
		if err != nil {
			return 0, err
		}
		defer fp.Close()
		fileName := filepath.Base(v.FileName)
		res.addByte(0) // TODO判断类型
		res.addBytes(p.terminal.ManufacturerId)
		fnbytes := util.ToGBK(fileName)
		res.addByte(uint8(len(fnbytes)))
		res.addBytes(fnbytes)
		i, _ := fp.Stat()
		res.addDWord(uint32(i.Size()))
		bs := make([]byte, 10*1024)
		for {
			s, err := fp.Read(bs)
			if err == io.EOF {
				break
			}
			res.addBytes(bs[:s])
		}
		// 升级包数据
		if err := p.writeMessage(0x8108, res.Bytes); err != nil {
			return 0, err
		}
		p.key.Upgrade = p.sequence
		return p.key.Upgrade, nil
	}
	return 0, xproto.ErrParam
}

func (p *Client) doLiveAVStream(s interface{}) (uint16, error) {
	v, ok := s.(*xproto.LiveStream)
	if !ok || v.Server == "" {
		return 0, xproto.ErrParam
	}
	addr := strings.Split(v.Server, ":")
	if len(addr) != 2 {
		return 0, xproto.ErrParam
	}
	var res jtBytes
	res.addByte(uint8(len(addr[0])))
	res.addString(addr[0])
	port, _ := strconv.Atoi(addr[1])
	res.addWord(uint16(port))
	res.addWord(0)
	res.addByte(uint8(v.Channel))
	res.addByte(0)
	res.addByte(uint8(v.StreamType))
	return p.sequence, p.writeMessage(0x9101, res.Bytes)
}

// 语音
func (p *Client) doLiveControl(s interface{}) (uint16, error) {
	return 0, xproto.ErrUnSupport
}

func (p *Client) doPlayback(s interface{}) (uint16, error) {
	v, ok := s.(*xproto.Playback)
	if !ok || v.Server == "" {
		return 0, xproto.ErrParam
	}
	addr := strings.Split(v.Server, ":")
	if len(addr) != 2 {
		return 0, xproto.ErrParam
	}
	var res jtBytes
	res.addByte(uint8(len(addr[0])))
	res.addString(addr[0])
	port, _ := strconv.Atoi(addr[1])
	res.addWord(uint16(port))
	res.addWord(0)
	chls := strings.Split(v.Channels, ";")
	if len(chls) > 1 {
		res.addString(chls[0])
	} else {
		res.addString(v.Channels)
	}
	res.addByte(0)
	res.addByte(0)
	res.addByte(0)
	res.addByte(0) //
	res.addByte(0) //
	res.addByte(0) //
	res.addDateTime(v.StartTime)
	res.addDateTime(v.EndTime)
	return p.sequence, p.writeMessage(0x9201, res.Bytes)
}

//{"codeId": "0x8107", "data": {}}

func (p *Client) doJt808(s interface{}) (uint16, error) {
	c, ok := s.(*xproto.User)
	if !ok {
		return 0, xproto.ErrParam
	}
	codeId := util.HexStringToUint16(c.CodeId)
	if codeId == 0x8201 {
		return p.sequence, p.writeMessage(codeId, nil)
	} else if codeId == 0x8107 {
		if err := p.writeMessage(codeId, nil); err != nil {
			return 0, err
		}
		p.key.Property = p.sequence
		return p.key.Property, nil
	}
	sdata, _ := jsoniter.Marshal(s)
	iter := jsoniter.Get(sdata, "data")
	if err := iter.LastError(); err != nil {
		return 0, err
	}
	var res jtBytes
	switch codeId {
	case 0x8202:
		var param Jt8202
		iter.ToVal(&param)
		res.addWord(param.Interval)
		res.addDWord(param.EffectiveTime)
	case 0x8203:
		var param Jt8203
		iter.ToVal(&param)
		res.addWord(param.AlarmSerial)
		res.addDWord(param.AlarmType)
	case 0x8302:
		var param Jt8302
		iter.ToVal(&param)
		res.addByte(param.Flag)
		question := util.ToGBK(param.Question)
		res.addByte(uint8(len(question)))
		res.addBytes(question)
		for _, v := range param.Answers {
			res.addByte(v.Id)
			content := util.ToGBK(v.Content)
			res.addByte(uint8(len(content)))
			res.addBytes(content)
		}
		if err := p.writeMessage(codeId, res.Bytes); err != nil {
			return 0, err
		}
		p.key.OnDemand = p.sequence
		return p.key.OnDemand, nil
	case 0x8303:
		var param Jt8303
		iter.ToVal(&param)
		res.addByte(param.Type)
		res.addByte(uint8(len(param.Message)))
		for _, v := range param.Message {
			res.addByte(v.Type)
			name := util.ToGBK(v.Name)
			res.addWord(uint16(len(name)))
			res.addBytes(name)
		}
	case 0x8304:
		var param Jt8304
		iter.ToVal(&param)
		res.addByte(param.Type)
		content := util.ToGBK(param.Content)
		res.addWord(uint16(len(content)))
		res.addBytes(content)
	case 0x8400:
		var param Jt8400
		iter.ToVal(&param)
		res.addByte(param.Type)
		res.addString(param.PhoneNumber) // TODO check max 20
	case 0x8401:
		var param Jt8401
		iter.ToVal(&param)
		res.addByte(param.Type)
		res.addByte(uint8(len(param.Contacts)))
		for _, v := range param.Contacts {
			res.addByte(v.Type)
			pn := util.ToGBK(v.PhoneNumber)
			res.addByte(uint8(len(pn)))
			res.addBytes(pn)
			name := util.ToGBK(v.Name)
			res.addByte(uint8(len(name)))
			res.addBytes(name)
		}
	case 0x8600:
		var param Jt8600
		iter.ToVal(&param)
		res.addByte(param.Action)
		res.addByte(uint8(len(param.Regions)))
		for _, v := range param.Regions {
			res.addDWord(v.Id)
			res.addWord(v.Attr)
			res.addDWord(uint32(v.La * 1000000))
			res.addDWord(uint32(v.Lo * 1000000))
			res.addDWord(v.Ra)
			res.addDateTime(v.St)
			res.addDateTime(v.Et)
			res.addWord(v.MaxSpd)
			res.addByte(v.KeepTime)
			res.addWord(v.MaxNightSpd)
			res.addWord(uint16(len(v.Name)))
			res.addString(v.Name)
		}
	case 0x8602:
		var param Jt8602
		iter.ToVal(&param)
		res.addByte(param.Action)
		res.addByte(uint8(len(param.Regions)))
		for _, v := range param.Regions {
			res.addDWord(v.Id)
			res.addWord(v.Attr)
			res.addDWord(uint32(v.LeftTopLa * 1000000))
			res.addDWord(uint32(v.LeftTopLo * 1000000))
			res.addDWord(uint32(v.RightBotLa * 1000000))
			res.addDWord(uint32(v.RightBotLo * 1000000))
			res.addDateTime(v.St)
			res.addDateTime(v.Et)
			res.addWord(v.MaxSpd)
			res.addByte(v.KeepTime)
			res.addWord(v.MaxNightSpd)
			res.addWord(uint16(len(v.Name)))
			res.addString(v.Name)
		}
	case 0x8604:
		var param Jt8604
		iter.ToVal(&param)
		res.addDWord(param.Id)
		res.addWord(param.Attr)
		res.addDateTime(param.St)
		res.addDateTime(param.Et)
		res.addWord(param.MaxSpd)
		res.addByte(param.KeepTime)
		res.addWord(uint16(len(param.Dots)))
		for _, v := range param.Dots {
			res.addDWord(uint32(v.La * 1000000))
			res.addDWord(uint32(v.Lo * 1000000))
		}
		res.addWord(param.MaxNightSpd)
		res.addWord(uint16(len(param.Name)))
		res.addString(param.Name)
	case 0x8601, 0x8603, 0x8605:
		var param Jt8601O3
		iter.ToVal(&param)
		res.addByte(uint8(len(param.Ids)))
		for _, v := range param.Ids {
			res.addDWord(v)
		}
	default:
		return 0, nil
	}
	return p.sequence, p.writeMessage(codeId, res.Bytes)
}

func (p *Client) doControl(s interface{}) (uint16, error) {
	v, ok := s.(*xproto.Control)
	if !ok {
		return 0, xproto.ErrParam
	}
	data, _ := jsoniter.Marshal(v.Data)
	j := jsoniter.Get(data)
	var res jtBytes
	switch v.Type {
	case xproto.Ctrl_Text:
		var txt Jt8300
		if err := jsoniter.Unmarshal(data, &txt); err != nil {
			return 0, err
		}
		res.addByte(txt.Flag)
		res.addString(txt.Text)
		return p.sequence, p.writeMessage(0x8300, res.Bytes)
	case xproto.Ctrl_Vehi:
		code := xproto.VehiCode(j.Get("code").ToInt())
		if code == xproto.Ctrl_VehiOpenDoor {
			res.addByte(0x00)
		} else if code == xproto.Ctrl_VehiCloseDoor {
			res.addByte(0x01)
		} else {
			return 0, xproto.ErrParam
		}
		return p.sequence, p.writeMessage(0x8500, res.Bytes)
	case xproto.Ctrl_Reset:
		res.addByte(4)
		return 0, p.writeMessage(0x8105, res.Bytes)
	case xproto.Ctrl_Capture:
	case xproto.Ctrl_Live:
		var param xproto.LiveControl
		j.ToVal(param)
		res.addByte(param.Channel)
		res.addByte(param.Type)
		res.addByte(param.CloseMeidaType)
		res.addByte(param.ToggleStreamType)
		return 0, p.writeMessage(0x9102, res.Bytes)
	}
	return 0, nil
}

// 查询请求
func (p *Client) doQuery(s interface{}) (uint16, error) {
	v, ok := s.(*xproto.Query)
	if !ok {
		return 0, xproto.ErrParam
	}
	var data jtBytes
	// 这里判断如果是1;2;3;4多于1个通道搜索全部
	pos := strings.Index(v.ChannelList, ";")
	if pos > 0 {
		data.addByte(0)
	} else {
		ch, _ := strconv.Atoi(v.ChannelList)
		data.addByte(uint8(ch))
	}
	data.addDateTime(v.StartTime)
	data.addDateTime(v.EndTime)
	if v.FileType == xproto.File_AlarmVideo || v.FileType == xproto.File_AlarmVoice {
		data.add64BITS(0xffffffff)
	} else {
		data.add64BITS(0)
	}
	switch v.FileType {
	case xproto.File_NormalVideo, xproto.File_AlarmVideo:
		data.addByte(2)
	// case xproto.File_NormalPic, xproto.File_AlarmPic:
	// 	data.addByte(0)
	case xproto.File_NormalVoice, xproto.File_AlarmVoice:
		data.addByte(1)
	default:
		return 0, xproto.ErrParam
	}
	data.addByte(uint8(v.StreamType))
	data.addByte(uint8(v.StoreType))
	return p.sequence, p.writeMessage(0x9205, data.Bytes)
}
