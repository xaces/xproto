package xproto

// ReqType 请求类型
type ReqType int

const (
	Req_None              = ReqType(0x00) //  空指令
	Req_Query             = ReqType(0x01) //  查询
	Req_Parameters        = ReqType(0x02) //  参数
	Req_LiveStream        = ReqType(0x03) //  实时视频
	Req_FtpTransfer       = ReqType(0x04) //  Ftp升级
	Req_Control           = ReqType(0x05) //  控制
	Req_SerialTransparent = ReqType(0x06) //  串口透传设置
	Req_SerialTransfer    = ReqType(0x07) //  串口传输数据
	Req_Playback          = ReqType(0x08) //  录像文件回放
	Req_FileTransfer      = ReqType(0x09) //  文件传输
	Req_Voice             = ReqType(0x0A)
	Req_WriteData         = ReqType(0x0B) //  裸数据（文件、语音...）
	Req_User              = ReqType(0x0C) //  自定义请求
	Req_Close             = ReqType(0xFF) //  关闭链接
)

// Response 通用请求应答
type Response struct {
	Session string      `json:"ss"`
	Details interface{} `json:"details"`
}

type FrameType int

const (
	Frame_Invalid    = 0 //  无效
	Frame_VideoI     = 1 //  视频I帧
	Frame_VideoP     = 2 //  视频P帧
	Frame_Audio      = 3 //  hisi g726音频帧
	Frame_SerialPort = 4 //  串口数据帧
	Frame_File       = 5 //  文件数据帧
	Frame_Status     = 6 //  状态数据
)

// LiveStream 实时流
type LiveStream struct {
	Session    string `json:"session"`
	Channel    int    `json:"channel"`
	StreamType int    `json:"streamType"` // 0-主码流 1-子码流
	On         int    `json:"on"`         // 0-停止 1-开启
	Frames     string `json:"frames"`     // 1;2;3
	Server     string `json:"server"`
}

// Voice 语音业务
type Voice struct {
	Session  string `json:"session"`
	WorkMode int    `json:"workMode"` // 0-监听 1-对讲 2-广播
	Channel  int    `json:"channel"`  // 通道
	On       int    `json:"on"`       // 0-停止 1-开启
	Server   string `json:"server"`
}

// Playback 录像文件回放
type Playback struct {
	Session   string `json:"session"`
	Name      string `json:"name"`
	Server    string `json:"server"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Channels  string `json:"channels"`
	Frames    string `json:"frames"` // 1;2;3
	Action    int    `json:"action"` // 0-下载式回放，1-流式回放
}

type FileType = int

const (
	// 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	File_Unknow      = FileType(0)
	File_NormalVideo = FileType(1)
	File_AlarmVideo  = FileType(2)
	File_NormalPic   = FileType(3)
	File_AlarmPic    = FileType(4)
	File_Upgrade     = FileType(5)
	File_Log         = FileType(6)
	File_Configure   = FileType(7)
	File_BlackBox    = FileType(8)
	File_NormalVoice = FileType(9)
	File_AlarmVoice  = FileType(10)
)

// Query 查找
type Query struct {
	Server      string   `json:"server"`
	Session     string   `json:"session"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	FileType    FileType `json:"fileType"`    // 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	ChannelList string   `json:"channelList"` //
	StreamType  int      `json:"streamType"`  //
	StoreType   int      `json:"storeType"`   // 1 主录像 2
}

// File 查找结果
type File struct {
	FileName     string   `json:"fileName"`
	FileType     FileType `json:"fileType"`
	Channel      int      `json:"channel"`
	StartTime    string   `json:"startTime"`
	EndTime      string   `json:"endTime"`
	FileSize     int      `json:"fileSize"`
	FileDuration int      `json:"fileDuration"`
}

// SerialTransparent 串口透传设置
type SerialTransparent struct {
	Session  string `json:"session"`
	Port     int    `json:"port"`
	BaudRate int    `json:"baudRate"`
	DatBit   int    `json:"dataBit"`
	CheckBit int    `json:"checkBit"`
	StopBit  int    `json:"stopBit"`
}

// SerialData 透传数据
type SerialData struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

// AcitonCode 动作代码
type AcitonCode int

const (
	Act_Download AcitonCode = iota // 从设备下载文件
	Act_Upload                     // 上传文件到设备
)

// FtpTransfer 文件传输
type FtpTransfer struct {
	Session  string     `json:"session"`
	FtpURL   string     `json:"ftpUrl"`   // ftp://admin:admin@192.168.1.101:21
	FileType FileType   `json:"fileType"` // 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	FileSrc  string     `json:"fileSrc"`  // 源文件
	FileDst  string     `json:"fileDst"`  // 目标文件
	Action   AcitonCode `json:"action"`   //
}

