package util

import (
	"strconv"

	"github.com/axgle/mahonia"
)

var (
	dec = mahonia.NewDecoder("gbk")
	enc = mahonia.NewEncoder("gbk")
)

func ToUTF8(data []byte) string {
	return dec.ConvertString(string(data))
}
func ToGBK(str string) []byte {
	return []byte(enc.ConvertString(str))
}

func HexStringToUint16(s string) uint16 {
	v, _ := strconv.ParseUint(s[2:], 16, 32)
	return uint16(v)
}
