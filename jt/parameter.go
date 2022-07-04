package jt

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/xaces/xproto/jt/util"
)

type jtParamsData struct {
	Id    string      `json:"codeId"`
	Value interface{} `json:"value,omitempty"`
}

type jtParameter struct {
	Data []jtParamsData `json:"jtparameter"`
}

type jtChannelParam struct {
	Live struct {
		Encode        uint8  `json:"encode"`
		Res           uint8  `json:"res"`
		FrameInterval uint16 `json:"frameInterval"`
		FrameRate     uint8  `json:"frameRate"`
		BitRate       uint16 `json:"bitRate"`
	} `json:"live"`
	Storage struct {
		Encode        uint8  `json:"encode"`
		Res           uint8  `json:"res"`
		FrameInterval uint16 `json:"frameInterval"`
		FrameRate     uint8  `json:"frameRate"`
		BitRate       uint16 `json:"bitRate"`
	} `json:"storage"`
	Osd uint16 `json:"osd"`
}

// 音视频参数设置
type jt0075 struct {
	jtChannelParam
	AudioOut uint8 `json:"audioOut"`
}

type jt0076Channel struct {
	Physical uint8 `json:"physical"`
	Logical  uint8 `json:"logical"`
	Type     uint8 `json:"type"`
	Cloud    uint8 `json:"cloud"`
}

// 音视频通道列表设置
type jt0076 struct {
	Total    uint8           `json:"total"`
	ANum     uint8           `json:"anum"`
	VNum     uint8           `json:"vnum"`
	Channels []jt0076Channel `json:"channels"`
}

// 单独通道视频参数
type jt0077 struct {
	SetChannels []jtChannelParam `json:"setChannels"`
}

// 特殊报警录像参数设置
type jt0079 struct {
	StoThreshold uint8 `json:"stoThreshold"`
	Duration     uint8 `json:"duration"`
	StartTime    uint8 `json:"startTime"`
}

// 视频分析报警参数定义
type jt007B struct {
	PeopleNum    uint8 `json:"peopleNum"`
	FatThreshold uint8 `json:"fatThreshold"`
	StartTime    uint8 `json:"startTime"`
}

// 终端休眠唤醒模式
type jt007C struct {
	WakeMode  uint8 `json:"wakeMode"`
	WakeType  uint8 `json:"wakeType"`
	WakeTimed uint8 `json:"wakeTimed"`
}

