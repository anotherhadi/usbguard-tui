package ui

import (
	"strings"

	"github.com/anotherhadi/ilovetui/style"
)

const (
	iconUSB       = " "
	iconKeyboard  = ""
	iconPointer   = ""
	iconCamera    = ""
	iconStorage   = ""
	iconAudio     = ""
	iconBluetooth = ""
	iconNetwork   = ""
	iconPrint     = ""
	iconHub       = ""
	iconPhone     = ""
	iconGamepad   = ""
	iconSecurity  = ""
)

type deviceIconCategory struct {
	icon     string
	keywords []string
}

var deviceIconCategories = []deviceIconCategory{
	{iconKeyboard, []string{"keyboard"}},
	{iconPointer, []string{"mouse", "touchpad", "trackpad"}},
	{iconCamera, []string{"webcam", "camera"}},
	{iconStorage, []string{"storage", "disk", "drive", "ssd", "data", "flash", "card reader"}},
	{iconAudio, []string{"headset", "headphone", "audio", "speaker", "microphone"}},
	{iconBluetooth, []string{"bluetooth"}},
	{iconNetwork, []string{"wifi", "wi-fi", "wireless", "ethernet", "network adapter", "network card"}},
	{iconPrint, []string{"printer", "scanner"}},
	{iconHub, []string{"hub"}},
	{iconPhone, []string{"phone", "android", "iphone"}},
	{iconGamepad, []string{"gamepad", "joystick"}},
	{iconSecurity, []string{"security key", "smartcard", "smart card", "yubikey", "fingerprint"}},
}

func deviceIcon(name string) string {
	if !style.S.NerdFonts {
		return ""
	}
	lower := strings.ToLower(name)
	for _, c := range deviceIconCategories {
		for _, kw := range c.keywords {
			if strings.Contains(lower, kw) {
				return c.icon + " "
			}
		}
	}
	return iconUSB
}
