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
	if l.infoKV != nil {
		l.infoKV = l.infoKV[:0]
	}
	if l.errorKV != nil {
		l.errorKV = l.errorKV[:0]
	}
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
