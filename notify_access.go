package xproto

// NetType 网络类型
type NetType int

const (
	Net_Unknow NetType = 0x00 //  未知
	Net_Cable  NetType = 0x01 //  有线
	Net_Wifi   NetType = 0x02 //  Wifi
	Net_2G     NetType = 0x03 //  2G
	Net_3G     NetType = 0x04 //  3G
	Net_4G     NetType = 0x05 //  4G
	Net_5G     NetType = 0x06 //  5G
)

// LinkType 链路类型
type LinkType int

const (
	Link_Unknow             = LinkType(0x00) //  未知
	Link_Signal             = LinkType(0x01) //  信令
	Link_Interactive        = LinkType(0x02) //  交互链路（文件查询等）
	Link_Query              = LinkType(0x03) //  查询链路
	Link_LiveStream         = LinkType(0x04) //  实时预览
	Link_Playback           = LinkType(0x05) //  回放
	Link_FileTransfer       = LinkType(0x06) //  文件传输
	Link_SerialTransmission = LinkType(0x07) //  串口透传
	Link_Voice              = LinkType(0x08) //  语音
)

// Access 链路注册
type Access struct {
	DeviceID      uint     `json:"deviceId"`
	DeviceNo      string   `json:"deviceNo"`   // 设备号
	DeviceTime    string   `json:"deviceTime"` // 设备时间
	Online        bool     `json:"online"`
	RemoteAddress string   `json:"remoteAddress"` // 设备网络地址
	Session       string   `json:"session"`       // 链路会话号
	LinkType      LinkType `json:"linkType"`      // 链路类型
	UpTraffic     int64    `json:"upTraffic"`     // 上行流量
	DownTraffic   int64    `json:"downTraffic"`   // 下行流量
	NetType       NetType  `json:"netType"`       // 网络类型
	Version       string   `json:"version"`       // 版本信息
	DevType       string   `json:"devType"`       // 设备类型
}
