package rctx

const (
	LoggerInfoBufferCap  = 20
	LoggerErrorBufferCap = 4
)

type log struct {
	//暂存log数据
	infoKV  []interface{}
	errorKV []interface{}
}

func (l *log) reset() {
	l.infoKV = nil
	l.errorKV = nil
}

func (l *log) Info(kvs ...interface{}) {
	if l.infoKV == nil {
		l.infoKV = make([]interface{}, 0, LoggerInfoBufferCap)
	}
	l.infoKV = append(l.infoKV, kvs...)
}

func (l *log) Error(kvs ...interface{}) {
	if l.errorKV == nil {
		l.errorKV = make([]interface{}, 0, LoggerErrorBufferCap)
	}
	l.errorKV = append(l.errorKV, kvs...)
}
