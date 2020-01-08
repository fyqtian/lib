package viper

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
	"os"
	"path/filepath"
	"testing"
)

func TestNewViper(t *testing.T) {
	Convey("test newviper", t, func() {
		path, _ := filepath.Abs(".")
		op := &Options{[]string{path}, "config"}
		config, err := NewViper(op)
		So(err, ShouldEqual, nil)

		So(config.GetString("info.name"), ShouldEqual, "viper")
		So(config.GetInt("info.age"), ShouldEqual, 16)
		So(config.GetBool("info.isMan"), ShouldEqual, true)
	})
}

func TestGetSingleton(t *testing.T) {
	Convey("test newviper", t, func() {
		So(GetSingleton(), ShouldNotEqual, nil)
		So(GetSingleton().GetString("info.name"), ShouldEqual, "viper")
	})
}

func TestHelper_GetString(t *testing.T) {
	Convey("test get string", t, func() {
		var key = "info.name"
		var value = "envValue"

		v, _ := NewViper(DefaultOptions())
		os.Setenv(key, value)
		So(v.GetString(key), ShouldEqual, "viper")

		v.ReadFromEnv()
		So(v.GetString(key), ShouldEqual, value)

	})
}

func TestHelper_GetInt(t *testing.T) {
	Convey("test get string", t, func() {
		var key = "test.age"
		var age = "100"

		v, _ := NewViper(DefaultOptions())
		os.Setenv(key, age)

		So(v.GetInt(key), ShouldEqual, 22)

		v.ReadFromEnv()
		So(v.GetInt(key), ShouldEqual, cast.ToInt(age))
	})
}
