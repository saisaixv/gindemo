package log

import (
	. "github.com/smartystreets/goconvey/convey"
	"log/syslog"
	"os"
	"testing"
)

func Test_Log(t *testing.T) {
	Convey("test log", t, func() {
		Convey("test", func() {
			l := New(LOG_INFO, NewLogMux(os.Stderr, "kk", LstdFlags|Lshortfile))
			So(l, ShouldNotBeNil)
			l.Infof("%d %s\n", 5, "bb")
			l.Debugf("bbbb\n")
			l.Notice("kkk\n")
		})
		Convey("syslog", func() {
			l := New(LOG_INFO)
			m := NewSyslogMux("", "", syslog.LOG_LOCAL1, "cydex_ts", false, 0)
			So(m, ShouldNotBeNil)
			l.AddMuxer(m)
			l.Infof("%d %s\n", 5, "bb")
			l.Debug("this is debug")
			l.Notice("kkk")
		})
		Convey("syslog async", func() {
			l := New(LOG_INFO)
			m := NewSyslogMux("", "", syslog.LOG_LOCAL1, "cydex_ts", true, 100)
			So(m, ShouldNotBeNil)
			l.AddMuxer(m)
			l.Infof("%d %s\n", 5, "bb")
			l.Debug("this is debug")
			l.Notice("kkk async")
		})
	})
}

func Test_Std(t *testing.T) {
	Convey("test std", t, func() {
		Convey("test", func() {
			Info("this is info")
		})
		Convey("add nil mux", func() {
			AddMuxer(nil)
			Info("thi is info")
		})
	})

}
