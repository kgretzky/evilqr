package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/fatih/color"
)

var stdout io.Writer = color.Output
var debug_output = true
var mtx_log *sync.Mutex = &sync.Mutex{}

const (
	DEBUG = iota
	INFO
	IMPORTANT
	WARNING
	ERROR
	FATAL
	SUCCESS
)

var LogLabels = map[int]string{
	DEBUG:     "dbg",
	INFO:      "inf",
	IMPORTANT: "imp",
	WARNING:   "war",
	ERROR:     "err",
	FATAL:     "!!!",
	SUCCESS:   "+++",
}

func DebugEnable(enable bool) {
	debug_output = enable
}

func SetOutput(o io.Writer) {
	stdout = o
}

func GetOutput() io.Writer {
	return stdout
}

func NullLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

func Debug(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	if debug_output {
		fmt.Fprint(stdout, format_msg(DEBUG, format+"\n", args...))
	}
}

func Info(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(INFO, format+"\n", args...))
}

func Important(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(IMPORTANT, format+"\n", args...))
}

func Warning(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(WARNING, format+"\n", args...))
}

func Error(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(ERROR, format+"\n", args...))
}

func Fatal(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(FATAL, format+"\n", args...))
}

func Success(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprint(stdout, format_msg(SUCCESS, format+"\n", args...))
}

func Printf(format string, args ...interface{}) {
	mtx_log.Lock()
	defer mtx_log.Unlock()

	fmt.Fprintf(stdout, format, args...)
}

func format_msg(lvl int, format string, args ...interface{}) string {
	t := time.Now()
	var sign, msg *color.Color
	switch lvl {
	case DEBUG:
		sign = color.New(color.FgBlack, color.BgHiBlack)
		msg = color.New(color.Reset, color.FgHiBlack)
	case INFO:
		sign = color.New(color.FgGreen, color.BgBlack)
		msg = color.New(color.Reset)
	case IMPORTANT:
		sign = color.New(color.FgWhite, color.BgHiBlue)
		//msg = color.New(color.Reset, color.FgHiBlue)
		msg = color.New(color.Reset)
	case WARNING:
		sign = color.New(color.FgHiYellow, color.BgBlack)
		//msg = color.New(color.Reset, color.FgYellow)
		msg = color.New(color.Reset)
	case ERROR:
		sign = color.New(color.FgWhite, color.BgRed)
		msg = color.New(color.Reset, color.FgRed)
	case FATAL:
		sign = color.New(color.FgBlack, color.BgRed)
		msg = color.New(color.Reset, color.FgRed, color.Bold)
	case SUCCESS:
		sign = color.New(color.FgWhite, color.BgGreen)
		msg = color.New(color.Reset, color.FgGreen)
	}
	time_clr := color.New(color.Reset)
	return "\r[" + time_clr.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()) + "] [" + sign.Sprintf("%s", LogLabels[lvl]) + "] " + msg.Sprintf(format, args...)
}
