package xproto

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xaces/xproto/internal"
)

// prefix 文件前缀
func prefix() string {
	return time.Now().Format("20060102150405")
}

// FileSess 文件传输会话号
func FileSess(act AcitonCode, file string) string {
	ss := prefix() + "." + file
	if act == Act_Upload {
		ss = "Up" + ss
	}
	return ss
}

func FileOfSess(ss string) (string, AcitonCode) {
	if strings.HasPrefix(ss, "Up") {
		return ss[17:], Act_Upload
	} else {
		return ss[15:], Act_Download
	}
}

// UUID 生成Guid字串
func UUID() string {
	u := uuid.NewV4()
	return u.String()
}

func NewUUID(s *string) {
	if *s == "" {
		*s = UUID()
	}
}

func alarmUUID(a *Alarm) {
	if a.UUID != "" {
		return
	}
	var tt int64
	if a.StartTime == "" {
		tt = time.Now().Unix()
	} else {
		tt = internal.ToUnixTime(a.StartTime)
	}
	zerol := 12 - len(a.DeviceNo)
	a.UUID = fmt.Sprintf("%10d%03d%03d%0*d%s", tt, a.Type, a.SubType, zerol, 0, a.DeviceNo)
}
