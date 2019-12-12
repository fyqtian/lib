package rotate

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewRotate(filepath string, maxsize, maxbackups, maxage int, compress bool) *lumberjack.Logger {
	hook := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    maxsize, // megabytes
		MaxBackups: maxbackups,
		MaxAge:     maxage,   //days
		Compress:   compress, // disabled by default
	}
	return hook
}
