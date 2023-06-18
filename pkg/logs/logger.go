package logs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger interface {
	Log(s string)
}

// Logs lines of text with the following format:
//
// "%s (%s) %q\n" with timestamp, source code location and log message between quotes.
type TextLogger struct {
	mutex   *sync.Mutex
	writers []io.Writer
}

func NewTextLogger(writers ...io.Writer) *TextLogger {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stderr}
	}
	return &TextLogger{writers: writers, mutex: &sync.Mutex{}}
}

func (l *TextLogger) Log(s string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for _, w := range l.writers {
		timestr := time.Now().Format("2006-01-02 15:04:05.000 Z07:00")
		logstr := fmt.Sprintf("%s (%s) %q\n", timestr, getStackLevel(2), s)
		_, err := w.Write([]byte(logstr))
		if err != nil {
			panic(err)
		}
	}
}

// Opens a log file with the appropriate flag and mode.
func MustOpenLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0200)
	if err != nil {
		panic(err)
	}
	return f
}

// returns the filename and line of a calling function in the stack.
func getStackLevel(offset int) string {
	_, f, line, ok := runtime.Caller(offset)
	if !ok {
		panic(errors.New("invalid runtime caller offset"))
	}
	return fmt.Sprintf("%s:%d", f, line)
}
