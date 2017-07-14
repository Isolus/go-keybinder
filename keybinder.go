// Go bindings for keybinder-3.0.
//
// Functions use the same names as the native C function calls, but use CamelCase.
//
// The keybinder-3.0 documentation can be very useful for understanding how the
// functions in this package work. This documentation can be found at
// https://github.com/kupferlauncher/keybinder.
package keybinder

// #cgo pkg-config: keybinder-3.0
// #include "keybinder.go.h"
import "C"
import (
	"reflect"
	"sync"
	"unsafe"
)

func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func gobool(b C.gboolean) bool {
	return b != C.FALSE
}

type KeybinderHandler func(keystring string, userData unsafe.Pointer)

type keybinderHandlerData struct {
	fn        KeybinderHandler
	keystring string
	data      unsafe.Pointer
}

var (
	keybinderHandlerRegistry = struct {
		sync.RWMutex
		next int
		m    map[int]keybinderHandlerData
	}{
		next: 1,
		m:    make(map[int]keybinderHandlerData),
	}
)

// Init() is a wrapper around keybinder_init().
func Init() {
	C.keybinder_init()
}

// SetUseCookedAccelerators() is a wrapper around keybinder_set_use_cooked_accelerators().
func SetUseCookedAccelerators(cooked bool) {
	C.keybinder_set_use_cooked_accelerators(gbool(cooked))
}

// Bind() is a wrapper around keybinder_bind().
func Bind(keystring string, cb KeybinderHandler, data unsafe.Pointer) bool {

	keybinderHandlerRegistry.Lock()
	id := keybinderHandlerRegistry.next
	keybinderHandlerRegistry.next++
	keybinderHandlerRegistry.m[id] =
		keybinderHandlerData{fn: cb, keystring: keystring, data: data}
	keybinderHandlerRegistry.Unlock()

	cstr := C.CString(keystring)
	defer C.free(unsafe.Pointer(cstr))
	c := C._keybinder_bind(cstr, unsafe.Pointer(uintptr(id)))

	return gobool(c)
}

/*
// BindFull() is a wrapper around keybinder_bind_full().
func BindFull(keystring string, handler KeybinderHandler, data uintptr, notify GDestroyNotify) bool {
	// TODO: Implement GDestroyNotify
}
*/

// Unbind() is a wrapper around keybinder_unbind().
func Unbind(keystring string, handler KeybinderHandler) {
	cstr := C.CString(keystring)
	defer C.free(unsafe.Pointer(cstr))

	var (
		handlers int = 0
		id       int = -1
	)

	keybinderHandlerRegistry.Lock()
	for k, v := range keybinderHandlerRegistry.m {
		if v.keystring == keystring {
			handlers++
		}
		h1 := reflect.ValueOf(v.fn).Pointer()
		h2 := reflect.ValueOf(handler).Pointer()
		if v.keystring == keystring && h1 == h2 {
			id = k
		}
	}
	keybinderHandlerRegistry.Unlock()

	if id > -1 {
		if handlers <= 1 {
			C._keybinder_unbind(cstr)
		}
		keybinderHandlerRegistry.Lock()
		delete(keybinderHandlerRegistry.m, id)
		keybinderHandlerRegistry.Unlock()
	}

}

// UnbindAll() is a wrapper around keybinder_unbind_all().
func UnbindAll(keystring string) {
	cstr := C.CString(keystring)
	defer C.free(unsafe.Pointer(cstr))
	C.keybinder_unbind_all(cstr)

	keybinderHandlerRegistry.Lock()
	for k, v := range keybinderHandlerRegistry.m {
		if v.keystring == keystring {
			delete(keybinderHandlerRegistry.m, k)
		}
	}
	keybinderHandlerRegistry.Unlock()
}

// GetCurrentEventTime() is a wrapper around keybinder_get_current_event_time().
func GetCurrentEventTime() uint32 {
	return uint32(C.keybinder_get_current_event_time())
}

// Supported() is a wrapper around keybinder_supported().
func Supported() bool {
	return gobool(C.keybinder_supported())
}

//export goKeybinderHandler
func goKeybinderHandler(keystring *C.gchar,
	data unsafe.Pointer) {

	id := int(uintptr(data))

	keybinderHandlerRegistry.Lock()
	r := keybinderHandlerRegistry.m[id]
	keybinderHandlerRegistry.Unlock()

	r.fn(C.GoString((*C.char)(keystring)), r.data)

}
