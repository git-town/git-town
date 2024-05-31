package windows

import "golang.org/x/sys/windows"

var NewLazySystemDLL = windows.NewLazySystemDLL

type Handle = windows.Handle

//sys GetKeyboardLayout(threadId uint32) (hkl Handle) = user32.GetKeyboardLayout
//sys ToUnicodeEx(vkey uint32, scancode uint32, keystate *byte, pwszBuff *uint16, cchBuff int32, flags uint32, hkl Handle) (ret int32) = user32.ToUnicodeEx
