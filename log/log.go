package log

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hmerritt/autocost/utils"
	"github.com/hmerritt/autocost/version"
)

const (
	LOG_FILE  = version.AppName + ".log" // @TODO Default log file (can be overridden by config file)
	LOG_LEVEL = 4                        // 5 = debug, 4 = info, 3 = warn, 2 = error
)

type Logger struct {
	// The logging level the logger should log at. This is typically (and defaults
	// to) `Info`, which allows Info(), Warn(), Error() and Fatal() to be logged.
	Level uint32

	// Timestamp of Logger initiation.
	InitTimestamp time.Time

	// Timestamp of the most recent log. Used to calculate and show the time in
	// milliseconds since last log.
	PrevTimestamp time.Time

	// Log to file
	File         *os.File
	Chan         chan []interface{}
	ChanWg       sync.WaitGroup
	ChanIsClosed bool
}

func NewLogger() *Logger {
	return &Logger{
		Level:         LOG_LEVEL,
		InitTimestamp: time.Now(),
		PrevTimestamp: time.Now(),
		Chan:          make(chan []interface{}),
		ChanWg:        sync.WaitGroup{},
		ChanIsClosed:  false,
	}
}

func (l *Logger) FileStart(filePath string) {
	logFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil && LOG_LEVEL >= 5 {
		color.Set(color.FgRed)
		fmt.Println("error opening log file:", err)
		color.Unset()
	}

	l.File = logFile
	l.ChanIsClosed = false
	l.Chan = make(chan []interface{})

	for a := range l.Chan {
		if l.ChanIsClosed {
			return
		}

		if l.File != nil {
			_, err := l.File.WriteString(fmt.Sprintln(a...))
			if err != nil && LOG_LEVEL >= 5 {
				color.Set(color.FgRed)
				fmt.Println("error writting to log file:", err)
				color.Unset()
			}
			l.ChanWg.Done()
		}
	}
}

func (l *Logger) FileClose() error {
	// Wait for unwritten logs to finish writing
	l.ChanWg.Wait()

	if l.File != nil {
		err := l.File.Close()

		if err != nil {
			return err
		}
	}

	l.ChanIsClosed = true
	close(l.Chan)

	return nil
}

func (l *Logger) log(level uint32, a ...interface{}) {
	if l.Level < level {
		return
	}
	if !l.ChanIsClosed {
		l.ChanWg.Add(1)
		go func() { l.Chan <- a }()
	}
	l.PrevTimestamp = time.Now()
	fmt.Println(a...)
}

func (l *Logger) SetLevel(level uint32) {
	l.Level = level
}

func (l *Logger) Error(messages ...interface{}) error {
	color.Set(color.FgRed)
	defer color.Unset()
	l.log(2, messages...)
	return errors.New(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(messages)), " "), "[]"))
}
func (l *Logger) Errorf(format string, messages ...interface{}) error {
	return l.Error(fmt.Sprintf(format, messages...))
}

func (l *Logger) Warn(messages ...interface{}) {
	color.Set(color.FgYellow)
	defer color.Unset()
	l.log(3, messages...)
}
func (l *Logger) Warnf(format string, messages ...interface{}) {
	l.Warn(fmt.Sprintf(format, messages...))
}

func (l *Logger) Success(messages ...interface{}) {
	color.Set(color.FgHiGreen)
	defer color.Unset()
	l.log(4, messages...)
}
func (l *Logger) Successf(format string, messages ...interface{}) {
	l.Success(fmt.Sprintf(format, messages...))
}

func (l *Logger) Info(messages ...interface{}) {
	l.log(4, messages...)
}
func (l *Logger) Infof(format string, messages ...interface{}) {
	l.Info(fmt.Sprintf(format, messages...))
}

func (l *Logger) Debug(messages ...interface{}) {
	// Function name
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	funcName = funcName[strings.LastIndex(funcName, ".")+1:] // Removes package name

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	toLog := fmt.Sprintf("%s (%s) +%9s => ", formattedTime, funcName, utils.DurationSince(l.PrevTimestamp))

	a := make([]interface{}, 0)
	a = append(a, toLog)
	a = append(a, messages...)

	l.log(5, a...)
}
func (l *Logger) Debugf(format string, messages ...interface{}) {
	l.Debug(fmt.Sprintf(format, messages...))
}

func (l *Logger) End() {
	color.Set(color.FgCyan)
	defer color.Unset()
	l.log(4, fmt.Sprintf("took %s", utils.DurationSince(l.InitTimestamp)))
}
