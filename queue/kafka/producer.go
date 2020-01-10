package kafka

import (
	"context"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"sync"
	"time"
)

type ProducerOptions = kafka.WriterConfig
type Message = kafka.Message

type ProducerHelper struct {
	options *ProducerOptions
	*kafka.Writer
}

var (
	ErrNotExists = errors.New("kafka not exists")
	ErrPushData  = errors.New("push data length is zero")
	producerOnce sync.Once
	Producer     *ProducerHelper
)

func balancer(b string) kafka.Balancer {
	var tmp kafka.Balancer
	switch b {
	case "roundrobin":
		tmp = &kafka.RoundRobin{}
	case "leastbytes":
		tmp = &kafka.LeastBytes{}
	case "hash":
		tmp = &kafka.Hash{}
	default:
		tmp = &kafka.RoundRobin{}
	}
	return tmp
}

func loadFromConfiger(prefix string, c config.ConfigerSlice) *ProducerOptions {
	p := prefix + "."
	return &ProducerOptions{
		Brokers:       c.GetStringSlice(utils.CombineString(p, "brokers")),
		Topic:         c.GetString(utils.CombineString(p, "topic")),
		MaxAttempts:   c.GetInt(utils.CombineString(p, "maxattempts")),
		QueueCapacity: c.GetInt(utils.CombineString(p, "queuecapacity")),
		BatchTimeout:  c.GetDuration(utils.CombineString(p, "batchtimeout")) * time.Millisecond,
		Async:         c.GetBool(utils.CombineString(p, "async")),
		RequiredAcks:  c.GetInt(utils.CombineString(p, "ack")),
		Balancer:      balancer(c.GetString(utils.CombineString(p, "balancer"))),
	}
}

func SampleProducerOptions(prefix string, c config.ConfigerSlice) *ProducerOptions {
	op := loadFromConfiger(prefix, c)

	if c.GetBool(utils.CombineString(prefix, ".", "compression")) {
		op.CompressionCodec = snappy.NewCompressionCodec()
	}
	return op
}

func NewProducer(options *ProducerOptions) *ProducerHelper {
	h := &ProducerHelper{}
	h.options = options
	h.Writer = kafka.NewWriter(*options)
	return h
}

type PushMessage struct {
	Key   []byte
	Value []byte
}

func (s *ProducerHelper) Push(data ...PushMessage) error {
	_, err := s.PushTimeout(0, data...)
	return err
}

func (s *ProducerHelper) PushTimeout(timeout time.Duration, data ...PushMessage) (context.CancelFunc, error) {
	if len(data) == 0 {
		return nil, ErrPushData
	}
	tmp := make([]kafka.Message, len(data))
	i := 0
	for _, m := range data {
		tmp[i] = kafka.Message{Key: m.Key, Value: m.Value}
		i++
	}
	var ctx = context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}
	return cancel, s.Writer.WriteMessages(ctx, tmp...)
}

func DefaultProducer() *ProducerHelper {
	producerOnce.Do(func() {
		Producer = NewProducer(SampleProducerOptions("kafka-producer", viper.GetSingleton()))
	})
	return Producer
}
