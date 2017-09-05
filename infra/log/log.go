/*
Package log provides provides a minimal interface for structured logging
*/
package log

import (
	"log"
	"os"

	gklog "github.com/go-kit/kit/log"
	gklevels "github.com/go-kit/kit/log/deprecated_levels"
)

const (
	loggerFormatJSON   = "json"
	loggerFormatLogFmt = "logfmt"
)

// Logger provides provides a minimal interface for structured logging
// It supplies leveled logging function which create a log event from keyvals,
// a variadic sequence of alternating keys and values.
type Logger interface {
	Debug(keyvals ...interface{})
	Info(keyvals ...interface{})
	Error(keyvals ...interface{})
	Warn(keyvals ...interface{})
	Crit(keyvals ...interface{})

	With(keyvals ...interface{}) Logger
}

// newLogger takes the name of the file and format of the logger as an argument, creates the file and returns
// a leveled logger that logs to the file.
// @format can have values logfmt or json. Default value is logfmt.
func newLogger(file string, format string) Logger {
	fw, err := GetFile(file)
	if err != nil {
		log.Fatal("error opening log file", err)
	}

	var l gklog.Logger
	if format == loggerFormatJSON {
		l = gklog.NewJSONLogger(fw)
	} else {
		l = gklog.NewLogfmtLogger(fw)
	}

	kitlevels := gklevels.New(
		l,

		// Fudge values so that switching between debug/info levels does not
		// mess with the log justification
		gklevels.DebugValue("dbug"),
		gklevels.ErrorValue("errr"),
	)

	kitlevels = kitlevels.With("ts", gklog.DefaultTimestampUTC)

	return levels{kitlevels}
}

//NewJSONLogger returns a Json Logger
func NewJSONLogger(fle string) Logger {
	return newLogger(fle, loggerFormatJSON)
}

//NewLogfmtLogger returns a LogfmtLogger Logger
func NewLogfmtLogger(fle string) Logger {
	return newLogger(fle, loggerFormatLogFmt)
}

//GetFile opens a file in read/write to append data to it
func GetFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

type levels struct {
	kit gklevels.Levels
}

func (l levels) Debug(keyvals ...interface{}) {
	if err := l.kit.Debug().Log(keyvals...); err != nil {
		log.Println("Error while logging(debug):", err)
	}
}

func (l levels) Info(keyvals ...interface{}) {
	if err := l.kit.Info().Log(keyvals...); err != nil {
		log.Println("Error while logging(info):", err)
	}
}

func (l levels) Error(keyvals ...interface{}) {
	if err := l.kit.Error().Log(keyvals...); err != nil {
		log.Println("Error while logging(error):", err)
	}
}
func (l levels) Warn(keyvals ...interface{}) {
	if err := l.kit.Warn().Log(keyvals...); err != nil {
		log.Println("Error while logging(warn):", err)
	}
}
func (l levels) Crit(keyvals ...interface{}) {
	if err := l.kit.Crit().Log(keyvals...); err != nil {
		log.Println("Error while logging(crit):", err)
	}
}

func (l levels) With(keyvals ...interface{}) Logger {
	return levels{l.kit.With(keyvals...)}
}
