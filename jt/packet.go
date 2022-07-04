package jt

import (
	"encoding/binary"

	"github.com/xaces/xproto/jt/util"
)

type jtBytes struct {
	jtVer uint8 // jt对应版本 1: jt808-2019
	pos   int
	Bytes []byte
}

// addDateTime 添加时间
// 2006-01-02 15:04:05 -> 060102150405 -> Bcd
func (b *jtBytes) addDateTime(v string) {
	tm := v[2:4] + v[5:7] + v[8:10] + v[11:13] + v[14:16] + v[17:]
	b.addBCD(util.Number2bcd(tm))
}

func (b *jtBytes) addByte(v uint8) {
	b.Bytes = append(b.Bytes, v)
}

func (b *jtBytes) addWord(v uint16) {
	b.Bytes = append(b.Bytes, byte(v>>8))
	b.Bytes = append(b.Bytes, byte(v))
}

func (b *jtBytes) addDWord(v uint32) {
	b.Bytes = append(b.Bytes, byte(v>>24))
	b.Bytes = append(b.Bytes, byte(v>>16))
	b.Bytes = append(b.Bytes, byte(v>>8))
	b.Bytes = append(b.Bytes, byte(v))
}

func (b *jtBytes) add64BITS(v uint64) {
	b.Bytes = append(b.Bytes, byte(v>>56))
	b.Bytes = append(b.Bytes, byte(v>>48))
	b.Bytes = append(b.Bytes, byte(v>>40))
	b.Bytes = append(b.Bytes, byte(v>>32))
	b.Bytes = append(b.Bytes, byte(v>>24))
	b.Bytes = append(b.Bytes, byte(v>>16))
	b.Bytes = append(b.Bytes, byte(v>>8))
	b.Bytes = append(b.Bytes, byte(v))
}

func (b *jtBytes) addBCD(bcd []byte) {
	b.addBytes(bcd)
}

func (b *jtBytes) addBytes(bs []byte) {
	b.Bytes = append(b.Bytes, bs...)
}

func (b *jtBytes) addString(data string) {
	bs := util.ToGBK(data)
	b.Bytes = append(b.Bytes, bs...)
}

func (b *jtBytes) toByte() uint8 {
	val := b.Bytes[b.pos]
	b.pos++
	return val
}

func (b *jtBytes) toWord() uint16 {
	val := binary.BigEndian.Uint16(b.Bytes[b.pos:])
	b.pos += 2
	return val
}

func (b *jtBytes) toDWord() uint32 {
	val := binary.BigEndian.Uint32(b.Bytes[b.pos:])
	b.pos += 4
	return val
}

func (b *jtBytes) toArrry(size int) []byte {
	b.pos += size
	return b.Bytes[b.pos-size : b.pos]
}

func (b *jtBytes) toBCD(size int) []byte {
	b.pos += size
	return b.Bytes[b.pos-size : b.pos]
}

func (b *jtBytes) toBytes() []byte {
	return b.Bytes[b.pos:]
}

const (
	jtprotocolFlag         = 0x7e
	jtprotocolHeaderLength = 12
	jtprotocolAuthority    = "abcdef"
)

func (p *Client) jtPacketBytes(msgId uint16, b []byte) []byte {
	srcLen := len(b)
	var data jtBytes
	data.addWord(msgId)
	data.addWord(uint16(srcLen))
	data.addBCD(p.PhoneNumberBcd)
	p.sequence++
	data.addWord(p.sequence)
	data.addBytes(b)
	// 计算checksum
	data.addByte(util.CheckSum(data.Bytes))
	var encBytes []byte
	encBytes = append(encBytes, 0x7E)
	for _, v := range data.Bytes {
		if v == 0x7D {
			encBytes = append(encBytes, 0x7D)
			encBytes = append(encBytes, 0x01)
		} else if v == 0x7E {
			encBytes = append(encBytes, 0x7D)
			encBytes = append(encBytes, 0x02)
		} else {
			encBytes = append(encBytes, v)
		}
	}
	encBytes = append(encBytes, 0x7E)
	return encBytes
}

