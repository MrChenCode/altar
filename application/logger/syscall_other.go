//+build !windows

package logger

import "syscall"

func syscallAccess(path string, mode uint32) (err error) {
	return syscall.Access(path, mode)
}
