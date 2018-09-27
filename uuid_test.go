package ble

import (
	"testing"
)

func TestValidUUID(t *testing.T) {
	cases := []struct {
		uuid  string
		valid bool
	}{
		{"1234", true},
		{"abcd", true},
		{"12345678", true},
		{"00001234-0000-1000-8000-00805f9b34fb", true},
		{"12345", false},
		{"ABCD", false},
		{"abcx", false},
		{"123456789", false},
		{"1234567z", false},
		{"00001234-0000-1000-8000-00805F9B34FB", false},
		{"0000123400001000800000805f9b34fb", false},
		{"g0001234-0000-1000-8000-00805f9b34fb", false},
	}
	for _, c := range cases {
		t.Run(c.uuid, func(t *testing.T) {
			valid := ValidUUID(c.uuid)
			if valid != c.valid {
				t.Errorf("ValidUUID(%s) == %v, want %v", c.uuid, valid, c.valid)
			}
		})
	}
}

func TestLongUUID(t *testing.T) {
	cases := []struct {
		uuid     string
		longUUID string
	}{
		{"1234", "00001234-0000-1000-8000-00805f9b34fb"},
		{"00001234", "00001234-0000-1000-8000-00805f9b34fb"},
		{"00001234-0000-1000-8000-00805f9b34fb", "00001234-0000-1000-8000-00805f9b34fb"},
	}
	for _, c := range cases {
		t.Run(c.uuid, func(t *testing.T) {
			longUUID := LongUUID(c.uuid)
			if longUUID != c.longUUID {
				t.Errorf("LongUUID(%s) == %v, want %v", c.uuid, longUUID, c.longUUID)
			}
		})
	}
}

func TestShortUUID(t *testing.T) {
	cases := []struct {
		uuid      string
		shortUUID string
	}{
		{"00001234-0000-1000-8000-00805f9b34fb", "1234"},
		{"00001234", "1234"},
		{"1234", "1234"},
		{"12345678-0000-1000-8000-00805f9b34fb", "12345678"},
		{"12345678-9abc-1000-8000-00805f9b34fb", "12345678-9abc-1000-8000-00805f9b34fb"},
		{"12345678-9abc-def0-8000-00805f9b34fb", "12345678-9abc-def0-8000-00805f9b34fb"},
		{"12345678-9abc-def0-1234-56789abcdef0", "12345678-9abc-def0-1234-56789abcdef0"},
	}
	for _, c := range cases {
		t.Run(c.uuid, func(t *testing.T) {
			shortUUID := ShortUUID(c.uuid)
			if shortUUID != c.shortUUID {
				t.Errorf("ShortUUID(%s) == %v, want %v", c.uuid, shortUUID, c.shortUUID)
			}
		})
	}
}
