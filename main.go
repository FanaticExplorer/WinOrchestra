package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

var (
	procEnumWindows              = user32.NewProc("EnumWindows")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowLongW           = user32.NewProc("GetWindowLongW")
	procGetWindow                = user32.NewProc("GetWindow")
	procGetClassNameW            = user32.NewProc("GetClassNameW")

	procCreateToolhelp32Snapshot = kernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32FirstW          = kernel32.NewProc("Process32FirstW")
	procProcess32NextW           = kernel32.NewProc("Process32NextW")
	procCloseHandle              = kernel32.NewProc("CloseHandle")

	procIsIconic            = user32.NewProc("IsIconic")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
)

const (
	GWL_EXSTYLE   = ^uintptr(19) // -20
	GWL_STYLE     = ^uintptr(15) // -16
	GW_OWNER      = 4
	WS_EX_TOOLWIN = 0x00000080
	WS_CHILD      = 0x40000000
	WS_CAPTION    = 0x00C00000

	TH32CS_SNAPPROCESS = 0x00000002
)

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

type windowEntry struct {
	PID       uint32 `json:"pid"`
	Exe       string `json:"exe"`
	Class     string `json:"class"`
	Title     string `json:"title"`
	Minimized bool   `json:"minimized"`
	Focused   bool   `json:"focused"`
}

type window struct {
	handle  syscall.Handle
	title   string
	pid     uint32
	class   string
	style   uintptr
	exStyle uintptr
	owner   uintptr
}

func processNames() map[uint32]string {
	snapshot, _, _ := procCreateToolhelp32Snapshot.Call(TH32CS_SNAPPROCESS, 0)
	if snapshot == 0 {
		return nil
	}
	defer procCloseHandle.Call(snapshot)

	result := make(map[uint32]string)
	var pe PROCESSENTRY32W
	pe.Size = uint32(unsafe.Sizeof(pe))

	r, _, _ := procProcess32FirstW.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
	for r != 0 {
		result[pe.ProcessID] = syscall.UTF16ToString(pe.ExeFile[:])
		pe.Size = uint32(unsafe.Sizeof(pe))
		r, _, _ = procProcess32NextW.Call(snapshot, uintptr(unsafe.Pointer(&pe)))
	}
	return result
}

var (
	flagTitle   string
	flagProcess string
	flagPID     int
	flagList    bool
	flagRaw     bool
)

func init() {
	flag.StringVar(&flagTitle, "title", "", "Filter by window title (partial, case-insensitive)")
	flag.StringVar(&flagProcess, "process", "", "Filter by process .exe name (partial, case-insensitive)")
	flag.IntVar(&flagPID, "pid", 0, "Filter by exact process ID")
	flag.BoolVar(&flagList, "list", false, "List matching windows as JSON")
	flag.BoolVar(&flagRaw, "raw", false, "Output JSON without indentation (compact)")
}

func enumerateWindows() []window {
	var windows []window

	cb := syscall.NewCallback(func(hwnd syscall.Handle, _ uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
		if vis == 0 {
			return 1
		}

		textLen, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
		if textLen == 0 {
			return 1
		}

		style, _, _ := procGetWindowLongW.Call(uintptr(hwnd), GWL_STYLE)
		exStyle, _, _ := procGetWindowLongW.Call(uintptr(hwnd), GWL_EXSTYLE)
		owner, _, _ := procGetWindow.Call(uintptr(hwnd), GW_OWNER)

		if exStyle&WS_EX_TOOLWIN != 0 || style&WS_CHILD != 0 || style&WS_CAPTION == 0 || owner != 0 {
			return 1
		}

		var pid uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&pid)))

		titleBuf := make([]uint16, textLen+1)
		procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&titleBuf[0])), uintptr(len(titleBuf)))

		classBuf := make([]uint16, 256)
		procGetClassNameW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&classBuf[0])), uintptr(len(classBuf)))

		windows = append(windows, window{
			handle:  hwnd,
			title:   syscall.UTF16ToString(titleBuf),
			pid:     pid,
			class:   syscall.UTF16ToString(classBuf),
			style:   style,
			exStyle: exStyle,
			owner:   owner,
		})
		return 1
	})
	procEnumWindows.Call(cb, 0)

	return windows
}

func filterWindows(windows []window, title, process string, pid int, names map[uint32]string) []window {
	if title == "" && process == "" && pid == 0 {
		return windows
	}

	var result []window
	for _, w := range windows {
		exeName := names[w.pid]
		if matchWindow(w, exeName, title, process, pid) {
			result = append(result, w)
		}
	}
	return result
}

func matchWindow(w window, exeName, title, process string, pid int) bool {
	if title != "" && !strings.Contains(strings.ToLower(w.title), strings.ToLower(title)) {
		return false
	}
	if process != "" && !strings.Contains(strings.ToLower(exeName), strings.ToLower(process)) {
		return false
	}
	if pid != 0 && w.pid != uint32(pid) {
		return false
	}
	return true
}

func toEntry(w window, names map[uint32]string, foregroundHwnd syscall.Handle) windowEntry {
	minimized, _, _ := procIsIconic.Call(uintptr(w.handle))

	return windowEntry{
		PID:       w.pid,
		Exe:       names[w.pid],
		Class:     w.class,
		Title:     w.title,
		Minimized: minimized != 0,
		Focused:   syscall.Handle(w.handle) == foregroundHwnd,
	}
}

func main() {
	flag.Parse()

	if !flagList {
		flag.Usage()
		os.Exit(0)
	}

	windows := enumerateWindows()
	names := processNames()

	matched := filterWindows(windows, flagTitle, flagProcess, flagPID, names)
	foregroundHwnd, _, _ := procGetForegroundWindow.Call()
	entries := make([]windowEntry, 0, len(matched))
	for _, w := range matched {
		entries = append(entries, toEntry(w, names, syscall.Handle(foregroundHwnd)))
	}

	var out []byte
	var err error
	if flagRaw {
		out, err = json.Marshal(entries)
	} else {
		out, err = json.MarshalIndent(entries, "", "  ")
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
