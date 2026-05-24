package main

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func main() {
	exe, err := os.Executable()
	if err != nil {
		showMessage("Seall could not find its install folder.")
		return
	}

	root := filepath.Dir(exe)
	target := filepath.Join(root, "seall-denshi", "dist", "win-unpacked", "Seall.exe")
	if _, err := os.Stat(target); err != nil {
		showMessage("Seall could not find the desktop app.\n\nMissing:\n" + target)
		return
	}

	if err := start(target); err != nil {
		showMessage("Seall could not start.\n\n" + err.Error())
	}
}

func start(target string) error {
	argv := []string{target}
	attr := &os.ProcAttr{
		Dir:   filepath.Dir(target),
		Files: []*os.File{nil, nil, nil},
	}

	process, err := os.StartProcess(target, argv, attr)
	if err != nil {
		return err
	}
	return process.Release()
}

func showMessage(message string) {
	user32 := syscall.NewLazyDLL("user32.dll")
	messageBox := user32.NewProc("MessageBoxW")

	titlePtr, _ := syscall.UTF16PtrFromString("Seall")
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	messageBox.Call(
		0,
		uintptr(unsafe.Pointer(messagePtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		0x00000010,
	)
}
