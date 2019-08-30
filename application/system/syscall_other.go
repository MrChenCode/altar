// +build !windows

//兼容windows和非windows syscall包操作
package system

import "syscall"

const (
	O_RDWR = syscall.O_RDWR
)

func SyscallAccess(path string, mode uint32) (err error) {
	return syscall.Access(path, mode)
}
