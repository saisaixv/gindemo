package log

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = log.LstdFlags
)

const (
	LOG_EMERG   = int(syslog.LOG_EMERG)
	LOG_ALERT   = int(syslog.LOG_ALERT)
	LOG_CRIT    = int(syslog.LOG_CRIT)
	LOG_ERR     = int(syslog.LOG_ERR)
	LOG_WARNING = int(syslog.LOG_WARNING)
	LOG_NOTICE  = int(syslog.LOG_NOTICE)
	LOG_INFO    = int(syslog.LOG_INFO)
	LOG_DEBUG   = int(syslog.LOG_DEBUG)
)

var (
	levels = map[int]string{
		LOG_EMERG:   "EMERG",
		LOG_ALERT:   "ALERT",
		LOG_CRIT:    "CRITICAL",
		LOG_ERR:     "ERROR",
		LOG_WARNING: "WARNING",
		LOG_NOTICE:  "NOTICE",
		LOG_INFO:    "INFO",
		LOG_DEBUG:   "DEBUG",
	}
)

const (
	namePrefix = "LEVEL"
	levelDepth = 4
)

func AddBracket() {
	for k, v := range levels {
		levels[k] = "[" + v + "]"
	}
}

func AddColon() {
	for k, v := range levels {
		levels[k] = v + ":"
	}
}

func SetLevelName(level int, name string) {
	levels[level] = name
}

func LevelName(level int) string {
	if name, ok := levels[level]; ok {
		return name
	}
	return namePrefix + strconv.Itoa(level)
}

func NameLevel(name string) int {
	for k, v := range levels {
		if v == name {
			return k
		}
	}
	var level int
	if strings.HasPrefix(name, namePrefix) {
		level, _ = strconv.Atoi(name[len(namePrefix):])
	}
	return level
}

type Muxer interface {
	Output(level, calldepth int, s string) error
}

type Logger struct {
	level  int
	muxers []Muxer
	mu     sync.Mutex
}

func New(level int, muxers ...Muxer) *Logger {
	l := new(Logger)
	l.level = level
	for _, muxer := range muxers {
		if muxer != nil {
			l.muxers = append(l.muxers, muxer)
		}
	}
	return l
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) AddMuxer(m Muxer) {
	if m != nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.muxers = append(l.muxers, m)
	}
}

func (l *Logger) output(level, calldepth int, s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, mux := range l.muxers {
		mux.Output(level, calldepth+1, s)
	}
	return nil
}

func (l *Logger) Err(level, calldepth int, err error) error {
	if err != nil && level <= l.level {
		return l.output(level, calldepth, err.Error())
	}
	return nil
}

func (l *Logger) Output(level, calldepth int, a ...interface{}) error {
	if level <= l.level {
		return l.output(level, calldepth, fmt.Sprint(a...))
	}
	return nil
}

func (l *Logger) Outputf(level, calldepth int, format string, a ...interface{}) error {
	if level <= l.level {
		return l.output(level, calldepth, fmt.Sprintf(format, a...))
	}
	return nil
}

func (l *Logger) Outputln(level, calldepth int, a ...interface{}) error {
	if level <= l.level {
		return l.output(level, calldepth, fmt.Sprintln(a...))
	}
	return nil
}

func (l *Logger) Debug(a ...interface{}) {
	l.Output(LOG_DEBUG, levelDepth, a...)
}

