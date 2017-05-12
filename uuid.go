package ble

import "strings"

func hexMatch(s, pattern string) bool {
	const hexDigits = "0123456789abcdef"
	if len(s) != len(pattern) {
		return false
	}
	for i := range s {
		switch pattern[i] {
		case 'x':
			if strings.IndexByte(hexDigits, s[i]) == -1 {
				return false
			}
		default:
			if s[i] != pattern[i] {
				return false
			}
		}
	}
	return true
}

// ValidUUID checks whether a string is a valid UUID.
func ValidUUID(u string) bool {
	switch len(u) {
	case 4:
		return hexMatch(u, "xxxx")
	case 36:
		return hexMatch(u, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	default:
		return false
	}
}
