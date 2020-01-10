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
	//if op.MinBytes == 0 {
	//	op.MinBytes = defaultMinbytes
	//}
	//if op.MaxBytes == 0 {
	//	op.MaxBytes = defaultMaxBytes
	//}

	return op
}

func NewConsumer(options *ConsumerOptions) *ConsumerHelper {
	h := &ConsumerHelper{}
	h.options = options
	h.Reader = kafka.NewReader(*options)
	return h
}

type ReaddCallback func(kafka.Message, error)
type FetchCallback func(kafka.Message, error) error

func (s *ConsumerHelper) Read(f ReaddCallback, async bool) error {
	return s.ReadTimeOut(f, async, 0)
}

func (s *ConsumerHelper) ReadTimeOut(f ReaddCallback, async bool, timeout time.Duration) error {
	ctx := context.Background()
	for {
		if timeout != 0 {
			//ignore cancel func
			ctx, _ = context.WithTimeout(ctx, timeout)
		}
		m, err := s.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if async {
			go func() {
				defer recover()
				f(m, err)
			}()
		} else {
			f(m, err)
		}
	}
}

func (s *ConsumerHelper) Fetch(f FetchCallback) error {
	return s.FetchTimeout(f, 0)
}

//it will
func (s *ConsumerHelper) FetchTimeout(f FetchCallback, timeout time.Duration) error {
	ctx := context.Background()
	for {
		if timeout != 0 {
			ctx, _ = context.WithTimeout(ctx, timeout)
		}
		m, err := s.Reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		rsErr := f(m, err)
		if rsErr == nil {
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
