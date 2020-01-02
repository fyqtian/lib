package rabbitmq

import (
	"fmt"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
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

const (
	defaultHeartbeat         = 10 * time.Second
	defaultConnectionTimeout = 30 * time.Second
	defaultHost              = "/"
	defaultLocale            = "en_US"
	defaultReconnectTime     = 3 * time.Second
)

var (
	ErrNotExists      = errors.New("rabbitmq not exists")
	ErrLostConnection = errors.New("Connection has lost")
	once              sync.Once
	MQ                *Helper
)

func (s *Options) url() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", url.QueryEscape(s.User), url.QueryEscape(s.Passwd), s.Host, s.Port)
}

func (s *Options) createConfig() amqp.Config {
	config := amqp.Config{
		Heartbeat: defaultHeartbeat,
		Locale:    defaultLocale,
	}
	if s.Vhost != "" {
		config.Vhost = s.Vhost
	} else {
		config.Vhost = defaultHost
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
	options                  *Options
	conn                     *amqp.Connection
	lostConnectionNotifyChan chan error
}

func (s *Helper) connect() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	if conn, err = amqp.DialConfig(s.options.url(), s.options.createConfig()); err != nil {
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
				select {
				case s.lostConnectionNotifyChan <- err:
				default:
				}
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
		//todo
		select {
		case s.lostConnectionNotifyChan <- err:
		default:
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

func (s *Helper) ListenLostConnection() <-chan error {
	return s.lostConnectionNotifyChan
}

func NewRabbitmq(op *Options) (*Helper, error) {
	var (
		h    = &Helper{}
		err  error
		conn *amqp.Connection
	)
	h.options = op
	if conn, err = h.connect(); err != nil {
		return nil, err
	}
	h.conn = conn
	h.lostConnectionNotifyChan = make(chan error)
	go h.listen()
	return h, nil
}

func NewWithRetry(option *Options, attempts int, interval time.Duration) (*Helper, error) {
	var (
		h   *Helper
		err error
	)
	utils.Retry(func() error {
		if h, err = NewRabbitmq(option); err != nil {
			return err
		}
		return nil
	}, attempts, interval)
	return h, err
}

func SampleOptions(prefix string, c config.Configer) *Options {
	p := prefix + "."
	op := &Options{
		Host:              c.GetString(utils.CombineString(p, "host")),
		Port:              c.GetString(utils.CombineString(p, "port")),
		User:              c.GetString(utils.CombineString(p, "user")),
		Passwd:            c.GetString(utils.CombineString(p, "passwd")),
		Vhost:             c.GetString(utils.CombineString(p, "vhost")),
		Heartbeat:         c.GetDuration(utils.CombineString(p, "hearbeat")) * time.Second,
		ConnectionTimeout: c.GetDuration(utils.CombineString(p, "connectiontimeout")) * time.Second,
		Locale:            c.GetString(utils.CombineString(p, "locale")),
	}
	return op
}

func DefaultMQ() *Helper {
	once.Do(func() {
		//var err error
		MQ, _ = NewWithRetry(SampleOptions("mq", viper.GetSingleton()), 10, 5*time.Second)
		//if err != nil {
		//	panic(err)
		//}
	})
	return MQ
}
