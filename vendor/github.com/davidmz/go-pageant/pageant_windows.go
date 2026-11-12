// Package pageant provides an interface to the PuTTY Pageant (pageant.exe)
// SSH agent. This package is Windows-only.
package pageant

// see https://github.com/Yasushi/putty/blob/master/windows/winpgntc.c#L155
// see https://github.com/paramiko/paramiko/blob/master/paramiko/win_pageant.py

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// MaxMessageLen is the maximum size (including the 4-byte length prefix) of
// a request to or a response from Pageant. It matches AGENT_MAX_MSGLEN of
// modern PuTTY versions.
const MaxMessageLen = 256 * 1024

// Timeout limits how long a single request waits for Pageant's reply.
// It is generous because Pageant may ask the user for a key passphrase.
// Set it before using the agent.
var Timeout = 5 * time.Minute

var (
	// ErrPageantNotFound is returned when the Pageant window is not found.
	ErrPageantNotFound = errors.New("pageant process not found")
	// ErrSendMessage is returned when Pageant rejects or fails the request.
	ErrSendMessage = errors.New("error sending message")

	// ErrMessageTooLong is returned when the request exceeds MaxMessageLen.
	ErrMessageTooLong = errors.New("message too long")
	// ErrInvalidMessageFormat is returned when the request's length prefix
	// does not match its actual length.
	ErrInvalidMessageFormat = errors.New("invalid message format")
	// ErrResponseTooLong is returned when Pageant's response exceeds MaxMessageLen.
	ErrResponseTooLong = errors.New("response too long")
)

const (
	agentCopydataID = 0x804e50ba
	wmCopydata      = 74
	smtoAbortIfHung = 0x0002
)

type copyData struct {
	dwData uintptr
	cbData uint32
	lpData unsafe.Pointer
}

var (
	// Pointers must be converted to uintptr directly in the arguments of
	// LazyProc.Call: only then Go keeps the pointed objects alive and in place.
	user32                 = windows.NewLazySystemDLL("user32.dll")
	procFindWindow         = user32.NewProc("FindWindowW")
	procSendMessageTimeout = user32.NewProc("SendMessageTimeoutW")

	requestCounter uint32
)

// Available returns true if Pageant is running.
func Available() bool { return pageantWindow() != 0 }

// query sends msg to Pageant and returns its response.
// Both msg and the response are raw agent messages with the length prefix.
func query(msg []byte) ([]byte, error) {
	if len(msg) > MaxMessageLen {
		return nil, ErrMessageTooLong
	}

	if len(msg) < 4 || uint32(len(msg)-4) != binary.BigEndian.Uint32(msg) {
		return nil, ErrInvalidMessageFormat
	}

	paWin := pageantWindow()
	if paWin == 0 {
		return nil, ErrPageantNotFound
	}

	// Unique per request: after a timeout Pageant may still hold the old
	// mapping open, and a reused name would hand us its stale contents.
	mapName := fmt.Sprintf("PageantRequest%08x%08x", os.Getpid(), atomic.AddUint32(&requestCounter, 1))
	pMapName, err := windows.UTF16PtrFromString(mapName)
	if err != nil {
		return nil, err
	}

	sa, err := ownedByCurrentUser()
	if err != nil {
		return nil, err
	}

	mmap, err := windows.CreateFileMapping(windows.InvalidHandle, sa, windows.PAGE_READWRITE, 0, MaxMessageLen, pMapName)
	if mmap != 0 {
		defer windows.CloseHandle(mmap)
	}
	if err != nil {
		// ERROR_ALREADY_EXISTS: someone else created a mapping with our name
		// and could read or forge the exchange.
		return nil, err
	}

	ptr, err := windows.MapViewOfFile(mmap, windows.FILE_MAP_WRITE, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	defer windows.UnmapViewOfFile(ptr)

	// ptr is an address outside the Go heap, so the conversion is safe;
	// it is done through a pointer to keep go vet's unsafeptr check quiet.
	mmSlice := unsafe.Slice(*(**byte)(unsafe.Pointer(&ptr)), MaxMessageLen)

	copy(mmSlice, msg)

	mapNameBytesZ := append([]byte(mapName), 0)

	cds := copyData{
		dwData: agentCopydataID,
		cbData: uint32(len(mapNameBytesZ)),
		lpData: unsafe.Pointer(&(mapNameBytesZ[0])),
	}

	var result uintptr
	ok, _, err := procSendMessageTimeout.Call(
		paWin, wmCopydata, 0, uintptr(unsafe.Pointer(&cds)),
		smtoAbortIfHung, uintptr(Timeout.Milliseconds()), uintptr(unsafe.Pointer(&result)),
	)
	if ok == 0 {
		return nil, fmt.Errorf("pageant: %w", err) // windows.ERROR_TIMEOUT on timeout
	}
	if result == 0 {
		return nil, ErrSendMessage
	}

	respLen := binary.BigEndian.Uint32(mmSlice)
	if respLen > MaxMessageLen-4 {
		return nil, ErrResponseTooLong
	}

	respData := make([]byte, respLen+4)
	copy(respData, mmSlice)

	return respData, nil
}

// ownedByCurrentUser returns security attributes that make the current user
// (not the Administrators group, which is the default owner in an elevated
// process) the owner of the mapping: Pageant rejects mappings owned by others.
func ownedByCurrentUser() (*windows.SecurityAttributes, error) {
	tu, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, err
	}
	sid := tu.User.Sid.String()
	sd, err := windows.SecurityDescriptorFromString("O:" + sid + "D:P(A;;GA;;;" + sid + ")")
	if err != nil {
		return nil, err
	}
	sa := &windows.SecurityAttributes{SecurityDescriptor: sd}
	sa.Length = uint32(unsafe.Sizeof(*sa))
	return sa, nil
}

func pageantWindow() uintptr {
	// No user32.dll (e.g. Nano Server) means no Pageant, not a panic.
	if procFindWindow.Find() != nil || procSendMessageTimeout.Find() != nil {
		return 0
	}
	nameP, _ := windows.UTF16PtrFromString("Pageant")
	h, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(nameP)), uintptr(unsafe.Pointer(nameP)))
	return h
}