// 解析协议头
func toJtBytes(b []byte) *jtBytes {
	// 转码
	srcLen := len(b)
	decBytes := make([]byte, srcLen)
	decIdx := 0
	for i := 0; i < srcLen-1; i++ {
		if b[i] == 0x7d && b[i+1] == 0x01 {
			decBytes[decIdx] = 0x7d
			i++
		} else if b[i] == 0x7d && b[i+1] == 0x02 {
			decBytes[decIdx] = 0x7e
			i++
		} else {
			decBytes[decIdx] = b[i]
		}
		decIdx++
	}
	return &jtBytes{pos: 0, Bytes: decBytes}
}

// jtMsgProperty 属性
// 0-9bit 字节消息体长度
// 10-12bit 数据加密方式
// 13bit 分包
// 14-15bit 保留 （jt808-2019 bit14 版本表示）
type jtMsgProperty struct {
	Length int
	// 10-12bit
	// 三位都为0  消息体不加密
	// 10bit 为1 消息体RSA加密
	Encryption uint8
	SubFlag    uint8
	VerFlag    uint8
}

func toMsgProperty(property uint16) *jtMsgProperty {
	msgProperty := &jtMsgProperty{
		Length:     int(property & 0x03FF),
		Encryption: uint8((property >> 10) & 0x07),
		SubFlag:    uint8((property >> 13) & 0x01),
		VerFlag:    uint8((property >> 14) & 0x01),
	}
	return msgProperty
}

// jtMsgPacket
// 大端对齐
type jtMsgPacket struct {
	MsgId          uint16
	MsgProperty    jtMsgProperty
	PhoneNumberBcd []byte //6位bcd 响应时无需编码
	MsgSerialNum   uint16
	MsgTotal       uint16
	MsgCur         uint16
	MsgBody        []byte
	Version        uint8
}

func (b *jtBytes) toMsgPacket() *jtMsgPacket {
	pkt := &jtMsgPacket{}
	pkt.MsgId = b.toWord()
	pkt.MsgProperty = *toMsgProperty(b.toWord())
	if pkt.MsgProperty.VerFlag == 1 { // jt808-2019
		b.jtVer = pkt.MsgProperty.VerFlag
		pkt.Version = b.toByte()
		pkt.PhoneNumberBcd = b.toBCD(10)
	} else {
		pkt.PhoneNumberBcd = b.toBCD(6)
	}
	pkt.MsgSerialNum = b.toWord()
	// 分包
	if pkt.MsgProperty.SubFlag > 0 {
		pkt.MsgTotal = b.toWord()
		pkt.MsgCur = b.toWord()
	}
	pkt.MsgBody = b.toBytes()
	return pkt
}

// jtRegister 注册信息
type jtRegister struct {
	Provincial     uint16
	City           uint16
	ManufacturerId []byte
	Model          []byte
	Id             []byte
	PlateColor     uint8
	CarLogo        string
}

func (b *jtBytes) toRegister() *jtRegister {
	var reg jtRegister
	reg.Provincial = b.toWord()
	reg.City = b.toWord()
	if b.jtVer == 1 {
		reg.ManufacturerId = b.toArrry(11)
		reg.Model = b.toArrry(30)
		reg.Id = b.toArrry(30)
	} else {
		reg.ManufacturerId = b.toArrry(5)
		reg.Model = b.toArrry(20)
		reg.Id = b.toArrry(7)
	}
	reg.PlateColor = b.toByte()
	reg.CarLogo = util.ToUTF8(b.toBytes())
	return &reg
}

type jtTerminalProperty struct {
	Type           uint16 `json:"type"`
	ManufacturerId []byte `json:"manufacturerId"` // 制造商Id
	Model          []byte `json:"model"`          // 型号
	Id             []byte `json:"id"`             // 终端Id
	ICCID          []byte `json:"iccid"`          // 终端 SIM 卡 ICCID
	hardVerLength  byte   `json:"-"`
	HardVer        string `json:"hardVer"`
	fwVerLength    byte   `json:"-"`
	FwVer          string `json:"fwVer"`
	GNSS           byte   `json:"gnss"`
	GPRS           byte   `json:"gprs"`
}

func (b *jtBytes) toTerminalProperty() *jtTerminalProperty {
	var info jtTerminalProperty
	info.Type = b.toWord()
	info.ManufacturerId = b.toArrry(5)
	if b.jtVer == 1 {
		info.Model = b.toArrry(30)
		info.Id = b.toArrry(30)
	} else {
		info.Model = b.toArrry(20)
		info.Id = b.toArrry(7)
	}
	info.ICCID = b.toArrry(10)
	info.hardVerLength = b.toByte()
	info.HardVer = util.ToUTF8(b.toArrry(int(info.hardVerLength)))
	info.fwVerLength = b.toByte()
	info.FwVer = util.ToUTF8(b.toArrry(int(info.fwVerLength)))
	info.GNSS = b.toByte()
	info.GPRS = b.toByte()
	return &info
}

