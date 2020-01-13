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

type Message = sarama.ProducerMessage

var (
	asyncProducerOnce   sync.Once
	asyncProducerHelper *AsyncProducer
	syncProducerOnce    sync.Once
	syncProducerHelper  *SyncProducer
)

const (
	defaultProduceTimeout = 5 * time.Second
)

type ProducerOptions struct {
	Brokers []string
	//0 NONE
	//1 gzip
	//2 snappy
	//3 lz4
	//4 zstd
	CompressionType int
	//0  doesn't send any response, the TCP ACK is all you get.
	// -1  waits for only the local commit to succeed before responding.
	//1 waits for all in-sync replicas to commit before responding
	Ack      int
	ClientId string
	Config   *sarama.Config
}

func producerLoadFromConfig(prefix string, c config.ConfigerSlice) *ProducerOptions {
	p := prefix + "."
	return &ProducerOptions{
		Brokers:         c.GetStringSlice(utils.CombineString(p, "brokers")),
		Ack:             c.GetInt(utils.CombineString(p, "ack")),
		CompressionType: c.GetInt(utils.CombineString(p, "compressiontype")),
		ClientId:        c.GetString(utils.CombineString(p, "clientid")),
	}
}

func SampleProducerOptions(prefix string, c config.ConfigerSlice) *ProducerOptions {
	op := producerLoadFromConfig(prefix, c)
	op.Config = sarama.NewConfig()
	op.Config.Producer.Timeout = defaultProduceTimeout
	op.Config.Producer.RequiredAcks = sarama.RequiredAcks(op.Ack)
	op.Config.Producer.Return.Successes = true
	if op.CompressionType > 0 {
		op.Config.Producer.Compression = sarama.CompressionCodec(op.CompressionType)
	}
	if op.ClientId != "" {
		op.Config.ClientID = op.ClientId
	}

	return op
}

type AsyncProducer struct {
	options *ProducerOptions
	sarama.AsyncProducer
}

type SyncProducer struct {
	options *ProducerOptions
	sarama.SyncProducer
}

func NewMessage(key, value []byte, topic string) *Message {
	return &Message{
		Topic:     topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Timestamp: time.Now(),
	}
}

func (s *AsyncProducer) PushMessage(msg *Message) (*sarama.ProducerMessage, *sarama.ProducerError) {
	s.AsyncProducer.Input() <- msg
	select {
	case err := <-s.AsyncProducer.Errors():
		return nil, err
	case m := <-s.AsyncProducer.Successes():
		return m, nil
	}
}

func NewAsyncProducer(options *ProducerOptions) (*AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(options.Brokers, options.Config)
	if err != nil {
		return nil, errors.WithMessage(err, "new async producer")
	}
	return &AsyncProducer{options, producer}, nil
}

func DefaultAsyncProducer() *AsyncProducer {
	asyncProducerOnce.Do(func() {
		tmp, err := NewAsyncProducer(SampleProducerOptions("kafka-producer", viper.GetSingleton()))
		if err != nil {
			panic(err)
		}
		asyncProducerHelper = tmp
	})
	return asyncProducerHelper
}

func NewSyncProducer(options *ProducerOptions) (*SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(options.Brokers, options.Config)
	if err != nil {
		return nil, errors.WithMessage(err, "new sync producer")
	}
	return &SyncProducer{options, producer}, nil
}

func DefaultSyncProducer() *SyncProducer {
	syncProducerOnce.Do(func() {
		tmp, err := NewSyncProducer(SampleProducerOptions("kafka-producer", viper.GetSingleton()))
		if err != nil {
			panic(err)
		}
		syncProducerHelper = tmp
	})
	return syncProducerHelper
}
