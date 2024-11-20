package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// 로그레벨 정의
// TRACE > DEBUG > INFO > WARN > ERROR > FATAL
type LogLevel int

const (
	LogLevelFatal LogLevel = iota * 10
	LogLevelError
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

var (
	instance *Log
	once     sync.Once
)

// 로거타입과 그 구조체 멤버들은 외부로 노출될 필요가 없으므로 모두 소문자로 시작.
// 로거의 기능들은 대문자로 시작하는 메서드를 통하여 제공.
type Log struct {
	level       LogLevel
	path        string
	mutex       *sync.Mutex
	logFile     *os.File
	logFileName string
}

// 로거인스턴스 생성
// 파라메터: 로그를 저장할 디렉토리, 레벨 그리고 컨텍스트로그를 저장할 컬렉션명
func GetInstance(dir string, level LogLevel) *Log {
	once.Do(func() {
		instance = &Log{
			level: level,
			mutex: &sync.Mutex{},
		}

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(fmt.Sprintf("failed to create log directory: %v", err))
		}

		path, err := filepath.Abs(dir)
		if err != nil {
			panic(fmt.Sprintf("failed to get absolute path: %v", err))
		}
		instance.path = path

		if err := instance.setOutput(); err != nil {
			panic(fmt.Sprintf("failed to set initial output: %v", err))
		}
	})
	return instance
}

// 로그레벨 반환
func LevelType(level string) LogLevel {
	switch level {
	case "TRACE":
		return LogLevelTrace
	case "DEBUG":
		return LogLevelDebug
	case "INFO":
		return LogLevelInfo
	case "WARN":
		return LogLevelWarn
	case "ERROR":
		return LogLevelError
	case "FATAL":
		return LogLevelFatal
	default:
		return LogLevelTrace
	}
}

// 로그파일이 저장될 디렉토리의 절대경로 반환
func (l *Log) Path() string {
	return l.path
}

// 추적용 로그
func (l *Log) Trace(format string, v ...interface{}) {
	if l.level >= LogLevelTrace {
		l.Log("[TRACE] "+format, v...)
	}
}

// 디버그용 로그
func (l *Log) Debug(format string, v ...interface{}) {
	if l.level >= LogLevelDebug {
		l.Log("[DEBUG] "+format, v...)
	}
}

// 정보성 로그
func (l *Log) Info(format string, v ...interface{}) {
	if l.level >= LogLevelInfo {
		l.Log("[INFO] "+format, v...)
	}
}

// 경고수준 로그
func (l *Log) Warn(format string, v ...interface{}) {
	if l.level >= LogLevelWarn {
		l.Log("[WARN] "+format, v...)
	}
}

// 오류수준 로그
func (l *Log) Error(format string, v ...interface{}) {
	if l.level >= LogLevelError {
		l.Log("[ERROR] "+format, v...)
	}
}

// 치명적인 수준 로그
func (l *Log) Fatal(format string, v ...interface{}) {
	if l.level >= LogLevelFatal {
		l.Log("[FATAL] "+format, v...)
	}
}

// 로그레벨과 상관없는 로그
func (l *Log) Log(format string, v ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.setOutput()

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	file = filepath.Base(file)

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)
	fmt.Fprintf(l.logFile, "%s %s (%s:%d)\n", timestamp, msg, file, line)
}

func (l *Log) setOutput() error {
	newFileName := time.Now().Format("2006_01_02.log")
	if l.logFileName == newFileName {
		return nil
	}

	newFile, err := os.OpenFile(
		filepath.Join(l.path, newFileName),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	if l.logFile != nil {
		if err := l.logFile.Close(); err != nil {
			newFile.Close()
			return fmt.Errorf("failed to close old log file: %w", err)
		}
	}

	l.logFile = newFile
	l.logFileName = newFileName
	log.SetOutput(l.logFile)
	return nil
}

func (l *Log) Filename() string {
	return l.logFileName
}

func (l *Log) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}
