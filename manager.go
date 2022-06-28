package xproto

import (
	"sync"
)

// GUID 生成管理GUID
func generateGUID(args []string) string {
	if len(args) > 1 && args[1] != "" {
		return args[0] + "." + args[1]
	}
	return args[0]
}

// gMnger 默认列表
var gMnger sync.Map

// SyncAdd 添加
func SyncAdd(client IClient, args ...string) error {
	guid := generateGUID(args)
	if _, loaded := gMnger.LoadOrStore(guid, client); loaded {
		return ErrObjectExist
	}
	return nil
}

// SyncRemove 从映射表中删除链接
func SyncRemove(args ...string) {
	guid := generateGUID(args)
	gMnger.Delete(guid)
}

// SyncGet 从映射表中获取链接
func SyncGet(args ...string) IClient {
	guid := generateGUID(args)
	if v, ok := gMnger.Load(guid); ok {
		ctx, _ := v.(IClient)
		return ctx
	}
	return nil
}

// SyncStop 停止服务
func SyncStop(args ...string) {
	if ctx := SyncGet(args...); ctx != nil {
		ctx.Request(Req_Close, nil, nil)
	}
}

// SyncSend 发送信息到设备
func SyncSend(cmd ReqType, s interface{}, r interface{}, args ...string) error {
	if ctx := SyncGet(args...); ctx != nil {
		return ctx.Request(cmd, s, r)
	}
	return ErrObjectNoExist
}

// SyncWrite 发送信息到设备
func SyncWrite(s *RawData) error {
	return SyncSend(Req_WriteData, s, nil, s.DeviceNo, s.Session)
}

// SyncStopAll 停止所有链接
func SyncStopAll() {
	gMnger.Range(func(key, value interface{}) bool {
		ctx, _ := value.(IClient)
		ctx.Request(Req_Close, nil, nil)
		return true
	})
}
