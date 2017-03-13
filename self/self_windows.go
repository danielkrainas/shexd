package self

import (
	"fmt"
	"github.com/kardianos/osext"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

const (
	HWND_BROADCAST   uintptr = 0xffff
	WM_SETTINGCHANGE uint    = 0x1a
)

var (
	libUser32                 = syscall.NewLazyDLL("user32.dll")
	defaultInstallPath string = `${APPDATA}\shex`
	pathRegKey         string = "Path"
)

func broadcastSettingChange() {
	rawParam := "ENVIRONMENT"
	param := uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(rawParam)))
	sendMessageProcedure := libUser32.NewProc("SendMessageW")
	sendMessageProcedure.Call(uintptr(HWND_BROADCAST), uintptr(WM_SETTINGCHANGE), 0, param)
}

func getRegEnvValue(key string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}

	defer k.Close()
	s, _, err := k.GetStringValue(key)
	return s, err
}

func saveRegEnvValue(key string, value string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer k.Close()
	return k.SetStringValue(key, value)
}

func addInstallationToPath(installPath string) error {
	v, err := getRegEnvValue(pathRegKey)
	if err != nil {
		return err
	}

	v = v + string(os.PathListSeparator) + installPath
	// TODO: do logging here
	return saveRegEnvValue(pathRegKey, v)
}

func removePathInstallation(installPath string) error {
	sep := string(os.PathListSeparator)
	v, err := getRegEnvValue(pathRegKey)
	if err != nil {
		return err
	}

	paths := strings.Split(v, sep)
	for i, p := range paths {
		if p == installPath {
			continue
		}

		if i == 0 {
			v = p
		} else {
			v += sep + p
		}
	}

	// TODO: do logging here
	return saveRegEnvValue(pathRegKey, v)
}

func removeSelf(installPath string) error {
	cmd := exec.Command("cmd.exe", "/C", "choice", "/C", "Y", "/N", "/D", "Y", "/T", "3", "&", "rd", "/s", "/q", installPath)
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func getDefaultInstallPath() string {
	return os.ExpandEnv(defaultInstallPath)
}

func installSelf(dest string) error {
	src, err := osext.Executable()
	if err != nil {
		return err
	}

	if dest == "" {
		dest = getDefaultInstallPath()
	}

	existed := dirExists(dest)
	err = os.MkdirAll(dest, 0777)
	if err != nil {
		return err
	}

	selfDest := filepath.Join(dest, filepath.Base(src))
	if selfDest == src {
		fmt.Printf("already installed at %s\n", dest)
		return nil
	}

	_, err = copyFile(src, selfDest)
	if err != nil {
		return err
	}

	if err = os.Chmod(selfDest, 0777); err != nil {
		if !existed {
			_ = os.RemoveAll(dest)
		} else {
			_ = os.Remove(selfDest)
		}

		// TODO: add switch to disable symlink/$PATH creation
		_ = removePathInstallation(dest)
		return err
	}

	// TODO: add switch to disable symlink/$PATH creation
	if err = addInstallationToPath(dest); err != nil {
		return err
	}

	broadcastSettingChange()
	fmt.Printf("installed:\n %s => %s\n", src, selfDest)
	return nil
}

func uninstallSelf(installPath string) error {
	selfPath, err := osext.Executable()
	if err != nil {
		return err
	}

	if installPath == "" {
		installPath = getDefaultInstallPath()
	}

	if !dirExists(installPath) {
		fmt.Printf("no installation found at %s\n", installPath)
		return nil
	}

	// TODO: add switch to disable symlink stuff
	// TODO: add some logging for this
	err = removePathInstallation(installPath)
	if err != nil {
		return err
	}

	broadcastSettingChange()
	if filepath.Dir(selfPath) == installPath {
		if err = removeSelf(installPath); err != nil {
			return err
		}
	} else if err = os.RemoveAll(installPath); err != nil {
		return err
	}

	fmt.Printf("uninstalled: %s\n", installPath)
	return nil
}
