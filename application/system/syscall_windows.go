// +build windows

//兼容windows和非windows syscall包操作
package system

const (
	O_RDWR = 0
)

func SyscallAccess(path string, mode uint32) (err error) {
	return nil
}
