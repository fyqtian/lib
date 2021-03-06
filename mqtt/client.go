package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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
	topicChan  map[string]chan MQTT.Message
	sync.Mutex
	lostConnectionNotifyChan chan error
}

type topicInfo struct {
	qos       byte
	casllback MQTT.MessageHandler
}

type Client interface {
	PubSimple(topic string, payload interface{}) error
	SubSimple(string) (<-chan MQTT.Message, error)
}

var (
	ErrNotExists   = errors.New("mqtt not exists")
	ErrLostConnect = errors.New("mqtt connection lost")
	once           sync.Once
	Mqtt           *Helper
)

func (s *Helper) setOption(options *Options) {
	if options.OnConnect == nil {
		options.SetOnConnectHandler(s.onConnectHandler)
	}
	options.SetConnectionLostHandler(s.onConnectionLostHandler)

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

//监听断链
func (s *Helper) ListenLostConnection() <-chan error {
	return s.lostConnectionNotifyChan
}

func (s *Helper) onConnectionLostHandler(c MQTT.Client, err error) {
	select {
	case s.lostConnectionNotifyChan <- err:
	default:
	}
}

//注册处理重链后 sub topic
func (s *Helper) onConnectHandler(client MQTT.Client) {
	//todo
	//重连了 自动订阅 可以改成按需 先这样处理了
	s.topicStore.Range(func(key, value interface{}) bool {
		switch keyValue := key.(type) {
		case *map[string]byte:
			s.SubMultiple(*keyValue, value.(MQTT.MessageHandler))
		case string:
			val := value.(*topicInfo)
			s.Sub(keyValue, val.qos, val.casllback)
		}
		return true
	})
}
func (s *Helper) Pub(topic string, qos byte, retained bool, payload interface{}) error {
	if !s.client.IsConnectionOpen() {
		return ErrLostConnect
	}
	if token := s.client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Helper) PubSimple(topic string, payload interface{}) error {
	return s.Pub(topic, 0, false, payload)
}

func (s *Helper) Sub(topic string, qos byte, callback MQTT.MessageHandler) error {
	if !s.client.IsConnectionOpen() {
		return ErrLostConnect
	}
	if token := s.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	//断链后被重新订阅
	s.topicStore.Store(topic, &topicInfo{qos, callback})
	return nil
}

func (s *Helper) SubSimple(topic string) (<-chan MQTT.Message, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.topicChan[topic] = make(chan MQTT.Message, 4)
	err := s.Sub(topic, 0, func(client MQTT.Client, message MQTT.Message) {
		s.topicChan[topic] <- message
	})
	return s.topicChan[topic], err
}

func (s *Helper) SubMultiple(topics map[string]byte, callback MQTT.MessageHandler) error {
	if !s.client.IsConnectionOpen() {
		return ErrLostConnect
	}
	if token := s.client.SubscribeMultiple(topics, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	//todo
	s.topicStore.Store(&topics, callback)
	return nil
}

func (s *Helper) Unsubscribe(topics ...string) error {
	if !s.client.IsConnectionOpen() {
		return ErrLostConnect
	}
	if token := s.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Helper) GetOptionsReader() MQTT.ClientOptionsReader {
	return s.GetClient().OptionsReader()
}
func (s *Helper) GetClient() MQTT.Client {
	return s.client
}

func (s *Helper) Disconnect(i uint) {
	s.client.Disconnect(i)
	for _, v := range s.topicChan {
		close(v)
	}
}

func NewMqtt(option *Options) (*Helper, error) {
	var (
		h = &Helper{}
	)
	h.setOption(option)
	if err := h.connect(); err != nil {
		return nil, err
	}
	h.topicChan = make(map[string]chan MQTT.Message, 8)
	h.lostConnectionNotifyChan = make(chan error)
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

func SampleOptions(prefix string, c config.Configer) *Options {
	p := prefix + "."
	op := NewOptions()

	op.AddBroker(c.GetString(utils.CombineString(p, "addr")))

	if t := c.GetString(utils.CombineString(p, "clientid")); t == "" {
		op.SetClientID(uuid.NewV4().String())
	} else {
		op.SetClientID(t)
	}
	op.SetUsername(c.GetString(utils.CombineString(p, "user")))
	op.SetPassword(c.GetString(utils.CombineString(p, "passwd")))

	if t := c.GetDuration(utils.CombineString(p, "keepalive")); t != 0 {
		op.SetKeepAlive(t * time.Second)
	}
	if t := c.GetDuration(utils.CombineString(p, "pingtimeout")); t != 0 {
		op.SetPingTimeout(t * time.Second)
	}
	//op.SetAutoReconnect(c.GetBool(utils.CombineString(p, "reconnect")))
	op.SetCleanSession(c.GetBool(utils.CombineString(p, "cleansession")))
	return op
}

func DefaultMqtt() *Helper {
	var err error
	once.Do(func() {
		Mqtt, err = NewWithRetry(SampleOptions("emq", viper.GetSingleton()), 10, 5*time.Second)
		if err != nil {
			panic(err)
		}
	})
	return Mqtt
}
