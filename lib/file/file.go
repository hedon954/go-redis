package file

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// CheckNotExist returns true means not exist
func CheckNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// CheckPermission return true means permission denied
func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// MkDir creates dir
func MkDir(src string) error {
	syscall.Umask(0)
	return os.MkdirAll(src, os.ModePerm)
}

// IsNotExistMkdir makes dir if it not exists
func IsNotExistMkdir(src string) error {
	if notExist := CheckNotExist(src); notExist {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// OpenFile opens file
func OpenFile(filename, dir string) (*os.File, error) {
	if CheckPermission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	if err := IsNotExistMkdir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir: %s, err: %s", dir, err)
	}

	syscall.Umask(0)

	f, err := os.OpenFile(filepath.Join(dir, filename), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}
	return f, nil
}