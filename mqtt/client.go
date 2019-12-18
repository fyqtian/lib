package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
	"time"
)

type Options = MQTT.ClientOptions

type Helper struct {
	client     MQTT.Client
	options    *Options
	topicStore sync.Map
}
type topicInfo struct {
	qos       byte
	casllback MQTT.MessageHandler
}

var (
	store     = sync.Map{}
	NotExists = errors.New("mqtt not exists")
	LostError = errors.New("Connection has lost")
)

func (s *Helper) setOption(options *Options) {
	if options.OnConnect == nil {
		options.SetOnConnectHandler(s.onConnectHandler)
	}
	s.options = options
}

func (s *Helper) connect() error {
	client := MQTT.NewClient(s.options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	s.client = client
	return nil
}

func (s *Helper) onConnectHandler(client MQTT.Client) {
	//todo
	//重连了 自动订阅 可以改成按需 先这样处理了
	s.topicStore.Range(func(key, value interface{}) bool {
		val := value.(*topicInfo)
		keyStr := key.(string)
		s.Sub(keyStr, val.qos, val.casllback)
		return true
	})
}
func (s *Helper) Pub(topic string, qos byte, retained bool, payload interface{}) error {
	if !s.client.IsConnectionOpen() {
		return LostError
	}
	if token := s.client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Helper) PubSample(topic string, payload interface{}) error {
	return s.Pub(topic, 0, false, payload)
}

//todo
//无法处理通配订阅
// /sub/+ /sub/# 订阅不可用只能自己实现messageHandler
// /sub/a 可以
func (s *Helper) Sub(topic string, qos byte, callback MQTT.MessageHandler) error {
	if !s.client.IsConnectionOpen() {
		return LostError
	}
	if token := s.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	s.topicStore.Store(topic, &topicInfo{qos, callback})
	//断链后被重新订阅，chan不同
	return nil
}

func (s *Helper) GetOptionsReader() MQTT.ClientOptionsReader {
	return s.GetClient().OptionsReader()
}
func (s *Helper) GetClient() MQTT.Client {
	return s.client
}

//func Get(prefix string) (*Helper, error) {
//	if v, ok := store.Load(prefix); !ok {
//		return nil, NotExists
//	} else {
//		val, _ := v.(*Helper)
//		return val, nil
//	}
//}

func NewMqtt(option *Options) (*Helper, error) {
	var (
		h = &Helper{}
	)
	//if s, err := Get(prefix); err == nil {
	//	return s, err
	//}
	h.setOption(option)
	if err := h.connect(); err != nil {
		return nil, err
	}
	//store.Store(prefix, h)
	return h, nil
}

//todo
func NewWithRetry(option *MQTT.ClientOptions, attempt int, interval time.Duration) (*Helper, error) {
	var (
		h   *Helper
		err error
	)
	utils.Retry(func() error {
		if h, err = NewMqtt(option); err != nil {
			return err
		}
		return nil
	}, attempt, interval)
	return h, err
}

func NewOptions() *Options {
	return MQTT.NewClientOptions()
}

func Debug() {
	MQTT.DEBUG = log.New(os.Stdout, "", 0)
	MQTT.ERROR = log.New(os.Stdout, "", 0)
}
