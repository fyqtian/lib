package zap

import (
	"errors"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/log/rotate"
	"github.com/fyqtian/lib/utils"
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
	//每个轮转日志大小 单位mb
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
	ErrNotExists = errors.New("zap not exists")
	once         sync.Once
	Logger       *Helper
)

func loadFromConfiger(prefix string, c config.Configer) *Options {
	p := prefix + "."
	return &Options{
		FilePath:     c.GetString(utils.CombineString(p, "filepath")),
		FileSize:     c.GetInt(utils.CombineString(p, "filesize")),
		FileBackup:   c.GetInt(utils.CombineString(p, "filebackup")),
		FileMaxAge:   c.GetInt(utils.CombineString(p, "filemaxage")),
		FileCompress: c.GetBool(utils.CombineString(p, "filecompress")),
		Debug:        c.GetBool(utils.CombineString(p, "debug")),
		Level:        c.GetString(utils.CombineString(p, "level")),
		Listen:       c.GetString(utils.CombineString(p, "listen")),
	}
}

func SampleOptions(prefix string, c config.Configer) *Options {
	op := loadFromConfiger(prefix, c)

	if op.FilePath == "" {
		op.FilePath = "Runtime/runtime.log"
	}

	if op.FileSize == 0 {
		op.FileSize = 128
	}

	if op.FileBackup == 0 {
		op.FileBackup = 10
	}

	if op.FileMaxAge == 0 {
		op.FileMaxAge = 7
	}
	return op
}

func DefaultZap() *Helper {
	once.Do(func() {
		Logger = NewZap(SampleOptions("log", viper.GetSingleton()))
	})
	return Logger
}

func NewZap(option *Options) *Helper {
	var h = &Helper{
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
	return h
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
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	return zapcore.NewJSONEncoder(c)
}
