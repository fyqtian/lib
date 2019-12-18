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
		DbName:   "test",
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
