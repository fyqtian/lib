package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var topic = "test"
var testMessage = "test-testMessage"

func createOption() *Options {
	op := NewOptions()
	op.AddBroker("tcp://ubuntuVM:1883")
	op.SetClientID(uuid.NewV4().String())
	op.SetUsername("pushCore")
	op.SetPassword("pushCore")
	return op
}

func TestNewMqtt(t *testing.T) {
	Convey("test NewMqtt", t, func() {
		op := createOption()
		_, err := NewMqtt(op)
		So(err, ShouldEqual, nil)
	})
}

func TestNewWithRetry(t *testing.T) {
	Convey("test NewWithRetry", t, func() {
		op := createOption()
		var err error
		_, err = NewWithRetry(op, 0, 10*time.Second)
		So(err, ShouldEqual, nil)

	})
}

func TestHelper_Sub(t *testing.T) {
	Convey("test mqtt sub", t, func() {
		op := createOption()
		if c, err := NewMqtt(op); err != nil {
			t.Fatal(err)
			return
		} else {
			var receive []byte
			var callback = func(client mqtt.Client, body mqtt.Message) {
				receive = body.Payload()
			}
			if err := c.Sub(topic, 0, callback); err != nil {
				t.Fatal(err)
			}

			time.Sleep(2e9)
			if err := c.PubSample(topic, testMessage); err != nil {
				t.Fatal(err)
			}
			time.Sleep(2e9)
			So(string(receive), ShouldEqual, testMessage)
		}
	})
}
func TestHelper_PubSample(t *testing.T) {
	Convey("test mqtt pub", t, func() {
		op := createOption()
		if c, err := NewMqtt(op); err != nil {
			t.Fatal(err)
			return
		} else {
			if err := c.PubSample(topic, testMessage); err != nil {
				t.Fatal(err)
			}
		}
	})
}
