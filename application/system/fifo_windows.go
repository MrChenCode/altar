// +build windows

package system

//fifo通信有名管道
var (
	PIPFile = ""
	PIPEOF  = "EOF"
	PIPBUF  = "BUF"
	PIPERR  = "ERR"
)

func PipNew() error {
	return nil
}

func PipClose() {
	return
}

func PipCloseNew() error {
	return nil
}

func PipRead() chan string {
	c := make(chan string)
	c <- PIPEOF
	return c
}

func PipWrite(s, te string) {
	_, _ = s, te
}
