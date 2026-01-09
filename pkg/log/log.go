/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	traceStr = "TRACE"
	debugStr = "DEBUG"
	infoStr  = "INFO"
	warnStr  = "WARN"
	errorStr = "ERROR"
	fatalStr = "FATAL"
)

var (
	// LogLevel can specify what level the system runs at
	logLevel    = 6
	logLevelStr = "TRACE"
	levels      = map[string]int{
		"TRACE": 6,
		"6":     6,
		"DEBUG": 5,
		"5":     5,
		"INFO":  4,
		"4":     4,
		"WARN":  3,
		"3":     3,
		"ERROR": 2,
		"2":     2,
		"FATAL": 1,
		"1":     1,
	}
	stlGoLog   = log.New(os.Stderr, "[stl-go] ", log.LstdFlags|log.Llongfile)
	writeMutex sync.Mutex
)

// FormatMilliseconds will output time in a more human readable fashion
func FormatMilliseconds(milliseconds int64) string {
	totalSeconds := milliseconds / 1000
	ms := milliseconds % 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	hours := minutes / 60
	minutes %= 60
	days := hours / 24
	hours %= 24
	return fmt.Sprintf("%dd:%02dh:%02dm:%02ds:%03dms", days, hours, minutes, seconds, ms)
}

// FunctionTimer is a deferable function call to time how long something takes
func FunctionTimer() func() {
	start := time.Now()

	functionName := "unknown"
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		name := details.Name()
		functionName = name[strings.LastIndex(name, ".")+1:]
	}

	return func() {
		duration := time.Since(start)
		writeLog(6, fmt.Sprintf("Function '%s' completed in %v", functionName, FormatMilliseconds(duration.Abs().Milliseconds())))
	}
}

// SetLogLevel will set the log level by string
func SetLogLevel(levelStr string) {
	levelStr = strings.ToUpper(levelStr)
	if level, ok := levels[levelStr]; ok {
		logLevel = level
		logLevelStr = levelStr
	}
	Debugf("log level set to %s", logLevelStr)
}

// Init initializes settings related to logging
func Init(levelStr string, out io.Writer) {
	SetLogLevel(levelStr)
	stlGoLog.SetOutput(out)
}

// Trace is a wrapper for log.Trace
func Trace(v ...interface{}) {
	writeLog(6, traceStr, v...)
}

// Tracef is a wrapper for log.Tracef
func Tracef(format string, v ...interface{}) {
	writeLog(6, traceStr, fmt.Sprintf(format, v...))
}

// Debug is a wrapper for log.Debug
func Debug(v ...interface{}) {
	writeLog(5, debugStr, v...)
}

// Debugf is a wrapper for log.Debugf
func Debugf(format string, v ...interface{}) {
	writeLog(5, debugStr, fmt.Sprintf(format, v...))
}

// Info is a wrapper for log.Info
func Info(v ...interface{}) {
	writeLog(4, infoStr, v...)
}

// Infof is a wrapper for log.Infof
func Infof(format string, v ...interface{}) {
	writeLog(4, infoStr, fmt.Sprintf(format, v...))
}

// Warn is a wrapper for log.Warn
func Warn(v ...interface{}) {
	writeLog(3, warnStr, v...)
}

// Warnf is a wrapper for log.Warnf
func Warnf(format string, v ...interface{}) {
	writeLog(3, warnStr, fmt.Sprintf(format, v...))
}

// Error is a wrapper for log.Error
func Error(v ...interface{}) {
	writeLog(2, errorStr, v...)
}

// Errorf is a wrapper for log.Errorf
func Errorf(format string, v ...interface{}) {
	writeLog(2, errorStr, fmt.Sprintf(format, v...))
}

// Fatal is a wrapper for log.Fatal
func Fatal(v ...interface{}) {
	writeLog(1, fatalStr, v...)
	os.Exit(-1)
}

// Fatalf is a wrapper for log.Fatalf
func Fatalf(format string, v ...interface{}) {
	writeLog(1, fatalStr, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

// Writer returns log output writer object
func Writer() io.Writer {
	return stlGoLog.Writer()
}

// Logger is used by things like net/http to overwrite their standard logging
func Logger() *log.Logger {
	return stlGoLog
}

// writeLog outputs the line to the configured output
func writeLog(level int, levelStr string, v ...interface{}) {
	// determine if we need to display the logs
	if level <= logLevel {
		writeMutex.Lock()
		defer writeMutex.Unlock()
		// the origionall caller of this is 3 steps back, the output will display who called it
		err := stlGoLog.Output(3, fmt.Sprintf("[%s] %v", levelStr, fmt.Sprint(v...)))
		if err != nil {
			stlGoLog.Print(v...)
			stlGoLog.Print(err)
		}
	}
}
