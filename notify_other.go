package xproto

type Frame struct {
	DeviceNo  string
	Session   string
	Channel   uint16
	Type      uint16
	Timestamp int64
	HeaderLen int
	Data      []byte // 包含H数据
	*Status
}

type RawData struct {
	DeviceNo  string
	Session   string
	Type      uint16
	Channel   uint16 //
	Timestamp int64  // 时间戳
	Data      []byte
}
