package zap

import (
	"errors"
	"github.com/fyqtian/lib/log/rotate"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"sync"
	"time"
)

type Helper struct {
	options *Options
	*zap.Logger
}
type Field = zap.Field

type Options struct {
	//存储位置
	FilePath string
	//每个轮转日志大小
	FileSize int
	//保存几份日志
	FileBackup int
	//保存时间
	FileMaxAge int
	//是否压缩
	FileCompress bool
	//是否打印到os.stdout
	Debug bool
	//输出级别
	Level string
	//动态调整debug级别端口
	Listen string
}

var (
	store     = sync.Map{}
	NotExists = errors.New("zap not exists")
)

func Get(prefix string) (*Helper, error) {
	if v, ok := store.Load(prefix); !ok {
		return nil, NotExists
	} else {
		val, _ := v.(*Helper)
		return val, nil
	}
}

func NewZap(prefix string, option *Options) (*Helper, error) {
	var err error
	var h = &Helper{}
	if h, err := Get(prefix); err == nil {
		return h, nil
	}
	h = &Helper{
		options: option,
	}
	//日志滚动
	rotate := rotate.NewRotate(h.options.FilePath, h.options.FileSize, h.options.FileBackup, h.options.FileMaxAge, h.options.FileCompress)
	var arr []zapcore.Core
	//可通过http请求调整level
	autoLevel := zap.NewAtomicLevelAt(h.logLever(h.options.Level))
	core := zapcore.NewCore(h.defaultEncoder(), zapcore.AddSync(rotate), autoLevel)

	//todo
	// 如果配置了端口认为开启了http接口调整日志级别
	if h.options.Listen != "" {
		h.registerHttp(&autoLevel, h.options.Listen)
	}
	arr = append(arr, core)
	if h.options.Debug {
		tmp := zapcore.NewCore(h.defaultEncoder(), os.Stdout, autoLevel)
		arr = append(arr, tmp)
	}

	h.Logger = zap.New(zapcore.NewTee(arr...), zap.AddCaller())
	store.Store(prefix, h)
	return h, err
}

//panic
func (s *Helper) registerHttp(l *zap.AtomicLevel, port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/zap/handle/level", l.ServeHTTP)
	go func() {
		err := http.ListenAndServe(port, mux)
		if err != nil {
			panic(err)
		}
	}()
}

func (s *Helper) logLever(l string) zapcore.Level {
	var r zapcore.Level
	switch l {
	case "debug":
		r = zapcore.DebugLevel
	case "info":
		r = zapcore.InfoLevel
	case "warn":
		r = zapcore.WarnLevel
	case "error":
		r = zapcore.ErrorLevel
	case "fatal":
		r = zapcore.FatalLevel
	default:
		r = zapcore.InfoLevel
	}
	return r
}

func (s *Helper) defaultEncoder() zapcore.Encoder {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("2006-01-02 15:04:05")
	}
	return zapcore.NewJSONEncoder(c)
}
