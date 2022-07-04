package jt

import (
	"errors"
)

const (
	kResSuc = iota
	kResFailed
	kResError
	kResUnsupport
	kResAlarmAck
)

func codeError(code uint8) error {
	var err error = nil
	switch code {
	case kResFailed:
		err = errors.New("failed")
	case kResError:
		err = errors.New("message error")
	case kResUnsupport:
		err = errors.New("unsupport")
	case kResAlarmAck:
		err = errors.New("alarm ack")
	}
	return err
}
