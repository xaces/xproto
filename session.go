package xproto

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/xaces/xproto/internal"
)

// up-1-filename
// up-文件类型-文件名
// FileSess 文件传输会话号
func FileSess(r *FileTransfer) string {
	ss := strconv.Itoa(r.FileType) + "@" + r.FileName
	if r.Action == Act_Upload {
		ss = "file://" + ss
		r.FileName = filepath.Base(r.FileName)
	}
	return ss
}

func FileOfSess(ss string) (r FileTransfer) {
	r.Action = Act_Download
	if strings.HasPrefix(ss, "file://") {
		r.Action = Act_Upload
		ss = ss[7:]
	}
	arrs := strings.Split(ss, "@")
	r.FileType, _ = strconv.Atoi(arrs[0])
	r.FileName = arrs[1]
	if r.Action == Act_Download {
		r.FileName = filepath.Base(r.FileName)
	}
	return
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
