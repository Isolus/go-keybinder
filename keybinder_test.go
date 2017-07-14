package keybinder

import (
	"testing"
	"unsafe"
)

func init() {
	Init()
}

func keybinderHandler(keystring string, data unsafe.Pointer) {}

func TestBind(t *testing.T) {
	Bind("<Ctrl>space", keybinderHandler, nil)
}
