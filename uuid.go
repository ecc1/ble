package ble

import (
	"log"
	"strings"
)

const (
	// BluetoothBaseUUID for service discovery.
	// See www.bluetooth.com/specifications/assigned-numbers/service-discovery
	BluetoothBaseUUID = "0000xxxx-0000-1000-8000-00805f9b34fb"

	hexDigits = "0123456789abcdef"
)

func hexMatch(s, pattern string) bool {
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
		// 16-bit UUID
		return hexMatch(u, "xxxx")
	case 8:
		// 32-bit UUID
		return hexMatch(u, "xxxxxxxx")
	case 36:
		// 128-bit UUID
		return hexMatch(u, "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
	default:
		return false
	}
}

// LongUUID returns the 128-bit UUID corresponding to a possibly-shorter UUID,
// which must be valid.
func LongUUID(u string) string {
	switch len(u) {
	case 4:
		return "0000" + u + BluetoothBaseUUID[8:]
	case 8:
		return u + BluetoothBaseUUID[8:]
	case 36:
		return u
	default:
		log.Panicf("invalid UUID %q has length %d", u, len(u))
	}
	panic("unreachable")
}

// ShortUUID returns the shortest UUID corresponding to the given UUID,
// which must be valid.
func ShortUUID(u string) string {
	switch len(u) {
	case 4:
		return u
	case 8:
		if u[0:4] == "0000" {
			return u[4:8]
		}
		return u
	case 36:
		if u[8:36] != BluetoothBaseUUID[8:36] {
			return u
		}
		if u[0:4] == "0000" {
			return u[4:8]
		}
		return u[0:8]
	default:
		log.Panicf("invalid UUID %q has length %d", u, len(u))
	}
	panic("unreachable")
}

// UUIDs represents a list of UUIDs.
type UUIDs []string

// The String method allows a list of UUIDs to be printed in short form.
func (uuids UUIDs) String() string {
	var b strings.Builder
	// nolint
	b.WriteByte('[')
	for i, u := range uuids {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(ShortUUID(u))
	}
	b.WriteByte(']')
	return b.String()
}
