package viper

import (
	. "github.com/smartystreets/goconvey/convey"
	"path/filepath"
	"testing"
)

func TestNewViper(t *testing.T) {
	Convey("test newviper", t, func() {
		path, _ := filepath.Abs(".")
		op := &Options{[]string{path}, "config"}
		if config, err := NewViper(op); err != nil {
			t.Fatal(err)
			return
		} else {
			So(config.GetString("info.name"), ShouldEqual, "viper")
			So(config.GetInt("info.age"), ShouldEqual, 16)
			So(config.GetBool("info.isMan"), ShouldEqual, true)
		}
	})
}

func TestGetSingleton(t *testing.T) {
	Convey("test newviper", t, func() {
		So(GetSingleton(), ShouldNotEqual, nil)
		So(GetSingleton().GetString("info.name"), ShouldEqual, "viper")
	})
}
