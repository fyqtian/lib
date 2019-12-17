package viper

import (
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestNewViper(t *testing.T) {
	Convey("test newviper", t, func() {
		path, _ := filepath.Abs(".")
		op := &Options{[]string{path}, "test-config"}
		if config, err := NewViper("test", op); err != nil {
			t.Fatal(err)
			return
		} else {
			So(config.GetString("info.name"), ShouldEqual, "viper")
			So(config.GetInt("info.age"), ShouldEqual, 16)
			So(config.GetBool("info.isMan"), ShouldEqual, true)

			configInstance, _ := NewViper("test", op)
			So(config, ShouldEqual, configInstance)

			configOther, _ := NewViper("test-other", op)
			So(config, ShouldNotEqual, configOther)

		}
	})
}
