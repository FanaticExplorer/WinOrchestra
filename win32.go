package main

import "syscall"

// DLLs
var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

// Win32 API function pointers
var (
	procEnumWindows              = user32.NewProc("EnumWindows")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowLongW           = user32.NewProc("GetWindowLongW")
	procGetWindow                = user32.NewProc("GetWindow")
	procGetClassNameW            = user32.NewProc("GetClassNameW")
	procIsIconic                 = user32.NewProc("IsIconic")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")

	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
)

// Win32 constants
const (
	GWL_EXSTYLE   = ^uintptr(19) // -20
	GWL_STYLE     = ^uintptr(15) // -16
	GW_OWNER      = 4
	WS_EX_TOOLWIN = 0x00000080
	WS_CHILD      = 0x40000000
	WS_CAPTION    = 0x00C00000

	TH32CS_SNAPPROCESS = 0x00000002
)

// Win32 structs
type PROCESSENTRY32W struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [260]uint16
}

// window holds raw Win32 data for one window.
type window struct {
	handle  syscall.Handle
	title   string
	pid     uint32
	class   string
	style   uintptr
	exStyle uintptr
	owner   uintptr
}

// windowEntry is the JSON-ready representation of a window.
type windowEntry struct {
	PID       uint32 `json:"pid"`
	Exe       string `json:"exe"`
	Class     string `json:"class"`
	Title     string `json:"title"`
	Minimized bool   `json:"minimized"`
	Focused   bool   `json:"focused"`
}
