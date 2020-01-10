package kafka

import (
	"fmt"
	"github.com/fyqtian/lib/config/viper"
	"github.com/segmentio/kafka-go"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

func TestSampleConsumerOptions(t *testing.T) {
	Convey("test SampleConsumerOptions", t, func() {
		op := SampleConsumerOptions("kafka-consumer", viper.GetSingleton())
		//fmt.Printf("%#v", op)
		So(op.GroupID, ShouldEqual, "config-group-id")
	})
}

func TestDefaultConsumer(t *testing.T) {
	Convey("test DefaultConsumer", t, func() {
		So(DefaultConsumer(), ShouldNotEqual, nil)
	})
}

func TestConsumerHelper_Read(t *testing.T) {
	Convey("test consumer read", t, func() {
		DefaultConsumer().Read(func(message kafka.Message, e error) {
			log.Println(string(message.Value), e)
		}, false)
	})
}

func TestConsumerHelper_ReadTimeOut(t *testing.T) {
	Convey("test consumer fetch", t, func() {
		DefaultConsumer().ReadTimeOut(func(message kafka.Message, e error) {
			log.Println(string(message.Value), e, message.Partition, message.Offset)
		}, false, time.Second*20)

	})
}

func TestConsumerHelper_Fetch(t *testing.T) {
	Convey("test consumer fetch", t, func() {
		fmt.Println(123123)
		fmt.Printf("%#v\n", DefaultConsumer().options)
		fmt.Println(DefaultConsumer().Fetch(func(message kafka.Message, e error) error {
			fmt.Println(string(message.Value), e, message.Partition, message.Offset)
			return nil
		}), 5555)

	})
}
