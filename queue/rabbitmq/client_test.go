package rabbitmq

import (
	"testing"
	"time"
)
import . "github.com/smartystreets/goconvey/convey"

func createOptions() *Options {
	return &Options{
		Host:              "ubuntuVM",
		Port:              "5672",
		User:              "guest",
		Passwd:            "guest",
		Vhost:             "test",
		Heartbeat:         0,
		ConnectionTimeout: 0,
		Locale:            "en_US",
	}
}

func TestNewRabbitmq(t *testing.T) {
	Convey("test NewRabbitmq", t, func() {
		_, err := NewRabbitmq("mq", createOptions())
		So(err, ShouldEqual, nil)
	})
}

func TestChannel_PushlishSample(t *testing.T) {
	Convey("test channel_publish", t, func() {
		helper, err := NewRabbitmq("mq", createOptions())
		So(err, ShouldEqual, nil)
		ch, err := helper.Channel()
		So(err, ShouldEqual, nil)

		q, err := ch.QueueSample("abc", false)
		So(err, ShouldEqual, nil)

		content := "abc"
		err = ch.PushlishSample("", q.Name, []byte(content))
		So(err, ShouldEqual, nil)
	})

}

func TestChannel_Consumer(t *testing.T) {
	Convey("test channel_consumer", t, func() {
		helper, err := NewRabbitmq("mq", createOptions())
		So(err, ShouldEqual, nil)

		ch, err := helper.Channel()
		So(err, ShouldEqual, nil)

		q, err := ch.QueueSample("abc", false)

		So(err, ShouldEqual, nil)
		content := "abc"

		msgs, err := ch.ConusmerSample(q.Name, true)
		So(err, ShouldEqual, nil)

		first := <-msgs
		So(string(first.Body), ShouldEqual, content)
	})
}

func TestHelper_listen(t *testing.T) {
	Convey("test reconnect", t, func() {
		helper, err := NewRabbitmq("mq", createOptions())
		So(err, ShouldEqual, nil)
		//restart rabbitmq by manual
		time.Sleep(20e9)
		_, err = helper.Channel()
		So(err, ShouldEqual, nil)

	})
}
