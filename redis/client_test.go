package redis

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"testing"
	"time"
)

var toml = []byte(`
[redis]
addr = "ubuntuVM:6379"
auth = ""
prefix = "tob_"
index = 0
retry = 10
dialtimeout = 3
`)

func configer() {
	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer(toml))
}

func init() {
	configer()
}

func TestNewRedis(t *testing.T) {
	Convey("test NewRedis", t, func() {
		h, err := NewRedis(SampleOptions("redis", viper.GetViper()))
		So(err, ShouldEqual, nil)
		key := "test-key"
		value := "test-value"
		h.Set(key, value, 2*time.Second)
		time.Sleep(time.Second)
		tmpVal, _ := h.Get(key).Result()
		So(tmpVal, ShouldEqual, value)
	})
}
