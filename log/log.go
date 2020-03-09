package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func Construct(cfg Config) {
	logger = log.New(out(cfg), "", log.LstdFlags)
}

var logger *log.Logger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)

// Configure log file
type Config = string

func Debug(v ...interface{})                 { println("DEBUG", v...) }
func Debugf(format string, v ...interface{}) { printf("DEBUG", format, v...) }

func Info(v ...interface{})                 { println("INFO", v...) }
func Infof(format string, v ...interface{}) { printf("INFO", format, v...) }

func Warn(v ...interface{})                 { println("WARN", v...) }
func Warnf(format string, v ...interface{}) { printf("WARN", format, v...) }

func Error(v ...interface{})                 { println("ERROR", v...) }
func Errorf(format string, v ...interface{}) { printf("ERROR", format, v...) }

func println(level string, v ...interface{}) {
	logger.Println(append([]interface{}{prefix(level)}, v...))
}
func printf(level string, format string, v ...interface{}) {
	logger.Printf(prefix(level)+format, v...)
}

func prefix(level string) string {
	buf := make([]byte, 32)
	runtime.Stack(buf[:], false)
	offset := len("goroutine ")
	goID := []byte{}
	for _, r := range buf[offset:] {
		if r < '0' {
			break
		}
		if r > '9' {
			break
		}
		goID = append(goID, r)
	}

	skipCaller, skipPath := 3, 5
	_, file, line, ok := runtime.Caller(skipCaller)
	if !ok {
		return "<???>:0"
	}

	separator := "/"
	subFile := strings.Split(file, separator)
	if len(subFile) > skipPath {
		subFile = subFile[len(subFile)-skipPath:]
	}
	filePath := fmt.Sprintf("%s:%d", strings.Join(subFile, separator), line)

	return fmt.Sprintf("Level:%s File:%s GoID:%s ",
		level, filePath, string(goID))
}
