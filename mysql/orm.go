package mysql

import (
	"errors"
	"fmt"
	"github.com/fyqtian/lib/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

type Helper struct {
	*gorm.DB
	options *Options
	name    string
}

type Options struct {
	Host     string
	User     string
	Passwd   string
	DbName   string
	Charset  string
	Port     string
	Location string
	MaxConn  int
	IdelConn int
	Debug    bool
}

var (
	store     sync.Map
	NotExists = errors.New("zap not exists")
)

func NewOrm(prefix string, op *Options) (*Helper, error) {
	var (
		h   = &Helper{}
		err error
	)
	if s, err := Get(prefix); err == nil {
		return s, nil
	}
	h.options = op
	db, err := gorm.Open("mysql", h.combineDSN())
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	if op.Debug {
		db = db.Debug()
	}
	h.DB = db
	h.DB.DB().SetMaxOpenConns(op.MaxConn)
	h.DB.DB().SetMaxIdleConns(op.IdelConn)
	store.Store(prefix, h)
	return h, err
}

func Get(prefix string) (*Helper, error) {
	if v, ok := store.Load(prefix); !ok {
		return nil, NotExists
	} else {
		val, _ := v.(*Helper)
		return val, nil
	}
}

func NewWithRetry(prefix string, option *Options, attempts int, interval time.Duration) (*Helper, error) {
	var (
		h   *Helper
		err error
	)
	utils.Retry(func() error {
		if h, err = NewOrm(prefix, option); err != nil {
			return err
		}
		return nil
	}, attempts, interval)
	return h, err
}

func (s *Helper) combineDSN() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		s.options.User,
		s.options.Passwd,
		s.options.Host,
		s.options.Port,
		s.options.DbName,
		s.options.Charset,
		s.options.Location)
	return dsn
}
