package kafka

import (
	"errors"
	"fmt"
	"github.com/fyqtian/lib/config/viper"
	"github.com/segmentio/kafka-go"
	. "github.com/smartystreets/goconvey/convey"
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
		go func() {
			DefaultConsumer().Read(func(message kafka.Message, e error) error {
				fmt.Println(string(message.Value), e)
				return nil
			})
		}()
		time.Sleep(600e9)
	})
}

func TestConsumerHelper_Fetch(t *testing.T) {
	Convey("test consumer fetch", t, func() {
		go func() {
			DefaultConsumer().Fetch(func(message kafka.Message, e error) error {
				fmt.Println(string(message.Value), e, message.Partition, message.Offset)
				return errors.New("ttt")
			})
		}()

		time.Sleep(600e9)
	})
}
