package mysql

import (
	"github.com/fyqtian/lib/config/viper"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func createOption() *Options {
	return SampleOptions("db", viper.GetSingleton())
}
func TestNewOrm(t *testing.T) {
	op := createOption()
	Convey("test NewOrm", t, func() {
		_, err := NewOrm(op)
		So(err, ShouldEqual, nil)
	})
}

func TestNewWithRetry(t *testing.T) {
	Convey("test NewWithRetry", t, func() {
		op := createOption()
		_, err := NewWithRetry(op, 0, 10*time.Second)
		So(err, ShouldEqual, nil)
	})
}

func TestConvenienceOrm(t *testing.T) {
	Convey("test  ConvenienceOrm", t, func() {
		tmp := ConvenienceOrm("db")
		So(tmp, ShouldNotEqual, nil)
	})

}
