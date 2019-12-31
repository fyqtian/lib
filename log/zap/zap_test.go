package zap

import (
	"github.com/fyqtian/lib/config/viper"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSampleOptions(t *testing.T) {
	Convey("test SampleOption", t, func() {
		op := SampleOptions("log", viper.GetSingleton())
		So(op.FilePath, ShouldEqual, viper.GetSingleton().GetString("log.filepath"))
	})
}

func TestNewZap(t *testing.T) {
	o := &Options{
		FilePath:     "Runtime/test/test.log",
		FileSize:     1024,
		FileBackup:   3,
		FileMaxAge:   3,
		FileCompress: false,
		Debug:        true,
		Level:        "debug",
		Listen:       "127.0.0.1:9999",
	}
	logger := NewZap(o)

	logger.Debug("debug", zap.String("name", "van"))
	logger.Info("info", zap.String("name", "van"))

	u := url.URL{
		Scheme:     "http",
		Opaque:     "",
		User:       nil,
		Host:       o.Listen,
		Path:       "/zap/handle/level",
		RawPath:    "",
		ForceQuery: false,
		RawQuery:   "",
		Fragment:   "",
	}
	json := `{"level":"info"}`
	if req, err := http.NewRequest(http.MethodPut, u.String(), strings.NewReader(json)); err != nil {
		t.Fatal(err)
		return
	} else {
		req.Header.Add("Content-Type", "application/json")
		if resp, err := http.DefaultClient.Do(req); err != nil {
			t.Fatal(err)
			return
		} else {
			if resp.StatusCode != http.StatusOK {
				msg, _ := ioutil.ReadAll(resp.Body)
				t.Fatal("http code is not 200", string(msg))
			}
			logger.Debug("debug", zap.String("name", "van"))
			logger.Info("info", zap.String("name", "van"))
		}

	}

}

func TestDefaultZap(t *testing.T) {
	Convey("test DefaultZap", t, func() {
		So(DefaultZap(), ShouldNotEqual, nil)

	})
}
