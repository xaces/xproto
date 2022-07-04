package jt

import (
	"github.com/xaces/xproto"
	"github.com/xaces/xproto/jt/util"
)

type jtBulkStatus struct {
	DataNumber uint16
	BlindSpot  uint8
	Length     uint16
	St         jtStatus
}

type jtStatus struct {
	AlarmFlag uint32
	Status    uint32
	Latitude  uint32 // *1000000
	Longitude uint32
	Altitude  uint16
	Speed     uint16
	Direction uint16
	DTU       [6]byte // 6个字节的BCD

}

func dtu(dtu []byte) string {
	tm := util.Bcd2Number(dtu)
	tm = "20" + tm[:2] + "-" + tm[2:4] + "-" + tm[4:6] + " " + tm[6:8] + ":" + tm[8:10] + ":" + tm[10:]
	return tm
}

func jtgps(val uint32) float32 {
	return float32(val) / 1000000
}

// 附加信息 ID 附加信息长度 描述及要求
// 0x01 4 里程，DWORD，1/10km，对应车上里程表读数
// 0x02 2 油量，WORD，1/10L，对应车上油量表读数
// 0x03 2 行驶记录功能获取的速度，WORD，1/10km/h
// 0x04 2 需要人工确认报警事件的 ID，WORD，从 1 开始计数
// 0x05-0x10 保留
// 0x11 1 或 5 超速报警附加信息见表 28
// 0x12 6 进出区域/路线报警附加信息见表 29
// 0x13 7 路段行驶时间不足/过长报警附加信息见表 30
// 0x14-0x24 保留
// 0x25 4 扩展车辆信号状态位，定义见表 31
// 0x2A 2 IO 状态位，定义见表 32
// 0x2B 4 模拟量，bit0-15，AD0；bit16-31，AD1。
// 0x30 1 BYTE，无线通信网络信号强度
// 0x31 1 BYTE，GNSS 定位卫星数
// 0xE0 后续信息长度 后续自定义信息长度
// 0xE1-0xFF 自定义区域

type jtStatusAttach struct {
	AttachID       uint8
	AttachLenght   uint8
	Mileage        uint32   // attachID 0x01 1/10km
	OilVolume      uint16   // attachID 0x02 1/10L
	ReSpeed        uint16   // attachID 0x03 1/10km/h
	ManualID       uint16   // attachID 0x04
	OverSpeedAlarm struct { // attachID 0x11 表28
		localType uint8 //
		AreaID    uint32
	}
	EntryAndExitAlarm struct { // attachID 0x12 表29
		localType uint8 //
		AreaID    uint32
		Direction uint8
	}
	TravelTimeAlarm struct { // attachID 0x13 表30
		AreaID     uint32
		TravelTime uint16 //
		Result     uint8
	}
	SignalState      uint32 // attachID 0x25 表31
	IoState          uint16 // attachID 0x2A 表32
	Analog           uint32 //attachID 0x2B  模拟量，bit0-15，AD0；bit16-31，AD1。
	WirelessStrength uint8  // attachID 0x30 无线通信网络信号强度
	GNSSBNumber      uint8  // attachID 0x31 GNSS 定位卫星数
	// FollowLength
}

func (b *jtBytes) toAttachStatus() *jtStatusAttach {
	var a jtStatusAttach
	a.AttachID = b.toByte()
	a.AttachLenght = b.toByte()
	switch a.AttachID {
	case 0x01:
		a.Mileage = b.toDWord() * 10
	case 0x02:
		a.OilVolume = b.toWord() * 10
	case 0x03:
		a.ReSpeed = b.toWord() * 10
	case 0x04:
		a.ManualID = b.toWord()
	case 0x11:
		a.OverSpeedAlarm.localType = b.toByte()
		if a.OverSpeedAlarm.localType > 0 {
			a.OverSpeedAlarm.AreaID = b.toDWord()
		}
	case 0x12:
		a.EntryAndExitAlarm.localType = b.toByte()
		a.EntryAndExitAlarm.AreaID = b.toDWord()
		a.EntryAndExitAlarm.Direction = b.toByte()
	case 0x13:
		a.TravelTimeAlarm.AreaID = b.toDWord()
		a.TravelTimeAlarm.TravelTime = b.toWord()
		a.TravelTimeAlarm.Result = b.toByte()
	case 0x25:
		a.SignalState = b.toDWord()
	case 0x2A:
		a.IoState = b.toWord()
	case 0x2B:
		a.Analog = b.toDWord()
	case 0x30:
		a.WirelessStrength = b.toByte()
	case 0x31:
		a.GNSSBNumber = b.toByte()
	}
	// log.Printf("attachID %d mileage %d\n", a.AttachID, a.Mileage)
	return &a
}

// status 表
// 位 状态
// 0 0：ACC 关；1： ACC 开
// 1 0：未定位；1：定位
// 2 0：北纬；1：南纬
// 3 0：东经；1：西经
// 4 0：运营状态；1：停运状态
// 5 0：经纬度未经保密插件加密；1：经纬度已经保密插件加密
// 6-7 保留
// 8-9 00：空车；01：半载；10：保留；11：满载 （可用于客车的空、重车及货车的空载、满载状态表示，人工输入或传感器获取）
// 10 0：车辆油路正常；1：车辆油路断开
// 11 0：车辆电路正常；1：车辆电路断开
// 12 0：车门解锁；1：车门加锁
// 13 0：门 1 关；1：门 1 开（前门）
// 14 0：门 2 关；1：门 2 开（中门）
// 15 0：门 3 关；1：门 3 开（后门）
// 16 0：门 4 关；1：门 4 开（驾驶席门）
// 17 0：门 5 关；1：门 5 开（自定义）
// 18 0：未使用 GPS 卫星进行定位；1：使用 GPS 卫星进行定位
// 19 0：未使用北斗卫星进行定位；1：使用北斗卫星进行定位
// 20 0：未使用 GLONASS 卫星进行定位；1：使用 GLONASS 卫星进行定位
// 21 0：未使用 Galileo 卫星进行定位；1：使用 Galileo 卫星进行定位
// 22-31 保留

// doStatusParse 解析位置信息 28byte
func (p *Client) doStatusParse(b *jtBytes) (st xproto.Status, alrflag uint32) {
	st.DeviceNo = p.DeviceNo
	alrflag = b.toDWord()
	status := b.toDWord()
	st.Acc = uint8(status & 0x01)
	st.Location.Type = uint8(status & 0x02)
	st.Location.Latitude = jtgps(b.toDWord())
	st.Location.Longitude = jtgps(b.toDWord())
	st.Location.Altitude = int(b.toWord())
	st.Location.Speed = float32(b.toWord()) / 100.0
	st.Location.Angle = float32(b.toWord())
	st.DTU = dtu(b.toBCD(6))
	if alrflag > 0 {
		return
	}
	b.toAttachStatus()
	return
}