func (l *Logger) Notice(a ...interface{}) {
	l.Output(LOG_NOTICE, levelDepth, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.Output(LOG_INFO, levelDepth, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.Output(LOG_WARNING, levelDepth, a...)
}

func (l *Logger) Error(a ...interface{}) {
	l.Output(LOG_ERR, levelDepth, a...)
}

func (l *Logger) Critical(a ...interface{}) {
	l.Output(LOG_CRIT, levelDepth, a...)
}

func (l *Logger) Panic(a ...interface{}) {
	s := fmt.Sprint(a...)
	if LOG_EMERG <= l.level {
		l.output(LOG_EMERG, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Output(LOG_EMERG, levelDepth, a...)
	os.Exit(1)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Outputf(LOG_DEBUG, levelDepth, format, a...)
}

func (l *Logger) Noticef(format string, a ...interface{}) {
	l.Outputf(LOG_NOTICE, levelDepth, format, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Outputf(LOG_INFO, levelDepth, format, a...)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Outputf(LOG_WARNING, levelDepth, format, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Outputf(LOG_ERR, levelDepth, format, a...)
}

func (l *Logger) Criticalf(format string, a ...interface{}) {
	l.Outputf(LOG_CRIT, levelDepth, format, a...)
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	if LOG_EMERG <= l.level {
		l.output(LOG_EMERG, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Outputf(LOG_EMERG, levelDepth, format, a...)
	os.Exit(1)
}

func (l *Logger) Debugln(a ...interface{}) {
	l.Outputln(LOG_DEBUG, levelDepth, a...)
}

func (l *Logger) Infoln(a ...interface{}) {
	l.Outputln(LOG_INFO, levelDepth, a...)
}

func (l *Logger) Noticeln(a ...interface{}) {
	l.Outputln(LOG_NOTICE, levelDepth, a...)
}

func (l *Logger) Warningln(a ...interface{}) {
	l.Outputln(LOG_WARNING, levelDepth, a...)
}

func (l *Logger) Errorln(a ...interface{}) {
	l.Outputln(LOG_ERR, levelDepth, a...)
}

func (l *Logger) Criticalln(a ...interface{}) {
	l.Outputln(LOG_CRIT, levelDepth, a...)
}

func (l *Logger) Panicln(a ...interface{}) {
	s := fmt.Sprintln(a...)
	if LOG_EMERG <= l.level {
		l.output(LOG_EMERG, levelDepth-1, s)
	}
	panic(s)
}

func (l *Logger) Fatalln(a ...interface{}) {
	l.Outputln(LOG_EMERG, levelDepth, a...)
	os.Exit(1)
}

var std = New(LOG_DEBUG, NewLogMux(os.Stderr, "", LstdFlags|Lshortfile))

func AddMuxer(m Muxer) {
	std.AddMuxer(m)
}

func SetLevel(level int) {
	std.level = level
}

// func Emerg(a ...interface{}) {
// 	std.Output(LOG_EMERG, levelDepth, a...)
// }
//
// func Alert(a ...interface{}) {
// 	std.Output(LOG_ALERT, levelDepth, a...)
// }

func Fatal(a ...interface{}) {
	std.Output(LOG_EMERG, levelDepth, a...)
	os.Exit(1)
}

func Critical(a ...interface{}) {
	std.Output(LOG_CRIT, levelDepth, a...)
}

func Error(a ...interface{}) {
	std.Output(LOG_ERR, levelDepth, a...)
}

func Warning(a ...interface{}) {
	std.Output(LOG_WARNING, levelDepth, a...)
}

func Notice(a ...interface{}) {
	std.Output(LOG_NOTICE, levelDepth, a...)
}

func Info(a ...interface{}) {
	std.Output(LOG_INFO, levelDepth, a...)
}

func Debug(a ...interface{}) {
	std.Output(LOG_DEBUG, levelDepth, a...)
}

func Fatalf(format string, a ...interface{}) {
	std.Outputf(LOG_EMERG, levelDepth, format, a...)
	os.Exit(1)
}

func Criticalf(format string, a ...interface{}) {
	std.Outputf(LOG_CRIT, levelDepth, format, a...)
}

func Errorf(format string, a ...interface{}) {
	std.Outputf(LOG_ERR, levelDepth, format, a...)
}

func Warningf(format string, a ...interface{}) {
	std.Outputf(LOG_WARNING, levelDepth, format, a...)
}

func Noticef(format string, a ...interface{}) {
	std.Outputf(LOG_NOTICE, levelDepth, format, a...)
}

func Infof(format string, a ...interface{}) {
	std.Outputf(LOG_INFO, levelDepth, format, a...)
}

func Debugf(format string, a ...interface{}) {
	std.Outputf(LOG_DEBUG, levelDepth, format, a...)
}

func Fatalln(a ...interface{}) {
	std.Outputln(LOG_EMERG, levelDepth, a...)
	os.Exit(1)
}

func Criticalln(a ...interface{}) {
	std.Outputln(LOG_CRIT, levelDepth, a...)
}

func Errorln(a ...interface{}) {
	std.Outputln(LOG_ERR, levelDepth, a...)
}

func Warningln(a ...interface{}) {
	std.Outputln(LOG_WARNING, levelDepth, a...)
}

func Noticeln(a ...interface{}) {
	std.Outputln(LOG_NOTICE, levelDepth, a...)
}

func Infoln(a ...interface{}) {
	std.Outputln(LOG_INFO, levelDepth, a...)
}

func Debugln(a ...interface{}) {
	std.Outputln(LOG_DEBUG, levelDepth, a...)
}
