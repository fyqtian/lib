package redis

import (
	"errors"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type Helper struct {
	*redis.Client
	options *Options
}

type Options = redis.Options

var (
	once          sync.Once
	Redis         *Helper
	SertNxErrFail = errors.New("setnx fail")
)

func SampleOptions(prefix string, c config.Configer) *Options {
	p := prefix + "."
	op := &redis.Options{
		Addr:       c.GetString(utils.CombineString(p, "addr")),
		Password:   c.GetString(utils.CombineString(p, "passwd")),
		DB:         c.GetInt(utils.CombineString(p, "index")),
		MaxRetries: c.GetInt(utils.CombineString(p, "retry")),
		//MinRetryBackoff: c.GetDuration(utils.CombineString(p, "minretrytime")),
		//MaxRetryBackoff: c.GetDuration(utils.CombineString(p, "maxretrytime")),
		DialTimeout:  c.GetDuration(utils.CombineString(p, "dialtimeout")) * time.Second,
		PoolSize:     c.GetInt(utils.CombineString(p, "poolsize")),
		MinIdleConns: c.GetInt(utils.CombineString(p, "minidelconn")),
	}
	return op
}

func NewRedis(options *Options) (*Helper, error) {

	c := redis.NewClient(options)
	if _, err := c.Ping().Result(); err != nil {
		return nil, err
	} else {
		h := &Helper{
			Client:  c,
			options: options,
		}
		return h, nil
	}
}

func DefaultRedis() *Helper {
	once.Do(func() {
		var err error
		Redis, err = NewRedis(SampleOptions("redis", viper.GetSingleton()))
		if err != nil {
			panic(err)
		}
	})
	return Redis
}

type SpinLocker struct {
	key        string
	token      string
	tryTimes   int
	expireTime time.Duration
	*Helper
}

//todo
//断线后不支持重入
func NewSpinLocker(key, token string, tryTimes int, expireTime time.Duration, c *Helper) *SpinLocker {
	return &SpinLocker{
		key:        key,
		token:      token,
		tryTimes:   tryTimes,
		expireTime: expireTime,
		Helper:     c,
	}
}

//todo
//还要处理下
func (s *SpinLocker) Lock() error {
	//先尝试拿锁
	val, err := s.Client.Get(s.key).Result()
	//当前没有锁
	if err == redis.Nil {
		ok, err := s.Client.SetNX(s.key, s.token, s.expireTime).Result()
		if err != nil {
			return err
		}
		if !ok {
			return SertNxErrFail
		}
	} else {
		//如果持有当前锁 刷新过期时间
		if val == s.token {
			_, err = s.Client.Set(s.key, s.token, s.expireTime).Result()
			return err
		} else {
			return err
		}
	}

	return nil
}
func (s *SpinLocker) Unlock() error {
	_, err := s.Client.Del(s.key).Result()
	if err != nil {
		return err
	}
	return nil
}
