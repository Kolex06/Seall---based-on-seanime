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
	packagedDesktop := filepath.Join(root, "seall-denshi", "dist", "win-unpacked", "Seall.exe")
	sourceElectron := filepath.Join(root, "seall-denshi", "node_modules", "electron", "dist", "electron.exe")
	denshiDir := filepath.Join(root, "seall-denshi")

	if _, err := os.Stat(packagedDesktop); err == nil {
		if err := start(packagedDesktop); err != nil {
			showMessage("Seall could not start.\n\n" + err.Error())
		}
		return
	}

	if _, err := os.Stat(sourceElectron); err == nil {
		if err := startWithArgs(sourceElectron, []string{sourceElectron, denshiDir}, denshiDir); err != nil {
			showMessage("Seall could not start.\n\n" + err.Error())
		}
		return
	}

	showMessage("Seall could not find the desktop app.\n\nMissing:\n" + packagedDesktop + "\n" + sourceElectron)
}

func start(target string) error {
	return startWithArgs(target, []string{target}, filepath.Dir(target))
}

func startWithArgs(target string, argv []string, dir string) error {
	process, err := os.StartProcess(target, argv, &os.ProcAttr{
		Dir:   dir,
		Files: []*os.File{nil, nil, nil},
	})
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
