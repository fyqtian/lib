package skafka

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
	"time"
)
import . "github.com/smartystreets/goconvey/convey"

func TestDefaultProducer(t *testing.T) {
	Convey("test default producer", t, func() {
		So(DefaultAsyncProducer(), ShouldNotEqual, nil)
	})
}

func TestProducer_Push(t *testing.T) {
	Convey("test producer push", t, func() {
		So(DefaultAsyncProducer(), ShouldNotEqual, nil)
		Debug()
		i := 0
		for {
			i++
			m := NewMessage([]byte("abs"), []byte(uuid.NewV4().String()), "test")
			s, err := DefaultAsyncProducer().PushMessage(m)
			fmt.Printf("%#v\n", s)
			fmt.Println(i, err, time.Now())
			time.Sleep(5e9)
		}
	})
}

func TestDefaultSyncProducer(t *testing.T) {
	Convey("test default producer", t, func() {
		So(DefaultSyncProducer(), ShouldNotEqual, nil)
		m := NewMessage(nil, []byte("99999"), "test")
		fmt.Println(DefaultSyncProducer().SendMessage(m))
	})
}
