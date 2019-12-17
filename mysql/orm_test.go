package mysql

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func createOption() *Options {
	return &Options{
		Host:     "ubuntuVM",
		User:     "root",
		Passwd:   "123456",
		Name:     "test",
		Charset:  "utf8",
		Port:     "3306",
		Location: "Local",
		MaxConn:  20,
		IdelConn: 10,
		Debug:    true,
	}
}
func TestNewOrm(t *testing.T) {
	op := createOption()
	Convey("test NewOrm", t, func() {
		if h, err := NewOrm("test", op); err != nil {
			t.Fatal(err)
			return
		} else {
			hInstance, _ := NewOrm("test", op)
			So(h, ShouldEqual, hInstance)

			hOther, _ := NewOrm("test-other", op)
			So(h, ShouldNotEqual, hOther)
		}
	})
}

func TestNewWithRetry(t *testing.T) {
	Convey("test NewWithRetry", t, func() {
		op := createOption()
		_, err := NewWithRetry("test", op, 0, 10*time.Second)
		So(err, ShouldEqual, nil)
	})
}
