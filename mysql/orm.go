package mysql

import (
	"errors"
	"fmt"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
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
	Port     string
	User     string
	Passwd   string
	DbName   string
	Charset  string
	Location string
	MaxConn  int
	IdelConn int
	Debug    bool
}

const (
	defaultLocation = "local"
	defaultCharset  = "utf8"
	defaultMaxConn  = 20
	defaultIdelConn = 10
)

var (
	store        sync.Map
	ErrNotExists = errors.New("orm not exists")
	once         sync.Once
	Orm          *Helper
)

func loadFromConfiger(prefix string, c config.Configer) *Options {
	p := prefix + "."
	return &Options{
		Host:     c.GetString(utils.CombineString(p, "host")),
		User:     c.GetString(utils.CombineString(p, "user")),
		Passwd:   c.GetString(utils.CombineString(p, "passwd")),
		DbName:   c.GetString(utils.CombineString(p, "dbname")),
		Charset:  c.GetString(utils.CombineString(p, "charset")),
		Port:     c.GetString(utils.CombineString(p, "port")),
		Location: c.GetString(utils.CombineString(p, "location")),
		MaxConn:  c.GetInt(utils.CombineString(p, "maxconn")),
		IdelConn: c.GetInt(utils.CombineString(p, "idelconn")),
		Debug:    c.GetBool(utils.CombineString(p, "debug")),
	}
}

func SampleOptions(prefix string, c config.Configer) *Options {
	op := loadFromConfiger(prefix, c)
	if op.Location == "" {
		op.Location = defaultLocation
	}
	if op.Charset == "" {
		op.Charset = defaultCharset
	}
	if op.MaxConn == 0 {
		op.MaxConn = defaultMaxConn
	}

	if op.IdelConn == 0 {
		op.IdelConn = defaultIdelConn
	}
	return op
}

func DefaultOrm() *Helper {
	once.Do(func() {
		//var err error
		Orm, _ = NewWithRetry(SampleOptions("db", viper.GetSingleton()), 10, 5*time.Second)
		//if err != nil {
		//	panic(err)
		//}
	})
	return Orm
}

func Get(prefix string) (*Helper, error) {
	if v, ok := store.Load(prefix); !ok {
		return nil, ErrNotExists
	} else {
		val, _ := v.(*Helper)
		return val, nil
	}
}

func ConvenienceOrm(prefix string) *Helper {
	if s, err := Get(prefix); err == nil {
		return s
	}
	obj, err := NewWithRetry(SampleOptions(prefix, viper.GetSingleton()), 99999, 5*time.Second)
	if err != nil {
		store.Store(prefix, obj)
	}
	return obj
}

func NewOrm(op *Options) (*Helper, error) {
	var (
		h   = &Helper{}
		err error
	)
	//if s, err := Get(prefix); err == nil {
	//	return s, nil
	//}
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
	//store.Store(prefix, h)
	return h, err
}

func NewWithRetry(option *Options, attempts int, interval time.Duration) (*Helper, error) {
	var (
		h   *Helper
		err error
	)
	utils.Retry(func() error {
		if h, err = NewOrm(option); err != nil {
			return err
		}
		return nil
	}, attempts, interval)
	return h, err
}

func (s *Helper) combineDSN() string {
	if s.options.Location == "" {
		s.options.Location = ""
	}
	if s.options.Charset == "" {
		s.options.Charset = defaultCharset
	}
	if s.options.MaxConn == 0 {
		s.options.MaxConn = defaultMaxConn
	}
	if s.options.IdelConn == 0 {
		s.options.IdelConn = defaultIdelConn
	}

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
