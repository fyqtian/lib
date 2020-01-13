package skafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/fyqtian/lib/config"
	"github.com/fyqtian/lib/config/viper"
	"github.com/fyqtian/lib/utils"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
)

var (
	singleOnce     sync.Once
	singleConsumer sarama.PartitionConsumer
	groupOnce      sync.Once
	groupConsumer  *Consumer
)

type ConsumerOptions struct {
	//kafka version
	Version     string
	Topics      []string
	Brokers     []string
	GroupId     string
	ClientId    string
	Assignor    string
	Debug       bool
	Compression bool
	Config      *sarama.Config
}

func (s *ConsumerOptions) balancer() sarama.BalanceStrategy {
	var tmp sarama.BalanceStrategy
	switch s.Assignor {
	case sarama.StickyBalanceStrategyName:
		tmp = sarama.BalanceStrategySticky
	case sarama.RoundRobinBalanceStrategyName:
		tmp = sarama.BalanceStrategyRoundRobin
	case sarama.RangeBalanceStrategyName:
		tmp = sarama.BalanceStrategyRange
	default:
		tmp = sarama.BalanceStrategyRoundRobin
	}
	return tmp
}

func loadFromConfig(prefix string, c config.ConfigerSlice) *ConsumerOptions {
	p := prefix + "."
	return &ConsumerOptions{
		Version:  c.GetString(utils.CombineString(p, "version")),
		Topics:   c.GetStringSlice(utils.CombineString(p, "topics")),
		Brokers:  c.GetStringSlice(utils.CombineString(p, "brokers")),
		GroupId:  c.GetString(utils.CombineString(p, "groupid")),
		ClientId: c.GetString(utils.CombineString(p, "clientid")),
		Assignor: c.GetString(utils.CombineString(p, "assignor")),
		Debug:    c.GetBool(utils.CombineString(p, "debug")),
	}
}

var defaultKafkaVersion = sarama.MaxVersion

func SampleOptions(prefix string, c config.ConfigerSlice) *ConsumerOptions {
	op := loadFromConfig(prefix, c)
	config := sarama.NewConfig()
	config.Version = defaultKafkaVersion
	if op.Version != "" {
		version, err := sarama.ParseKafkaVersion(op.Version)
		if err == nil {
			config.Version = version
		}
	}
	config.Consumer.Group.Rebalance.Strategy = op.balancer()
	if op.ClientId != "" {
		config.ClientID = op.ClientId
	}
	op.Config = config
	return op
}

type CallbackConsumerMessage func(*sarama.ConsumerMessage) error

type Consumer struct {
	ready    chan bool
	options  *ConsumerOptions
	client   sarama.ConsumerGroup
	callback CallbackConsumerMessage
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (s *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(s.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (s *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (s *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29

	for message := range claim.Messages() {
		if s.callback != nil {
			if err := s.callback(message); err == nil {
				session.MarkMessage(message, "")
			}
		} else {
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func (s *Consumer) SetCallback(f CallbackConsumerMessage) {
	s.callback = f
}

func (s *Consumer) Consumer(handler sarama.ConsumerGroupHandler) error {
	if handler == nil {
		handler = s
	}
	if err := s.client.Consume(context.Background(), s.options.Topics, handler); err != nil {
		return errors.WithMessage(err, "consumer topics")
	}
	return nil
}

func Debug() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}

func NewConsumerGroup(options *ConsumerOptions) (*Consumer, error) {
	consumer := &Consumer{
		ready:   make(chan bool),
		options: options,
	}

	client, err := sarama.NewConsumerGroup(options.Brokers, options.GroupId, options.Config)
	if err != nil {
		return nil, errors.WithMessage(err, "new consumerGroup")
	}
	consumer.client = client
	if options.Debug {
		Debug()
	}
	return consumer, nil
}

func DefaultConsumerGroup() *Consumer {
	groupOnce.Do(func() {
		tmp, err := NewConsumerGroup(SampleOptions("kafka-consumer", viper.GetSingleton()))
		if err != nil {
			panic(err)
		}
		groupConsumer = tmp
	})
	return groupConsumer
}

func NewConsumer(brokers []string, topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "new consumer")
	}
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, errors.WithMessage(err, "consumer partition")
	}
	return partitionConsumer, err
}

func DefaultConsumer() sarama.PartitionConsumer {
	singleOnce.Do(func() {
		c := viper.GetSingleton()
		p := "kafka-consumer."
		tmp, _ := NewConsumer(
			c.GetStringSlice(p+"brokers"),
			c.GetStringSlice(p + "topics")[0],
			c.GetInt32(p+"partition"),
			-1,
		)
		singleConsumer = tmp
	})
	return singleConsumer
}
