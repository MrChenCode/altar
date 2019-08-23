//+build windows

package logger

func syscallAccess(path string, mode uint32) (err error) {
	return nil
}
