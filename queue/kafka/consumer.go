package kafka

import (
	"context"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

var (
	Consumer     *ConsumerHelper
	consumerOnce sync.Once
)

type ConsumerOptions = kafka.ReaderConfig
type ConsumerHelper struct {
	options *ConsumerOptions
	*kafka.Reader
}

const (
	defaultGroupID  = "go-kafka-client"
	defaultMinbytes = 1e3
	defaultMaxBytes = 100e6
)

//both set partation>0 and groupid will cause panic
func loadConsumerFromConfiger(prefix string, c config.ConfigerSlice) *ConsumerOptions {
	p := prefix + "."
	return &ConsumerOptions{
		Brokers:        c.GetStringSlice(utils.CombineString(p, "brokers")),
		GroupID:        c.GetString(utils.CombineString(p, "groupid")),
		Topic:          c.GetString(utils.CombineString(p, "topic")),
		Partition:      c.GetInt(utils.CombineString(p, "partition")),
		MinBytes:       c.GetInt(utils.CombineString(p, "minbytes")),
		MaxBytes:       c.GetInt(utils.CombineString(p, "maxbytes")),
		CommitInterval: c.GetDuration(utils.CombineString(p, "commitinterval")) * time.Millisecond,
	}
}

func SampleConsumerOptions(prefix string, c config.ConfigerSlice) *ConsumerOptions {
	op := loadConsumerFromConfiger(prefix, c)
	if op.MinBytes == 0 {
		op.MinBytes = defaultMinbytes
	}
	if op.MaxBytes == 0 {
		op.MaxBytes = defaultMaxBytes
	}
	return op
}

func NewConsumer(options *ConsumerOptions) *ConsumerHelper {
	h := &ConsumerHelper{}
	h.options = options
	h.Reader = kafka.NewReader(*options)
	return h
}

type MessageCallback func(kafka.Message, error) error

func (s *ConsumerHelper) Read(f MessageCallback) {
	ctx := context.Background()
	for {
		m, err := s.ReadMessage(ctx)
		f(m, err)
	}
}

func (s *ConsumerHelper) Fetch(f MessageCallback) {
	ctx := context.Background()
	for {
		m, err := s.Reader.FetchMessage(ctx)
		rsErr := f(m, err)
		if err != nil && rsErr != nil {
			s.Reader.CommitMessages(ctx, m)
		}
	}
}

func DefaultConsumer() *ConsumerHelper {
	consumerOnce.Do(func() {
		Consumer = NewConsumer(SampleConsumerOptions("kafka-consumer", viper.GetSingleton()))
	})
	return Consumer
}