func jtParameterGet(id uint32, data []byte) jtParamsData {
	m := jtParamsData{}
	m.Id = fmt.Sprintf("0x%04x", id)
	jb := &jtBytes{pos: 0, Bytes: data}
	switch id {
	case 0x0001:
		m.Value = jb.toDWord()
	case 0x0002:
		m.Value = jb.toDWord()
	case 0x0003:
		m.Value = jb.toDWord()
	case 0x0004:
		m.Value = jb.toDWord()
	case 0x0005:
		m.Value = jb.toDWord()
	case 0x0006:
		m.Value = jb.toDWord()
	case 0x0007:
		m.Value = jb.toDWord()
	case 0x0010:
		m.Value = util.ToUTF8(data)
	case 0x0011:
		m.Value = util.ToUTF8(data)
	case 0x0012:
		m.Value = util.ToUTF8(data)
	case 0x0013:
		m.Value = util.ToUTF8(data)
	case 0x0014:
		m.Value = util.ToUTF8(data)
	case 0x0015:
		m.Value = util.ToUTF8(data)
	case 0x0016:
		m.Value = util.ToUTF8(data)
	case 0x0017:
		m.Value = util.ToUTF8(data)
	case 0x0018:
		m.Value = jb.toDWord()
	case 0x0019:
		m.Value = jb.toDWord()
	case 0x001A:
		m.Value = util.ToUTF8(data)
	case 0x001B:
		m.Value = jb.toDWord()
	case 0x001C:
		m.Value = jb.toDWord()
	case 0x001D:
		m.Value = util.ToUTF8(data)
	case 0x0020:
		m.Value = jb.toDWord()
	case 0x0021:
		m.Value = jb.toDWord()
	case 0x0022:
		m.Value = jb.toDWord()
	case 0x0023:
		m.Value = jb.toDWord()
	case 0x0024:
		m.Value = jb.toDWord()
	case 0x0025:
		m.Value = jb.toDWord()
	case 0x0026:
		m.Value = jb.toDWord()
	case 0x0027:
		m.Value = jb.toDWord()
	case 0x0028:
		m.Value = jb.toDWord()
	case 0x0029:
		m.Value = jb.toDWord()
	case 0x002A:
		m.Value = jb.toDWord()
	case 0x002B:
		m.Value = jb.toDWord()
	case 0x002C:
		m.Value = jb.toDWord()
	case 0x002D:
		m.Value = jb.toDWord()
	case 0x002E:
		m.Value = jb.toDWord()
	case 0x002F:
		m.Value = jb.toDWord()
	case 0x0030:
		m.Value = jb.toDWord()
	case 0x0031:
		m.Value = jb.toWord()
	case 0x0040:
		m.Value = util.ToUTF8(data)
	case 0x0041:
		m.Value = util.ToUTF8(data)
	case 0x0042:
		m.Value = util.ToUTF8(data)
	case 0x0043:
		m.Value = util.ToUTF8(data)
	case 0x0044:
		m.Value = util.ToUTF8(data)
	case 0x0045:
		m.Value = jb.toDWord()
	case 0x0046:
		m.Value = jb.toDWord()
	case 0x0047:
		m.Value = jb.toDWord()
	case 0x0048:
		m.Value = util.ToUTF8(data)
	case 0x0049:
		m.Value = util.ToUTF8(data)
	case 0x0050:
		m.Value = jb.toDWord()
	case 0x0051:
		m.Value = jb.toDWord()
	case 0x0052:
		m.Value = jb.toDWord()
	case 0x0053:
		m.Value = jb.toDWord()
	case 0x0054:
		m.Value = jb.toDWord()
	case 0x0055:
		m.Value = jb.toDWord()
	case 0x0056:
		m.Value = jb.toDWord()
	case 0x0057:
		m.Value = jb.toDWord()
	case 0x0058:
		m.Value = jb.toDWord()
	case 0x0059:
		m.Value = jb.toDWord()
	case 0x005A:
		m.Value = jb.toDWord()
	case 0x005B:
		m.Value = jb.toWord()
	case 0x005C:
		m.Value = jb.toWord()
	case 0x005D:
		m.Value = jb.toWord()
	case 0x005E:
		m.Value = jb.toWord()
	case 0x0064:
		m.Value = jb.toDWord()
	case 0x0065:
		m.Value = jb.toDWord()
	case 0x0070:
		m.Value = jb.toDWord()
	case 0x0071:
		m.Value = jb.toDWord()
	case 0x0072:
		m.Value = jb.toDWord()
	case 0x0073:
		m.Value = jb.toDWord()
	case 0x0074:
		m.Value = jb.toDWord()
	case 0x0075:
		var param jt0075
		param.Live.Encode = jb.toByte()
		param.Live.Res = jb.toByte()
		param.Live.FrameInterval = jb.toWord()
		param.Live.FrameRate = jb.toByte()
		param.Live.BitRate = jb.toWord()
		param.Storage.Encode = jb.toByte()
		param.Storage.Res = jb.toByte()
		param.Storage.FrameInterval = jb.toWord()
		param.Storage.FrameRate = jb.toByte()
		param.Storage.BitRate = jb.toWord()
		param.Osd = jb.toWord()
		param.AudioOut = jb.toByte()
		m.Value = param
	case 0x0076:
		var param jt0076
		param.Total = jb.toByte()
		param.ANum = jb.toByte()
		param.VNum = jb.toByte()
		chns := int(param.Total + param.ANum + param.VNum)
		for i := 0; i < chns; i++ {
			ch := jt0076Channel{}
			ch.Physical = jb.toByte()
			ch.Logical = jb.toByte()
			ch.Type = jb.toByte()
			ch.Cloud = jb.toByte()
			param.Channels = append(param.Channels, ch)
		}
		m.Value = param
	case 0x0077:
		count := int(jb.toByte())
		var param jt0077
		for i := 0; i < count; i++ {
			var ch jtChannelParam
			ch.Live.Encode = jb.toByte()
			ch.Live.Res = jb.toByte()
			ch.Live.FrameInterval = jb.toWord()
			ch.Live.FrameRate = jb.toByte()
			ch.Live.BitRate = jb.toWord()
			ch.Storage.Encode = jb.toByte()
			ch.Storage.Res = jb.toByte()
			ch.Storage.FrameInterval = jb.toWord()
			ch.Storage.FrameRate = jb.toByte()
			ch.Storage.BitRate = jb.toWord()
			ch.Osd = jb.toWord()
			param.SetChannels = append(param.SetChannels, ch)
		}
		m.Value = param
	case 0x0079:
		var param jt0079
		param.StoThreshold = jb.toByte()
		param.Duration = jb.toByte()
		param.StartTime = jb.toByte()
		m.Value = param
	case 0x007A:
		m.Value = jb.toDWord()
	case 0x007B:
		var param jt007B
		param.PeopleNum = jb.toByte()
		param.FatThreshold = jb.toByte()
		param.StartTime = jb.toByte()
		m.Value = param
	case 0x007C:
		var param jt007C
		param.WakeMode = jb.toByte()
		param.WakeType = jb.toByte()
		param.WakeTimed = jb.toByte()
		m.Value = param
	case 0x0080:
		m.Value = jb.toDWord()
	case 0x0081:
		m.Value = jb.toByte()
	case 0x0082:
		m.Value = jb.toByte()
	case 0x0083:
		m.Value = util.ToUTF8(data)
	case 0x0084:
		m.Value = jb.toByte()
	case 0x0090:
		m.Value = jb.toByte()
	case 0x0091:
		m.Value = jb.toByte()
	case 0x0092:
		m.Value = jb.toByte()
	case 0x0093:
		m.Value = jb.toDWord()
	case 0x0094:
		m.Value = jb.toByte()
	case 0x0095:
		m.Value = jb.toDWord()
	case 0x0100:
		m.Value = jb.toDWord()
	case 0x0101:
		m.Value = jb.toWord()
	case 0x0102:
		m.Value = jb.toDWord()
	case 0x0103:
		m.Value = jb.toWord()
	case 0x0110:
		m.Value = data
	default:
		m.Value = data
	}
	return m
}

