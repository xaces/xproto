package xproto

// NetCode 网络类型
type NetCode = int

const (
	Net_Unknow NetCode = iota //  未知
	Net_Cable                 //  有线
	Net_Wifi                  //  Wifi
	Net_2G                    //  2G
	Net_3G                    //  3G
	Net_4G                    //  4G
	Net_5G                    //  5G
)

// LinkCode 链路类型
type LinkCode = int

const (
	Link_Unknow             LinkCode = iota //  未知
	Link_Signal                             //  信令
	Link_Interactive                        //  交互链路（文件查询等）
	Link_Query                              //  查询链路
	Link_LiveStream                         //  实时预览
	Link_Playback                           //  回放
	Link_FileTransfer                       //  文件传输
	Link_SerialTransmission                 //  串口透传
	Link_Voice                              //  语音
)

// Access 链路注册
type Access struct {
	DeviceID      uint     `json:"deviceId"`
	DeviceNo      string   `json:"deviceNo"`   // 设备号
	DeviceTime    string   `json:"deviceTime"` // 设备时间
	Online        bool     `json:"online"`
	RemoteAddress string   `json:"remoteAddress"` // 设备网络地址
	Session       string   `json:"session"`       // 链路会话号
	LinkType      LinkCode `json:"linkType"`      // 链路类型
	UpTraffic     int64    `json:"upTraffic"`     // 上行流量
	DownTraffic   int64    `json:"downTraffic"`   // 下行流量
	NetType       NetCode  `json:"netType"`       // 网络类型
	Version       string   `json:"version"`       // 版本信息
	DevType       string   `json:"devType"`       // 设备类型
}
