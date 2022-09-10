package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type NullLogger struct {
}

func (n *NullLogger) Debug(args ...interface{}) {
}

func (n *NullLogger) Info(args ...interface{}) {
}

func (n *NullLogger) Warning(args ...interface{}) {

}

func (n *NullLogger) Error(args ...interface{}) {
}

func (n *NullLogger) Fatal(args ...interface{}) {
}

func (n *NullLogger) Debugf(format string, args ...interface{}) {
}

func (n *NullLogger) Infof(format string, args ...interface{}) {
}

func (n *NullLogger) Warningf(format string, args ...interface{}) {
}

func (n *NullLogger) Errorf(format string, args ...interface{}) {
}

func (n *NullLogger) Fatalf(format string, args ...interface{}) {
}

type LogLevel int

const (
	DEBUG   LogLevel = 1
	INFO    LogLevel = 2
	WARNING LogLevel = 3
	ERROR   LogLevel = 4
	FATAL   LogLevel = 5
)

type logger struct {
	out io.Writer
}

func NewLogger(output io.Writer) Logger {
	return &logger{
		out: output,
	}
}

func (l *logger) log(level string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprint(args...)
	logRecord := fmt.Sprintf("[%s]\t%s\t%s\n", now, level, message)
	// l.file.Write([]byte(logRecord))
	fmt.Print(logRecord)
}

func (l *logger) logf(level, format string, args ...interface{}) {
	now := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logRecord := fmt.Sprintf("[%s]\t%s\t%s\n", now, level, message)
	// l.file.Write([]byte(logRecord))
	fmt.Print(logRecord)
}

func (l *logger) Debug(args ...interface{}) {
	l.log("DEBUG", args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log("INFO", args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.log("WARNING", args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log("ERROR", args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.log("FATAL", args...)
	os.Exit(2)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logf("DEBUG", format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logf("INFO", format, args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.logf("WARNING", format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logf("ERROR", format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logf("FATAL", format, args...)
	os.Exit(2)
}