type jtTerminalUpgrade struct {
	Type   uint8 `json:"type"`
	Result uint  `json:"result"`
}

func (b *jtBytes) toTerminalUpgrade() *jtTerminalUpgrade {
	var info jtTerminalUpgrade
	info.Type = b.toByte()
	info.Result = uint(b.toByte())
	return &info
}

// 8.21 临时位置跟踪控制
type Jt8202 struct {
	Interval      uint16 `json:"interval"`
	EffectiveTime uint32 `json:"effectiveTime"`
}

// 8.22 人工确认报警消息
type Jt8203 struct {
	AlarmSerial uint16 `json:"alarmSerial"`
	AlarmType   uint32 `json:"alarmType"`
}

// 8.23 文本信息下发
type Jt8300 struct {
	Flag uint8  `json:"flag"`
	Text string `json:"text"`
}

// 8.23 文本信息下发
type Jt8302 struct {
	Flag     uint8  `json:"flag"`
	Question string `json:"text"`
	Answers  []struct {
		Id      uint8  `json:"id"`
		Content string `json:"content"`
	} `json:"answers"`
}

// 8.28 信息点播菜单设置
type Jt8303 struct {
	Type    uint8 `json:"type"`
	Message []struct {
		Type uint8  `json:"type"`
		Name string `json:"name"`
	} `json:"message"`
}

// 8.30 信息服务
type Jt8304 struct {
	Type    uint8  `json:"type"`
	Content string `json:"content"`
}

// 8.31 电话回拨
type Jt8400 struct {
	Type        uint8  `json:"type"`
	PhoneNumber string `json:"phoneNumber"`
}

// 8.32 设置电话本
type Jt8401 struct {
	Type     uint8 `json:"type"`
	Contacts []struct {
		Type        uint8  `json:"type"`
		PhoneNumber string `json:"phoneNumber"`
		Name        string `json:"name"`
	} `json:"contacts"`
}

// 圆形区域
type Jt8600 struct {
	Action  uint8 `json:"action"`
	Regions []struct {
		Id          uint32  `json:"id"`
		Attr        uint16  `json:"attr"`
		La          float32 `json:"latitude"`
		Lo          float32 `json:"longitude"`
		Ra          uint32  `json:"radius"`
		St          string  `json:"startTime"`
		Et          string  `json:"endTime"`
		MaxSpd      uint16  `json:"maxSpeed"`
		KeepTime    uint8   `json:"keepTime"`
		MaxNightSpd uint16  `json:"maxNightSpeed"`
		Name        string  `json:"name"`
	} `json:"regions"`
}

// 矩形区域
type Jt8602 struct {
	Action  uint8 `json:"action"`
	Regions []struct {
		Id          uint32  `json:"id"`
		Attr        uint16  `json:"attr"`
		LeftTopLa   float32 `json:"leftTopLatitude"`
		LeftTopLo   float32 `json:"leftTopLongitude"`
		RightBotLa  float32 `json:"rightBotLatitude"`
		RightBotLo  float32 `json:"rightBotLongitude"`
		St          string  `json:"startTime"`
		Et          string  `json:"endTime"`
		MaxSpd      uint16  `json:"maxSpeed"`
		KeepTime    uint8   `json:"keepTime"`
		MaxNightSpd uint16  `json:"maxNightSpeed"`
		Name        string  `json:"name"`
	} `json:"regions"`
}

// 多边形区域
type Jt8604 struct {
	Id       uint32 `json:"id"`
	Attr     uint16 `json:"attr"`
	St       string `json:"startTime"`
	Et       string `json:"endTime"`
	MaxSpd   uint16 `json:"maxSpeed"`
	KeepTime uint8  `json:"keepTime"`
	Dots     []struct {
		La float32 `json:"latitude"`
		Lo float32 `json:"longitude"`
	} `json:"dots"`
	MaxNightSpd uint16 `json:"maxNightSpeed"`
	Name        string `json:"name"`
}

type Jt8601O3 struct {
	Ids []uint32 `json:"id"`
}
