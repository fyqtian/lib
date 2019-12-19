package viper

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
)

var (
	store     = sync.Map{}
	NotExists = errors.New("config not exists")
	V         *Helper
)

type Helper struct {
	options *Options
	*viper.Viper
}

type Options struct {
	ConfigPath []string
	FileName   string
}

//func Get(prefix string) (*Helper, error) {
//	if v, ok := store.Load(prefix); !ok {
//		return nil, NotExists
//	} else {
//		val, _ := v.(*Helper)
//		return val, nil
//	}
//}

func DefaultOptions() *Options {
	path1, _ := filepath.Abs("config")
	path2, _ := filepath.Abs("configs")
	path3, _ := filepath.Abs(".")
	return &Options{
		ConfigPath: []string{path1, path2, path3},
		FileName:   os.Getenv("RUN_TIME"),
	}
}
func init() {
	//ignore error
	V, _ = NewViper(nil)
}
func GetSingleton() *Helper {
	return V
}

func NewViper(option *Options) (*Helper, error) {
	var err error
	var h = &Helper{}

	//if h, err := Get(prefix); err == nil {
	//	return h, nil
	//}
	if option == nil {
		option = DefaultOptions()
	}
	h.options = option
	v := viper.New()
	v.SetConfigName(option.FileName)

	for _, path := range option.ConfigPath {
		v.AddConfigPath(path)
	}

	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.WatchConfig()
	h.Viper = v
	//store.Store(prefix, h)
	return h, nil
}
