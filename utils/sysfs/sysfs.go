package sysfs

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var CreatePermissions os.FileMode = 0755

type SysFs interface {
	DirExists(dirPath string) bool
	FileExists(filePath string) bool
	Read(filePath string) (io.ReadCloser, error)
	Write(filePath string) (io.WriteCloser, error)
	DeleteFile(filePath string) error
	DeleteDir(dirPath string) error
	ReadDir(dirPath string) ([]os.FileInfo, error)
	CreateDir(dirPath string) error
}

type sysFs struct{}

var _ SysFs = &sysFs{}

func New() SysFs {
	return &sysFs{}
}

func (fs *sysFs) DeleteDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func (fs *sysFs) DirExists(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err != nil || !stat.IsDir() {
		return false
	}

	return true
}

func (fs *sysFs) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return true
}

func (fs *sysFs) Read(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

func (fs *sysFs) Write(filePath string) (io.WriteCloser, error) {
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, CreatePermissions)
}

func (fs *sysFs) DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

func (fs *sysFs) ReadDir(dirPath string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirPath)
}

func (fs *sysFs) CreateDir(dirPath string) error {
	return os.Mkdir(dirPath, CreatePermissions)
}

func WriteJson(fs SysFs, filePath string, v interface{}) error {
	w, err := fs.Write(filePath)
	if err != nil {
		return err
	}

	defer w.Close()
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func ReadJson(fs SysFs, filePath string, v interface{}) error {
	r, err := fs.Read(filePath)
	if err != nil {
		return err
	}

	defer r.Close()
	d := json.NewDecoder(r)
	if err := d.Decode(v); err != nil {
		return err
	}

	return nil
}

func ReadAll(fs SysFs, filePath string) ([]byte, error) {
	r, err := fs.Read(filePath)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return ioutil.ReadAll(r)
}

func ClearDir(fs SysFs, dirPath string) error {
	files, err := fs.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		p := filepath.Join(dirPath, f.Name())
		if f.IsDir() {
			err = fs.DeleteDir(p)
		} else {
			err = fs.DeleteFile(p)
		}

		if err != nil {
			return err
		}
	}

	return nil

}

func CopyFile(fs SysFs, srcPath, dstPath string) (int64, error) {
	src, err := fs.Read(srcPath)
	if err != nil {
		return -1, err
	}

	defer src.Close()
	dst, err := fs.Write(dstPath)
	if err != nil {
		return -1, err
	}

	defer dst.Close()
	return io.Copy(dst, src)
}

func SizeOf(r io.ReadCloser) (int64, error) {
	fr, ok := r.(*os.File)
	if !ok {
		return -1, errors.New("couldn't read size: not a valid file reader")
	}

	fi, err := fr.Stat()
	if err != nil {
		return -1, err
	}

	return fi.Size(), nil
}
