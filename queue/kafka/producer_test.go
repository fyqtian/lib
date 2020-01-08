package kafka

import (
	"fmt"
	"github.com/fyqtian/lib/config/viper"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSampleProducerOptions(t *testing.T) {
	Convey("test sampleProducerOptions", t, func() {
		op := SampleProducerOptions("kafka-producer", viper.GetSingleton())

		So(op.Brokers, ShouldNotEqual, nil)
		So(op.CompressionCodec, ShouldNotEqual, nil)
		//So(op.Async, ShouldEqual, true)

	})
}

func TestDefaultProducer(t *testing.T) {
	Convey("test default producer", t, func() {
		So(DefaultProducer(), ShouldNotEqual, nil)
	})
}

func TestProducerHelper_Push(t *testing.T) {
	Convey("test default push", t, func() {

		for i := 10; i > 0; i-- {
			fmt.Println(i)
			err := DefaultProducer().Push(PushMessage{[]byte(fmt.Sprintf("i%d", i)), []byte("dddd")})
			So(err, ShouldEqual, nil)
		}
	})
}
