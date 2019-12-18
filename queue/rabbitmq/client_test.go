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
		_, err := NewRabbitmq(createOptions())
		So(err, ShouldEqual, nil)
	})
}

func TestChannel_PushlishSample(t *testing.T) {
	Convey("test channel_publish", t, func() {
		helper, err := NewRabbitmq(createOptions())
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
		helper, err := NewRabbitmq(createOptions())
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
		helper, err := NewRabbitmq(createOptions())
		So(err, ShouldEqual, nil)
		//restart rabbitmq by manual
		time.Sleep(20e9)
		_, err = helper.Channel()
		So(err, ShouldEqual, nil)

	})
}

func TestNewWithRetry(t *testing.T) {
	type args struct {
		prefix   string
		option   *Options
		attempts int
		interval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Helper
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    *Helper
			wantErr bool
		}{name: "testRetry", args: args{"testRetry", createOptions(), 10, time.Second}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewWithRetry(tt.args.option, tt.args.attempts, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithRetry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
