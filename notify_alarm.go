package xproto

// AlarmType 报警类型
type AlarmType = int

const (
	Alarm_Unknow                  AlarmType = 0x00 //未知   //unknow
	Alarm_VideoLost               AlarmType = 0x01 //视频丢失   //video lost
	Alarm_MotionDetection         AlarmType = 0x02 //移动侦测   //motion detection
	Alarm_VideoMask               AlarmType = 0x03 //视频遮挡   //video blind
	Alarm_InputGenerate           AlarmType = 0x04 //输入触发   //input trigger
	Alarm_Urgency                 AlarmType = 0x05 //紧急告警   //emmergency alarm
	Alarm_LowSpeed                AlarmType = 0x06 //低速告警   //low speed alarm
	Alarm_OverSpeed               AlarmType = 0x07 //超速告警   //overspeed alarm
	Alarm_LowTemperature          AlarmType = 0x08 //低温告警   //low temperature alarm
	Alarm_OverTemperature         AlarmType = 0x09 //高温告警   //high temperature alarm
	Alarm_Humidity                AlarmType = 0x0a //湿度告警   //humidity alarm
	Alarm_ParkTimeout             AlarmType = 0x0b //超时停车   //park timeout alarm
	Alarm_Shake                   AlarmType = 0x0c //震动告警   //vibration alarm
	Alarm_Geofence                AlarmType = 0x0d //电子围栏   //electronic defence
	Alarm_GeoLine                 AlarmType = 0x0e //电子线路   //electronic line
	Alarm_DoorException           AlarmType = 0x0f //异常开关门   //door open/ close exception
	Alarm_StorageException        AlarmType = 0x10 //存储介质异常   //storage exception
	Alarm_FatigueDriving          AlarmType = 0x11 //疲劳驾驶   //fatigue driving
	Alarm_FuelException           AlarmType = 0x12 //油量异常   //fuel exception
	Alarm_IllegalFire             AlarmType = 0x13 //非法点火   //illegal ACC
	Alarm_LocationModuleException AlarmType = 0x14 //定位模块异常   //location module exception
	Alarm_FrontPanelPrise         AlarmType = 0x15 //前面板打开   //front panel prise
	Alarm_SwipeCard               AlarmType = 0x16 //刷卡		//swipe card
	Alarm_SwipeIbutton            AlarmType = 0x17 //IBUTTON	//ibutton
	Alarm_HarshAcceleration       AlarmType = 0x18 //急加速
	Alarm_HarshBraking            AlarmType = 0x19 //急减速
	Alarm_LowSpeedWarn            AlarmType = 0x1a //低速预警	26
	Alarm_OverSpeedWarn           AlarmType = 0x1b //高速预警
	Alarm_VoltageWarn             AlarmType = 0x1c //电压告警
	Alarm_PeopleCounting          AlarmType = 0x1d //人数统计
	Alarm_DmsAndadasAlm           AlarmType = 0x1e //dms和adas报警（主动安全）
	Alarm_AccOn                   AlarmType = 0x1f //acc on
)

// Channel 和通道相关的报警数据如: io/video...等，具体数据自定义
// 和通道相关的报警
type Channel struct {
	Channel int `json:"ch"`
}

// Threshold 和阈值相关的报警
type Threshold struct {
	Value int     `json:"vt"` // 触发
	Time  int     `json:"tt"`
	Min   float32 `json:"min"`
	Max   float32 `json:"max"`
	Avg   float32 `json:"avg"`
	Cur   float32 `json:"cur"`
	Pre   float32 `json:"pre"`
}

// Shake
type Shake struct {
	Type int `json:"dt"` // 1-x 2-y 3-z 4-hit 5-tilt
	Threshold
}

type Speed struct {
	Threshold
	Dur  uint `json:"dur"`  // 时长
	Spds int  `json:"spds"` // 速度来源
}

// Alarm 报警
type Alarm struct {
	Tag       string      `json:"-"`
	DeviceNo  string      `json:"deviceNo"`
	DTU       string      `json:"dtu"`
	UUID      string      `json:"uuid"`
	Type      AlarmType   `json:"type"`
	SubType   int         `json:"subType"`
	StartTime string      `json:"startTime"`
	EndTime   string      `json:"endTime"`
	Data      interface{} `json:"data"`
	Status    *Status     `json:"status"`
}
