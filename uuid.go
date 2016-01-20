package ble

import "strings"

const hexDigits = "0123456789abcdef"

func isHex(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune(hexDigits, c) {
			return false
		}
	}
	return true
}

func validUUID(u string) bool {
	switch len(u) {
	case 4:
		return isHex(u)
	case 36:
		return isHex(u[0:8]) && u[8] == '-' &&
			isHex(u[9:13]) && u[13] == '-' &&
			isHex(u[14:18]) && u[18] == '-' &&
			isHex(u[19:23]) && u[23] == '-' &&
			isHex(u[24:36])
	default:
		return false
	}
}
