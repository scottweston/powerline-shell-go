// +build linux darwin

package writetest

import (
	"golang.org/x/sys/unix"
	"syscall"
)

func IsWritableDir(dir string) bool {
	return syscall.Access(dir, unix.W_OK) == nil
}
