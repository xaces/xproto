package xproto

// EventType 事件类型
type EventType = int

const (
	Event_File             EventType = iota + 0x500 //  未知
	Event_FileTimedCapture                          //  文件
	Event_FileLittle                                //  ftp文件传输
	Event_FileTransfer
	Event_Upgrade
	Event_FtpTransfer = (0x5ff)
)

// EventFile
type EventFile struct {
	FileName string   `josn:"fileName"`
	FileType FileType `json:"fileType"`
}

type FileCapture struct {
	EventFile
	Channel   int     `json:"channel"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Speed     float32 `json:"speed"`
}

type FileLittle struct {
	EventFile
	Channel  int `json:"channel"`
	Size     int `json:"size"`
	Duration int `json:"duration"`
}

type FileFtp struct {
	EventFile
	Ftp string `json:"ftp"`
}

// Event 事件
type Event struct {
	Tag      string      `json:"-"`
	DeviceID uint        `json:"deviceId"`
	DeviceNo string      `json:"deviceNo"`
	Session  string      `json:"session"`
	DTU      string      `json:"dtu"`
	Type     EventType   `json:"type"`
	Data     interface{} `json:"data"` // 负载数据
	Status   *Status     `json:"status"`
}
