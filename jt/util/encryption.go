package util

import (
	"fmt"
	"strconv"
	"strings"
)

// Bcd2Number bcd解码
func Bcd2Number(bcd []byte) string {
	var number string
	for _, i := range bcd {
		number += fmt.Sprintf("%02X", i)
	}
	pos := strings.LastIndex(number, "F")
	if pos == 8 {
		return "0"
	}
	return number[pos+1:]
}

// Hex2Byte Hex转[]byte
func hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

// Number2bcd bcd解码
func Number2bcd(number string) []byte {
	var rNumber = number
	for i := 0; i < 8-len(number); i++ {
		rNumber = "f" + rNumber
	}
	bcd := hex2Byte(rNumber)
	return bcd
}

// CheckSum 计算
func CheckSum(b []byte) uint8 {
	ret := b[0]
	for i := 1; i < len(b); i++ {
		ret ^= b[i]
	}
	return ret
}