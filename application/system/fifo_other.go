package system

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

//fifo通信有名管道
var PIPFile = filepath.Join(os.TempDir(), "altar.ipc")
var PIPEOF = "EOF"
var PIPBUF = "BUF"
var PIPERR = "ERR"

func PipNew() error {
	return syscall.Mkfifo(PIPFile, 0666)
}

func PipClose() {
	_ = os.Remove(PIPFile)
}

func PipCloseNew() error {
	PipClose()
	return PipNew()
}

func PipRead() chan string {
	c := make(chan string)
	go func(c chan string) {
		f, err := os.OpenFile(PIPFile, os.O_RDWR, os.ModeNamedPipe)
		if err != nil {
			c <- PIPERR + err.Error()
			c <- PIPEOF
			return
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					c <- PIPERR + err.Error()
				}
				c <- PIPEOF
				return
			}
			if len(line) < 3 {
				c <- PIPERR + "无效的FIFO数据"
				c <- PIPEOF
				return
			}
			if string(line[:3]) == PIPEOF {
				c <- PIPEOF
				return
			}
			c <- string(line)
		}
	}(c)
	return c
}

func PipWrite(s, te string) {
	if te != PIPEOF && te != PIPBUF && te != PIPERR {
		return
	}
	f, err := os.OpenFile(PIPFile, os.O_RDWR, 0777)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(fmt.Sprintf("%s%s\n", te, s))
	if te == PIPERR {
		_, _ = f.WriteString(fmt.Sprintf("%s\n", PIPEOF))
	}
}