func setByte(b *jtBytes, val interface{}) {
	b.addByte(1)
	data := reflect.ValueOf(val).Interface().(float64)
	b.addByte(uint8(data))
}

func setWord(b *jtBytes, val interface{}) {
	b.addByte(2)
	data := reflect.ValueOf(val).Interface().(float64)
	b.addWord(uint16(data))
}

func setDWord(b *jtBytes, val interface{}) {
	b.addByte(4)
	data := reflect.ValueOf(val).Interface().(float64)
	b.addDWord(uint32(data))
}

func setString(b *jtBytes, val interface{}) {
	data := reflect.ValueOf(val).Interface().(string)
	b.addByte(uint8(len(data)))
	b.addString(data)
}

func setByteArray(b *jtBytes, val interface{}) {
	data := reflect.ValueOf(val).Interface().([]byte)
	b.addByte(uint8(len(data)))
	b.addBytes(data)
}

func jtParameterSet(id uint32, jb *jtBytes, val interface{}) {
	if val == nil {
		return
	}
	switch id {
	case 0x0001:
		setDWord(jb, val)
	case 0x0002:
		setDWord(jb, val)
	case 0x0003:
		setDWord(jb, val)
	case 0x0004:
		setDWord(jb, val)
	case 0x0005:
		setDWord(jb, val)
	case 0x0006:
		setDWord(jb, val)
	case 0x0007:
		setDWord(jb, val)
	case 0x0010:
		setString(jb, val)
	case 0x0011:
		setString(jb, val)
	case 0x0012:
		setString(jb, val)
	case 0x0013:
		setString(jb, val)
	case 0x0014:
		setString(jb, val)
	case 0x0015:
		setString(jb, val)
	case 0x0016:
		setString(jb, val)
	case 0x0017:
		setString(jb, val)
	case 0x0018:
		setDWord(jb, val)
	case 0x0019:
		setDWord(jb, val)
	case 0x001A:
		setString(jb, val)
	case 0x001B:
		setDWord(jb, val)
	case 0x001C:
		setDWord(jb, val)
	case 0x001D:
		setString(jb, val)
	case 0x0020:
		setDWord(jb, val)
	case 0x0021:
		setDWord(jb, val)
	case 0x0022:
		setDWord(jb, val)
	case 0x0023:
		setDWord(jb, val)
	case 0x0024:
		setDWord(jb, val)
	case 0x0025:
		setDWord(jb, val)
	case 0x0026:
		setDWord(jb, val)
	case 0x0027:
		setDWord(jb, val)
	case 0x0028:
		setDWord(jb, val)
	case 0x0029:
		setDWord(jb, val)
	case 0x002A:
		setDWord(jb, val)
	case 0x002B:
		setDWord(jb, val)
	case 0x002C:
		setDWord(jb, val)
	case 0x002D:
		setDWord(jb, val)
	case 0x002E:
		setDWord(jb, val)
	case 0x002F:
		setDWord(jb, val)
	case 0x0030:
		setDWord(jb, val)
	case 0x0031:
		setWord(jb, val)
	case 0x0040:
		setString(jb, val)
	case 0x0041:
		setString(jb, val)
	case 0x0042:
		setString(jb, val)
	case 0x0043:
		setString(jb, val)
	case 0x0044:
		setString(jb, val)
	case 0x0045:
		setDWord(jb, val)
	case 0x0046:
		setDWord(jb, val)
	case 0x0047:
		setDWord(jb, val)
	case 0x0048:
		setString(jb, val)
	case 0x0049:
		setString(jb, val)
	case 0x0050:
		setDWord(jb, val)
	case 0x0051:
		setDWord(jb, val)
	case 0x0052:
		setDWord(jb, val)
	case 0x0053:
		setDWord(jb, val)
	case 0x0054:
		setDWord(jb, val)
	case 0x0055:
		setDWord(jb, val)
	case 0x0056:
		setDWord(jb, val)
	case 0x0057:
		setDWord(jb, val)
	case 0x0058:
		setDWord(jb, val)
	case 0x0059:
		setDWord(jb, val)
	case 0x005A:
		setDWord(jb, val)
	case 0x005B:
		setWord(jb, val)
	case 0x005C:
		setWord(jb, val)
	case 0x005D:
		setWord(jb, val)
	case 0x005E:
		setWord(jb, val)
	case 0x0064:
		setDWord(jb, val)
	case 0x0065:
		setDWord(jb, val)
	case 0x0070:
		setDWord(jb, val)
	case 0x0071:
		setDWord(jb, val)
	case 0x0072:
		setDWord(jb, val)
	case 0x0073:
		setDWord(jb, val)
	case 0x0074:
		setDWord(jb, val)
	case 0x0075:
		data, _ := json.Marshal(val)
		var param jt0075
		if err := json.Unmarshal(data, &param); err != nil {
			return
		}
		jb.addByte(param.Live.Encode)
		jb.addByte(param.Live.Res)
		jb.addWord(param.Live.FrameInterval)
		jb.addByte(param.Live.FrameRate)
		jb.addWord(param.Live.BitRate)
		jb.addByte(param.Storage.Encode)
		jb.addByte(param.Storage.Res)
		jb.addWord(param.Storage.FrameInterval)
		jb.addByte(param.Storage.FrameRate)
		jb.addWord(param.Storage.BitRate)
		jb.addWord(param.Osd)
		jb.addByte(param.AudioOut)
	case 0x0080:
		setDWord(jb, val)
	case 0x0081:
		setByte(jb, val)
	case 0x0082:
		setByte(jb, val)
	case 0x0083:
		setString(jb, val)
	case 0x0084:
		setByte(jb, val)
	case 0x0090:
		setByte(jb, val)
	case 0x0091:
		setByte(jb, val)
	case 0x0092:
		setByte(jb, val)
	case 0x0093:
		setDWord(jb, val)
	case 0x0094:
		setByte(jb, val)
	case 0x0095:
		setDWord(jb, val)
	case 0x0100:
		setDWord(jb, val)
	case 0x0101:
		setWord(jb, val)
	case 0x0102:
		setDWord(jb, val)
	case 0x0103:
		setWord(jb, val)
	case 0x0110:
		setByteArray(jb, val)
	}
}
