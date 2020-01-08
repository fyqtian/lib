package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fyqtian/lib/config/viper"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var topic = "test"
var testMessage = "test-testMessage"

func createOption() *Options {
	return SampleOptions("emq", viper.GetSingleton())
}

func TestSampleOptions(t *testing.T) {
	Convey("test SampleOptions", t, func() {
		op := SampleOptions("emq", viper.GetSingleton())
		So(op.CleanSession, ShouldEqual, false)
	})
}

func TestNewMqtt(t *testing.T) {
	Convey("test NewMqtt", t, func() {
		_, err := NewMqtt(createOption())
		So(err, ShouldEqual, nil)
	})
}

func TestNewWithRetry(t *testing.T) {
	Convey("test NewWithRetry", t, func() {
		var err error
		_, err = NewWithRetry(createOption(), 0, 10*time.Second)
		So(err, ShouldEqual, nil)
	})
}

func TestHelper_Sub(t *testing.T) {
	Convey("test mqtt sub", t, func() {
		c, err := NewMqtt(createOption())
		So(err, ShouldEqual, nil)
		var receive []byte
		var callback = func(client mqtt.Client, body mqtt.Message) {
			receive = body.Payload()
		}

		err = c.Sub(topic, 0, callback)
		So(err, ShouldEqual, nil)
		time.Sleep(2e9)

		err = c.PubSimple(topic, testMessage)
		So(err, ShouldEqual, nil)
		time.Sleep(2e9)

		So(string(receive), ShouldEqual, testMessage)
	})
}
func TestHelper_PubSample(t *testing.T) {
	Convey("test mqtt pub", t, func() {
		c, err := NewMqtt(createOption())
		So(err, ShouldEqual, nil)
		err = c.PubSimple(topic, testMessage)
		So(err, ShouldEqual, nil)
	})
}

func TestHelper_Unsubscribe(t *testing.T) {
	Convey("test mqtt Unsubscribe", t, func() {
		c, err := NewMqtt(createOption())
		So(err, ShouldEqual, nil)
		err = c.Unsubscribe("test1", "test2")
		So(err, ShouldEqual, nil)
	})
}

func TestHelper_SubMultiple(t *testing.T) {
	Convey("test mqtt SubMultiple", t, func() {
		c, err := NewMqtt(createOption())
		So(err, ShouldEqual, nil)
		err = c.SubMultiple(map[string]byte{"test1": 0, "test2": 0}, func(client mqtt.Client, message mqtt.Message) {
			fmt.Println((message.Topic()), string(message.Payload()))
		})
		So(err, ShouldEqual, nil)
	})
}

func TestHelper_SubSimple(t *testing.T) {
	Convey("test mqtt subsimple", t, func() {
		//c, err := NewMqtt(createOption())
		//So(err, ShouldEqual, nil)
		//ch, err := c.SubSimple("test5")
		////for val := range ch {
		////	fmt.Print(string(val), ok)
		////}
		//So(err, ShouldEqual, nil)
		//So(ch, ShouldNotEqual, nil)
	})
}

func TestHelper_ListenLostConnection(t *testing.T) {
	Convey("test ListenLostConnection", t, func() {
		//c, err := NewMqtt(createOption())
		//So(err, ShouldEqual, nil)
		//for {
		//	if val, ok := <-c.lostConnectionNotifyChan; ok {
		//		fmt.Println(val, ok)
		//	}
		//}
	})
}
