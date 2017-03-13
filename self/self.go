// +build !windows

package self

import (
	"fmt"
	"github.com/kardianos/osext"
	"os"
	"path/filepath"
	"syscall"

	"github.com/danielkrainas/shex/fsutils"
)

var (
	defaultInstallPath string = "/var/lib/shex"
	symlinkPath        string = "/usr/local/bin/shex"
)

func createSymLink(shexPath string) error {
	if fsutils.FileExists(symlinkPath) {
		return fmt.Errorf("could not create symlink, file already exists: %s\n", symlinkPath)
	}

	return os.Symlink(shexPath, symlinkPath)
}

func removeSymLink() error {
	if !fsutils.FileExists(symlinkPath) {
		return nil
	}

	return syscall.Unlink(symlinkPath)
}

func Uninstall(installPath string) error {
	if installPath == "" {
		installPath = defaultInstallPath
	}

	if !fsutils.DirExists(installPath) {
		fmt.Printf("no installation found at %s\n", installPath)
		return nil
	}

	realSrc, err := os.Readlink(symlinkPath)
	if err == nil && filepath.Dir(realSrc) == installPath {
		// TODO: add switch to disable symlink stuff
		// TODO: add some logging for this
		removeSymLink()
	}

	if err = os.RemoveAll(installPath); err != nil {
		return err
	}

	fmt.Printf("uninstalled: %s\n", installPath)
	return nil
}

func Install(dest string) error {
	src, err := osext.Executable()
	if err != nil {
		return err
	}

	if dest == "" {
		dest = defaultInstallPath
	}

	existed := fsutils.DirExists(dest)
	err = os.MkdirAll(dest, 0777)
	if err != nil {
		return err
	}

	selfDest := filepath.Join(dest, filepath.Base(src))
	if selfDest == src || (src == symlinkPath && selfDest == defaultInstallPath) {
		fmt.Printf("already installed at %s\n", dest)
		return nil
	}

	_, err = fsutils.CopyFile(src, selfDest)
	if err != nil {
		return err
	}

	if err = os.Chmod(selfDest, 0777); err != nil {
		if !existed {
			_ = os.RemoveAll(dest)
		} else {
			_ = os.Remove(selfDest)
		}

		_ = removeSymLink()
		return err
	}

	// TODO: add switch to disable symlink stuff
	if err = createSymLink(selfDest); err != nil {
		return err
	}

	fmt.Printf("installed:\n %s => %s\n", src, selfDest)
	return nil
}
