package xproto

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xaces/xproto/internal"
)

// LogStatus 状态数据输出
func LogStatus(tag string, arg interface{}, s *Status) {
	fmt.Printf("%s >> | %s | %-16s | DTU %s Latitude %0.6f Longitude %0.6f speed %0.6f\n", tag, internal.TimeNow(), s.DeviceNo, s.DTU,
		s.Location.Latitude, s.Location.Longitude, s.Location.Speed)
}

// LogAlarm 报警输出
func LogAlarm(data []byte, arg interface{}, a *Alarm) {
	fmt.Printf("%s >> | %s | %-16s | %s\n", a.Tag, internal.TimeNow(), a.DeviceNo, string(data))
	if a.Status != nil {
		LogStatus(a.Tag, arg, a.Status)
	}
}

// LogEvent 事件输出
func LogEvent(data []byte, arg interface{}, e *Event) {
	fmt.Printf("%s >> | %s | %-16s | %s\n", e.Tag, internal.TimeNow(), e.DeviceNo, string(data))
	if e.Status != nil {
		LogStatus(e.Tag, arg, e.Status)
	}
}

// LogFrame 帧输出
func LogFrame(arg interface{}, f *Frame) error {
	switch f.Type {
	case Frame_File: //ss 是文件名
		if v, ok := arg.(*os.File); ok {
			v.Write(f.Data)
		}
		// (*os.File)(unsafe.Pointer(&arg)).Write(r.Data)
	case Frame_SerialPort:
		fmt.Printf("%s[%s] >> | %v\n", f.Session, f.DeviceNo, f.Data)
	default:
		fmt.Printf("%s >> | channel %d type: %d timestamp: %v length %d\n", f.Session, f.Channel, f.Type, f.Timestamp, len(f.Data))
	}
	return nil
}

func DownloadFile(filename string, v interface{}) (interface{}, error) {
	if filename != "" {
		os.MkdirAll(filepath.Dir(filename), os.ModePerm)
		return os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	}
	if v != nil {
		if file, ok := (v).(*os.File); ok {
			file.Close()
		}
	}
	return nil, nil
}

// UploadFile 上传文件参考
func UploadFile(a *Access, filename string, showProgress bool) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	// 启动线程读文件传输
	go func(fp *os.File, a *Access) {
		fi, _ := fp.Stat()
		defer fp.Close()
		bs := make([]byte, 10*1024)
		sendTotalSize := 0
		lastBar := 0
		for {
			s, err := fp.Read(bs)
			raw := &RawData{Type: Frame_File, DeviceNo: a.DeviceNo, Session: a.Session}
			if err == io.EOF {
				SyncWrite(raw) // 发送结束
				break
			}
			raw.Data = bs[:s]
			if err = SyncWrite(raw); err != nil {
				break
			}
			time.Sleep(time.Millisecond * 10)
			if !showProgress {
				continue
			}
			sendTotalSize += s
			bar := (sendTotalSize * 100) / int(fi.Size())
			if bar > lastBar {
				fmt.Printf("\r[%-100s] [%3d%%/%d]\n", strings.Repeat("#", bar), bar, fi.Size())
				lastBar = bar
			}
		}
	}(fp, a)
	return nil
}
