package skafka

import (
	"github.com/Shopify/sarama"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	producerOnce   sync.Once
	producerHelper *Producer
)

const (
	defaultProduceTimeout = 5 * time.Second
)

type ProducerOptions struct {
	Brokers []string
	//-1,0,1,
	Ack    int
	config *sarama.Config
}

func producerLoadFromConfig(prefix string, c config.ConfigerSlice) *ProducerOptions {
	p := prefix + "."
	return &ProducerOptions{
		Brokers: c.GetStringSlice(utils.CombineString(p, "brokers")),
		Ack:     c.GetInt(utils.CombineString(p, "ack")),
	}
}

func SampleProducerOptions(prefix string, c config.ConfigerSlice) *ProducerOptions {
	op := producerLoadFromConfig(prefix, c)
	op.config = sarama.NewConfig()
	op.config.Producer.Timeout = defaultProduceTimeout
	op.config.Producer.RequiredAcks = sarama.RequiredAcks(op.Ack)
	op.config.Producer.Return.Successes = true
	return op
}

type Producer struct {
	options *ProducerOptions
	sarama.AsyncProducer
}

type Message struct {
	key   []byte
	value string
	topic string
}

func NewMessage(key, value []byte, topic string) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}
}

func (s *Producer) PushMessage(msg *sarama.ProducerMessage) (*sarama.ProducerMessage, *sarama.ProducerError) {
	s.AsyncProducer.Input() <- msg
	select {
	case err := <-s.AsyncProducer.Errors():
		return nil, err
	case m := <-s.AsyncProducer.Successes():
		return m, nil
	}
}

func (s *Producer) PushMessages(msgs ...*sarama.ProducerMessage) {
	for _, msg := range msgs {
		s.PushMessage(msg)
	}
}

func NewProducer(options *ProducerOptions) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer(options.Brokers, options.config)
	if err != nil {
		return nil, errors.WithMessage(err, "new producer")
	}
	return &Producer{options, producer}, nil
}

func DefaultProducer() *Producer {
	producerOnce.Do(func() {
		tmp, _ := NewProducer(SampleProducerOptions("kafka-producer", viper.GetSingleton()))
		producerHelper = tmp
	})
	return producerHelper
}
