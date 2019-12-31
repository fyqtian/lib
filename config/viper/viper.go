package viper

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	ErrNotExists = errors.New("config not exists")
	V            *Helper
)

type Helper struct {
	options *Options
	*viper.Viper
}

type Options struct {
	ConfigPath []string
	FileName   string
}

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
	V, _ = NewViper(DefaultOptions())
}
func GetSingleton() *Helper {
	return V
}

func NewViper(option *Options) (*Helper, error) {
	var err error
	var h = &Helper{}

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
	return h, nil
}
