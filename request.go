package xproto

// type
type ReqCode = int

const (
	Req_None ReqCode = iota //  空指令
	Req_Query
	Req_Parameters
	Req_LiveStream        //  实时视频
	Req_FtpTransfer       //  Ftp升级
	Req_Control           //  控制
	Req_SerialTransparent //  串口透传设置
	Req_SerialTransfer    //  串口传输数据
	Req_Playback          //  录像文件回放
	Req_FileTransfer      //  文件传输
	Req_Voice
	Req_WriteData //  裸数据（文件、语音...）
	Req_User      //  自定义请求
	Req_Close     //  关闭链接
)

// Response 通用请求应答
type Response struct {
	Session string      `json:"ss"`
	Details interface{} `json:"details"`
}

type FrameCode = uint16

const (
	Frame_Invalid    FrameCode = iota //  无效
	Frame_VideoI                      //  视频I帧
	Frame_VideoP                      //  视频P帧
	Frame_Audio                       //  hisi g726音频帧
	Frame_SerialPort                  //  串口数据帧
	Frame_File                        //  文件数据帧
	Frame_Status                      //  状态数据
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

type FileCode = int

const (
	File_Unknow      FileCode = iota
	File_NormalVideo          // 1:录像
	File_AlarmVideo           // 1:报警录像
	File_NormalPic            // 1:照片
	File_AlarmPic             // 1:报警照片
	File_Upgrade              // 1:升级文件
	File_Log                  // 1:日志文件
	File_Configure            // 1:配置文件
	File_BlackBox             // 1:黑匣子文件
	File_NormalVoice          // 1:录像
	File_AlarmVoice           // 1:录像
)

// Query 查找
type Query struct {
	Server      string   `json:"server"`
	Session     string   `json:"session"`
	StartTime   string   `json:"startTime"`
	EndTime     string   `json:"endTime"`
	FileType    FileCode `json:"fileType"`
	ChannelList string   `json:"channelList"` //
	StreamType  int      `json:"streamType"`  //
	StoreType   int      `json:"storeType"`   // 1 主录像 2
}

// File 查找结果
type File struct {
	FileName     string   `json:"fileName"`
	FileType     FileCode `json:"fileType"`
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

type ActCode = int

const (
	Act_Download ActCode = iota // 从设备下载文件
	Act_Upload                  // 上传文件到设备
)

// FtpTransfer 文件传输
type FtpTransfer struct {
	Session  string  `json:"session"`
	FtpURL   string  `json:"ftpUrl"`   // ftp://admin:admin@192.168.1.101:21
	FileType int     `json:"fileType"` // 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	FileSrc  string  `json:"fileSrc"`  // 源文件
	FileDst  string  `json:"fileDst"`  // 目标文件
	Action   ActCode `json:"action"`   //
}

// FileTransfer 文件传输
type FileTransfer struct {
	Session  string  `json:"session"`  // 自动生成session
	FileType int     `json:"fileType"` // 1:录像;2:报警录像;3:照片;4:报警照片;5:升级文件;6:日志文件;7:配置文件;8:黑匣子文件
	FileName string  `json:"fileName"`
	Server   string  `json:"server"`
	Action   ActCode `json:"action"`   //
	FileSize int     `json:"fileSize"` // 文件大小，上传文件到设备时要设置做校验用
	Offset   int     `json:"offset"`
}

type User struct {
	CodeId string                 `json:"codeId"`
	Data   map[string]interface{} `json:"data"`
}

type CtrlCode = int

const (
	Ctrl_Unknow        CtrlCode = iota
	Ctrl_Reboot                 // nil
	Ctrl_PTZ                    // PTZControl
	Ctrl_Reset                  // nil
	Ctrl_Vehi                   // nil
	Ctrl_GsensorAdjust          // nil
	Ctrl_OsdSpeed               // {"osdSpeed": }
	Ctrl_Snapshot               // CtrlSnapshot
	Ctrl_Format                 // nil
	Ctrl_Text                   // nil
	Ctrl_Property
	Ctrl_Playback
	Ctrl_Live
	Ctrl_MileReset
)

type Control struct {
	Type CtrlCode    `json:"ctrlCode"`
	Data interface{} `json:"data"`
}

type PTZCode = int

const (
	PTZ_Unknow PTZCode = iota
	PTZ_Up
	PTZ_Down
	PTZ_Left
	PTZ_Right
	PTZ_LeftUp
	PTZ_LeftDown
	PTZ_RightUp
	PTZ_RightDown
	PTZ_PreSetGoto
	PTZ_PreSetSet
	PTZ_PreSetClear
	PTZ_IrisOpen
	PTZ_IrisClose
	PTZ_ZoomIn
	PTZ_ZoomOut
	PTZ_FocusNear
	PTZ_FocusFar
	PTZ_AutoScan
	PTZ_BrushOn
	PTZ_BrushOff
	PTZ_CruiseOpen
	PTZ_CruiseClose
	PTZ_LightOn
	PTZ_LightOff
)

type CtrlPTZ struct {
	Code    PTZCode `json:"code"`
	Channel int     `json:"channel"`
	Speed   int     `json:"speed"`
}

type CtrlOsd struct {
	Speed string `json:"osdSpeed"`
}

type VehiCode = int

const (
	Vehi_CutOil      VehiCode = iota // 切断油路
	Vehi_RecoveryOil                 // 恢复油路
	Vehi_CutCircuit                  // 切断电路
	Vehi_RecCircuit                  // 恢复电路
	Vehi_OpenDoor
	Vehi_CloseDoor
)

type CtrlVehi struct {
	Code VehiCode `json:"code"`          
	Door int      `json:"door,omitempty"` // Action=5或6有效
}

type CtrlSnapshot struct {
	ChannelList string `json:"channelList"` // 1;2;3
	Reslution   int    `json:"reslution"`   // 0-跟随录像 1-1080p 2-720p 4-D1

}
type CaptureResult struct {
	Channel  int    `json:"channel"`
	FileName string `json:"fileName"`
}

type CtrlText struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type CtrlLive struct {
	Channel          uint8 `json:"channel"`
	Type             uint8 `json:"type"`
	CloseMeidaType   uint8 `json:"closeMeidaType"`
	ToggleStreamType uint8 `json:"toggleStreamType"`
}
