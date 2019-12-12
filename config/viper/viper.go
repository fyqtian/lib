package viper

import (
	"errors"
	"github.com/spf13/viper"
	"sync"
)

var (
	store     = sync.Map{}
	NotExists = errors.New("config not exists")
)

type Helper struct {
	option *Option
	*viper.Viper
}

type Option struct {
	ConfigPath []string
	FileName   string
}

func get(prefix string) (*Helper, error) {
	if v, ok := store.Load(prefix); !ok {
		return nil, NotExists
	} else {
		val, _ := v.(*Helper)
		return val, nil
	}
}

func NewViper(prefix string, option *Option) (*Helper, error) {
	var err error
	var h = &Helper{}

	if h, err := get(prefix); err == nil {
		return h, nil
	}

	h = &Helper{
		option: option,
	}
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
	store.Store(prefix, h)
	return h, nil
}
