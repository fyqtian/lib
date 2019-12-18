package rabbitmq

import (
	"fmt"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
	"net/url"
	"sync"
	"time"
)

type Options struct {
	Host              string
	Port              string
	User              string
	Passwd            string
	Vhost             string
	Heartbeat         time.Duration
	ConnectionTimeout time.Duration
	Locale            string
}

var (
	store                    = sync.Map{}
	NotExists                = errors.New("rabbitmq not exists")
	LostError                = errors.New("Connection has lost")
	defaultHeartbeat         = 10 * time.Second
	defaultConnectionTimeout = 30 * time.Second
	defaultLocale            = "en_US"
	defaultReconnectTime     = 3 * time.Second
)

func (s *Options) url() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", url.QueryEscape(s.User), url.QueryEscape(s.Passwd), s.Host, s.Port)
}

func (s *Options) createConfig() *amqp.Config {
	config := &amqp.Config{
		Heartbeat: defaultHeartbeat,
		Locale:    defaultLocale,
	}
	if s.Vhost != "" {
		config.Vhost = s.Vhost
	}
	if s.Heartbeat != 0 {
		config.Heartbeat = s.Heartbeat
	}
	if s.ConnectionTimeout != 0 {
		config.Dial = amqp.DefaultDial(s.ConnectionTimeout)
	} else {
		config.Dial = amqp.DefaultDial(defaultConnectionTimeout)
	}
	if s.Locale != "" {
		config.Locale = s.Locale
	}
	return config
}

type Helper struct {
	options *Options
	conn    *amqp.Connection
}

func (s *Helper) connect() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	if conn, err = amqp.DialConfig(s.options.url(), *s.options.createConfig()); err != nil {
		return nil, err
	}
	return conn, nil
}

//todo
//if someone use unnamed queue after reconnect ,it will cause error what not found queue
func (s *Helper) Channel() (*Channel, error) {
	if ch, err := s.conn.Channel(); err != nil {
		return nil, err
	} else {
		tmp := &Channel{ch}
		//todo
		//is there exist error before register notify
		go func() {
			for {
				err, ok := <-ch.NotifyClose(make(chan *amqp.Error))
				if !ok {
					break
				}
				log.Println(err)
				for {
					// wait 3s for connection reconnect
					time.Sleep(defaultReconnectTime)
					var err error
					ch, err = s.conn.Channel()
					if err == nil {
						tmp.channel = ch
						break
					}
				}
			}

		}()
		return tmp, nil
	}
}

func (s *Helper) listen() {
	for {
		err, ok := <-s.conn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			break
		}
		log.Println(err)

		for {
			time.Sleep(defaultReconnectTime)
			conn, err := s.connect()
			if err == nil {
				s.conn = conn
				break
			}

		}
	}
}

func NewRabbitmq(prefix string, op *Options) (*Helper, error) {
	var (
		h    = &Helper{}
		err  error
		conn *amqp.Connection
	)
	if s, err := Get(prefix); err == nil {
		return s, nil
	}
	h.options = op

	if conn, err = h.connect(); err != nil {
		return nil, err
	}
	h.conn = conn
	go h.listen()
	store.Store(prefix, h)
	return h, nil
}

func NewWithRetry(prefix string, option *Options, attempts int, interval time.Duration) (*Helper, error) {
	var (
		h   *Helper
		err error
	)
	utils.Retry(func() error {
		if h, err = NewRabbitmq(prefix, option); err != nil {
			return err
		}
		return nil
	}, attempts, interval)
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