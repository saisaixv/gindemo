package log

import (
	"bytes"
	"log"
	"log/syslog"
)

type item struct {
	level     int
	calldepth int
	s         string
}

type SyslogMux struct {
	l        *log.Logger
	w        *syslog.Writer
	async    bool
	buf      bytes.Buffer
	c        chan *item
	bufSize  int
	network  string
	raddr    string
	tag      string
	priority syslog.Priority
}

func NewSyslogMux(network, raddr string, priority syslog.Priority, tag string, async bool, bufSize int) *SyslogMux {
	var err error
	m := &SyslogMux{}
	m.w, err = syslog.Dial(network, raddr, priority, tag)
	if err != nil {
		println(err.Error())
		return nil
	}
	m.l = log.New(&m.buf, "", log.Lshortfile)
	if m.async {
		m.c = make(chan *item, bufSize)
		go func() {
			for item := range m.c {
				m.output(item.level, item.calldepth, item.s)
			}
		}()
	}

	return m
}

func (m *SyslogMux) Output(level, calldepth int, s string) error {
	if m.c != nil {
		m.c <- &item{level, calldepth + 1, s}
		return nil
	}
	return m.output(level, calldepth+1, s)
}

func (m *SyslogMux) output(level, calldepth int, s string) error {
	err := m.l.Output(calldepth, s)
	if err != nil {
		return err
	}
	outs := m.buf.String()
	m.buf.Reset()

	switch level {
	case LOG_EMERG:
		return m.w.Emerg(outs)
	case LOG_ALERT:
		return m.w.Alert(outs)
	case LOG_CRIT:
		return m.w.Crit(outs)
	case LOG_ERR:
		return m.w.Err(outs)
	case LOG_WARNING:
		return m.w.Warning(outs)
	case LOG_NOTICE:
		return m.w.Notice(outs)
	case LOG_INFO:
		return m.w.Info(outs)
	case LOG_DEBUG:
		return m.w.Debug(outs)
	default:
		return nil
	}
}