// FileTransfer 文件传输
type FileTransfer struct {
	Session  string     `json:"session"`  // 自动生成session
	FileType FileType   `json:"fileType"` // 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	FileName string     `json:"fileName"`
	Server   string     `json:"server"`
	Action   AcitonCode `json:"action"`   //
	FileSize int        `json:"fileSize"` // 文件大小，上传文件到设备时要设置做校验用
	Offset   int        `json:"offset"`
}

type PtzCode int //

const (
	Ctrl_PtzUnknow      = PtzCode(0x00)
	Ctrl_PtzUp          = PtzCode(0x01)
	Ctrl_PtzDown        = PtzCode(0x02)
	Ctrl_PtzLeft        = PtzCode(0x03)
	Ctrl_PtzRight       = PtzCode(0x04)
	Ctrl_PtzLeftUp      = PtzCode(0x05)
	Ctrl_PtzLeftDown    = PtzCode(0x06)
	Ctrl_PtzRightUp     = PtzCode(0x07)
	Ctrl_PtzRightDown   = PtzCode(0x08)
	Ctrl_PtzPreSetGoto  = PtzCode(0x09)
	Ctrl_PtzPreSetSet   = PtzCode(0x0A)
	Ctrl_PtzPreSetClear = PtzCode(0x0B)
	Ctrl_PtzIrisOpen    = PtzCode(0x0C)
	Ctrl_PtzIrisClose   = PtzCode(0x0D)
	Ctrl_PtzZoomIn      = PtzCode(0x0E)
	Ctrl_PtzZoomOut     = PtzCode(0x0F)
	Ctrl_PtzFocusNear   = PtzCode(0x10)
	Ctrl_PtzFocusFar    = PtzCode(0x11)
	Ctrl_PtzAutoScan    = PtzCode(0x12)
	Ctrl_PtzBrushOn     = PtzCode(0x13)
	Ctrl_PtzBrushOff    = PtzCode(0x14)
	Ctrl_PtzCruiseOpen  = PtzCode(0x15)
	Ctrl_PtzCruiseClose = PtzCode(0x16)
	Ctrl_PtzLightOn     = PtzCode(0x17)
	Ctrl_PtzLightOff    = PtzCode(0x18)
)

type PtzControl struct {
	Code    PtzCode `json:"code"`
	Channel int     `json:"channel"`
	Speed   int     `json:"speed"`
}

type OsdControl struct {
	OsdSpeed string `json:"osdSpeed"` // obd/
}

type VehiCode int

const (
	Ctrl_VehiCutOil          = VehiCode(0x01) // 切断油路
	Ctrl_VehiRecoveryOil     = VehiCode(0x02) // 恢复油路
	Ctrl_VehiCutCircuit      = VehiCode(0x03) // 切断电路
	Ctrl_VehiRecoveryCircuit = VehiCode(0x04) // 恢复电路
	Ctrl_VehiOpenDoor        = VehiCode(0x05)
	Ctrl_VehiCloseDoor       = VehiCode(0x06)
)

type VehiControl struct {
	Code VehiCode `json:"code"`           // obd/
	Door int      `json:"door,omitempty"` // Action=5或6有效
}

type CaptureControl struct {
	ChannelList string `json:"channelList"` // 1;2;3
	Reslution   int    `json:"reslution"`   // 0-跟随录像 1-1080p 2-720p 4-D1
}

type CaptureResult struct {
	Channel  int    `json:"channel"`
	FileName string `json:"fileName"`
}

type TextControl struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type LiveControl struct {
	Channel          uint8 `json:"channel"`
	Type             uint8 `json:"type"`
	CloseMeidaType   uint8 `json:"closeMeidaType"`
	ToggleStreamType uint8 `json:"toggleStreamType"`
}

type CtrlCode int

const (
	Ctrl_Reboot        = CtrlCode(0x01) // nil
	Ctrl_PTZ           = CtrlCode(0x02) // PtzControl
	Ctrl_Reset         = CtrlCode(0x03) // nil
	Ctrl_Vehi          = CtrlCode(0x04) // nil
	Ctrl_GsensorAdjust = CtrlCode(0x05) // nil
	Ctrl_OsdSpeed      = CtrlCode(0x06) // OsdControl
	Ctrl_Capture       = CtrlCode(0x07) // nil
	Ctrl_Format        = CtrlCode(0x08) // nil
	Ctrl_Text          = CtrlCode(0x09) // nil
	Ctrl_Property      = CtrlCode(0x0A)
	Ctrl_Playback      = CtrlCode(0x0B)
	Ctrl_Live          = CtrlCode(0x0C)
)

type Control struct {
	Type CtrlCode    `json:"ctrlType"`
	Data interface{} `json:"data"`
}

type User struct {
	CodeId string                 `json:"codeId"`
	Data   map[string]interface{} `json:"data"`
}
