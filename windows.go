package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// processNames returns a map of PID → executable name for all running processes.
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

// enumerateWindows returns all visible, non-tool, non-child,
// captioned windows owned by the current desktop (not by another window).
func enumerateWindows() []window {
	var windows []window

	cb := syscall.NewCallback(func(hwnd syscall.Handle, _ uintptr) uintptr {
		vis, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
		if vis == 0 {
			return 1 // continue enumeration
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

// filterWindows returns windows that match all provided criteria.
// If no criteria are given the unfiltered slice is returned.
func filterWindows(windows []window, title, process, class string, pid int, names map[uint32]string) []window {
	if title == "" && process == "" && class == "" && pid == 0 {
		return windows
	}

	var result []window
	for _, w := range windows {
		exeName := names[w.pid]
		if matchWindow(w, exeName, title, process, class, pid) {
			result = append(result, w)
		}
	}
	return result
}

// matchWindow checks a single window against title, process, class, and PID filters.
func matchWindow(w window, exeName, title, process, class string, pid int) bool {
	if title != "" && !strings.Contains(strings.ToLower(w.title), strings.ToLower(title)) {
		return false
	}
	if process != "" && !strings.Contains(strings.ToLower(exeName), strings.ToLower(process)) {
		return false
	}
	if class != "" && !strings.Contains(strings.ToLower(w.class), strings.ToLower(class)) {
		return false
	}
	if pid != 0 && w.pid != uint32(pid) {
		return false
	}
	return true
}

// toEntry converts a raw window into its JSON-ready representation.
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

// findWindows validates filters, enumerates windows, and returns matching windows.
// With --all, all matches are returned. Otherwise, only the first match is returned.
func findWindows() ([]window, map[uint32]string, error) {
	if flagTitle == "" && flagProcess == "" && flagPID == 0 && flagClass == "" {
		return nil, nil, fmt.Errorf("at least one filter is required: -t, -p, --pid, or --class")
	}

	windows := enumerateWindows()
	names := processNames()
	matched := filterWindows(windows, flagTitle, flagProcess, flagClass, flagPID, names)

	if len(matched) == 0 {
		return nil, nil, fmt.Errorf("no matching window found")
	}

	if !flagAll {
		return matched[:1], names, nil
	}

	return matched, names, nil
}

// closeWindow sends a WM_CLOSE message to the window, asking it to close gracefully.
func closeWindow(hwnd syscall.Handle) {
	procPostMessageW.Call(uintptr(hwnd), WM_CLOSE, 0, 0)
}

// minimizeWindow minimizes the given window.
func minimizeWindow(hwnd syscall.Handle) {
	procShowWindow.Call(uintptr(hwnd), SW_MINIMIZE)
}

// focusWindow restores the window if minimized and brings it to the foreground.
func focusWindow(hwnd syscall.Handle) {
	iconic, _, _ := procIsIconic.Call(uintptr(hwnd))
	if iconic != 0 {
		procShowWindow.Call(uintptr(hwnd), SW_RESTORE)
	}
	procSetForegroundWindow.Call(uintptr(hwnd))
}
