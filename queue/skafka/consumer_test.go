package skafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/fyqtian/lib/config/viper"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSampleOptions(t *testing.T) {
	Convey("test sampleOptions", t, func() {

	})
}

func TestNewConsumerGroup(t *testing.T) {
	Convey("test newConsumerGroup", t, func() {
		op := SampleOptions("kafka-consumer", viper.GetSingleton())
		_, err := NewConsumerGroup(op)
		So(err, ShouldEqual, nil)
	})
}

func Test_loadFromConfig(t *testing.T) {
	Convey("test loadfromconfig", t, func() {
		op := loadFromConfig("kafka-consumer", viper.GetSingleton())
		So(op.Version, ShouldEqual, "2.2.1")
	})
}

func TestNewConsumer(t *testing.T) {
	Convey("test newConsumer", t, func() {

		Debug()
		So(DefaultConsumer(), ShouldNotEqual, nil)
		for msg := range DefaultConsumer().Messages() {
			fmt.Println(msg.Topic, msg.Offset, string(msg.Value))
		}
	})
}

func TestDefaultConsumerGroup(t *testing.T) {
	Convey("test default consumer group", t, func() {
		So(DefaultConsumerGroup(), ShouldNotEqual, nil)
		DefaultConsumerGroup().SetCallback(func(message *sarama.ConsumerMessage) error {
			fmt.Println(string(message.Value), message.Offset)
			return nil
		})
		fmt.Println(DefaultConsumerGroup().Consumer(nil))
	})
}
