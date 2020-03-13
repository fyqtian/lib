package viper

import (
	"errors"
	"github.com/fyqtian/lib/utils"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	ErrNotExists = errors.New("config not exists")
	V            *Helper
)

type Helper struct {
	options *Options
	*viper.Viper
	//read env first
	//it will not read again from config when read from env first
	readFromEnv bool
}

type Options struct {
	ConfigPath []string
	FileName   string
}

func currentDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func (s *Options) SetPaths(p ...string) {
	var tmp []string
	for _, v := range p {
		if filepath.IsAbs(v) {
			tmp = append(tmp, v)
		} else {
			abs, _ := filepath.Abs(v)
			tmp = append(tmp, filepath.Join(utils.ExecPath(), v), abs)
		}
	}
	s.ConfigPath = tmp
}

func DefaultOptions() *Options {
	//for unit test
	// lib/configs/
	unitPath := filepath.Dir(filepath.Dir(currentDir())) + "/configs"
	op := &Options{
		FileName: os.Getenv("RUN_TIME"),
	}
	op.SetPaths("config", "configs", ".", unitPath)
	return op
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

func (s *Helper) ReadFromEnv() {
	s.readFromEnv = true
}

func (s *Helper) ReadFromConfig() {
	s.readFromEnv = false
}

func (s *Helper) bind(key string) {
	if s.readFromEnv {
		s.Viper.BindEnv(key, key)
	}
}

func (s *Helper) GetString(key string) string {
	s.bind(key)
	return s.Viper.GetString(key)
}

func (s *Helper) GetInt(key string) int {
	s.bind(key)
	return s.Viper.GetInt(key)
}

func (s *Helper) GetBool(key string) bool {
	s.bind(key)
	return s.Viper.GetBool(key)
}

func (s *Helper) GetDuration(key string) time.Duration {
	s.bind(key)
	return s.Viper.GetDuration(key)
}

func (s *Helper) GetFloat64(key string) float64 {
	s.bind(key)
	return s.Viper.GetFloat64(key)
}
