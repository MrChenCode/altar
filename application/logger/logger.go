package logger

import (
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

//日志文件的切割周期
type Cattime int

const (
	CAT_HOUR Cattime = iota
	CAT_DAY
	CAT_WEEK
	CAT_MONTH
	//永久不切割
	CAT_PERMANENT
	//测试使用
	TEST_CAT_MIN
)

type LoggerOption struct {
	//日志的路径
	LogFilePath string
	//日志切割周期
	CatTime Cattime
	//日志保存时间(天)
	FileExpire int
}

type Logger struct {
	sync.Mutex
	*zap.SugaredLogger
	jack *lumberjack.Logger
}

//创建一个写入到临时位置
func NewTestWriteNull() *Logger {
	w := zapcore.AddSync(ioutil.Discard)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger := &Logger{
		SugaredLogger: zap.New(core).Sugar(),
	}
	return logger
}

func NewConfig(fp string, cattime Cattime, exp int) (*Logger, error) {
	if filepath.Dir(fp) == os.DevNull {
		return NewTestWriteNull(), nil
	}
	return Init(&LoggerOption{
		LogFilePath: fp,
		CatTime:     cattime,
		FileExpire:  exp,
	})
}

func Init(opt *LoggerOption) (*Logger, error) {
	if err := ckopt(opt); err != nil {
		return nil, err
	}
	jack := &lumberjack.Logger{
		Filename: opt.LogFilePath, // 日志文件路径
		MaxSize:  1024 * 100,
		MaxAge:   opt.FileExpire, //最多保留天数
	}
	w := zapcore.AddSync(jack)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "__id",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger := &Logger{
		SugaredLogger: zap.New(core).Sugar(),
		jack:          jack,
	}

	//如果是永久不切割的不做切割处理了
	if opt.CatTime != CAT_PERMANENT {
		crontab := cron.New()
		cronerr := crontab.AddFunc(crontxt(opt.CatTime), func() {
			logger.Lock()
			logger.Sync()
			logger.jack.Rotate()
			logger.Unlock()
		})
		if cronerr != nil {
			return nil, cronerr
		}
		crontab.Start()
	}

	return logger, nil
}

//切割
func (lg *Logger) Rotate() {
	lg.Lock()
	lg.Sync()
	if lg.jack != nil {
		lg.jack.Rotate()
	}
	lg.Unlock()
}

func ckopt(opt *LoggerOption) (err error) {
	//检测日志写入的文件和目录权限
	d := opt.LogFilePath
	if d == "" {
		return errors.New("invalid log file path.")
	}

	fileinfo, err := os.Stat(d)
	if err == nil && fileinfo.IsDir() {
		return fmt.Errorf("logpath(%s) is a directory.", d)
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	dir := filepath.Dir(d)
	dir, err = filepath.Abs(dir)
	if err != nil {
		return
	}
	fileinfo, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("logpath (%s) directory does not exist.", d)
			return
		}
		return
	}

	if syscallAccess(dir, syscall.O_RDWR) != nil {
		err = fmt.Errorf("logpath(%s) lack of directory permissions.", d)
		return
	}

	//检测文件切割周期
	if opt.CatTime < CAT_HOUR && opt.CatTime > TEST_CAT_MIN {
		opt.CatTime = CAT_DAY
	}

	//检测文件过期时间
	if opt.FileExpire < 0 {
		opt.FileExpire = 0
	}

	return nil
}

func crontxt(cat Cattime) string {
	var t string
	switch cat {
	case CAT_HOUR:
		t = "0 0 * * * *"
	case CAT_DAY:
		t = "0 0 0 * * *"
	case CAT_WEEK:
		t = "0 0 0 * * 0"
	case CAT_MONTH:
		t = "0 0 0 1 * *"
	case TEST_CAT_MIN:
		t = "*/2 * * * * *"
	default:
		//默认每天
		t = "0 0 0 * * *"
	}
	return t
}
