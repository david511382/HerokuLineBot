package logger

import (
	"heroku-line-bot/bootstrap"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	loggers     *Loggers
	writers     *Writers
	loggerTypes *LoggerTypes
)

func init() {
	loggers = &Loggers{}
	writers = &Writers{}
	loggerTypes = &LoggerTypes{}
}

// 會選擇 index 0 的 logger 來打印
func GetLogger() ILogger {
	return loggers.GetLogger()
}

type Loggers struct {
	mux     sync.RWMutex // 讀寫鎖
	loggers []ILogger
}

// 會選擇 index 0 的 logger 來打印
func (ls *Loggers) GetLogger() ILogger {
	ls.mux.RLock()
	isInited := ls.loggers != nil
	if isInited {
		var logger ILogger
		if len(ls.loggers) != 0 {
			logger = ls.loggers[0]
		}
		ls.mux.RUnlock()
		return logger
	}
	ls.mux.RUnlock()

	ls.mux.Lock()
	defer ls.mux.Unlock()

	if ls.loggers == nil {
		ls.loggers = make([]ILogger, 0)
		for _, loggerType := range getLoggerTypes() {
			if loggerType.writer != nil {
				ls.loggers = append(ls.loggers, newLogger(loggerType.writer))
				continue
			}
			if loggerType.writerCreator != nil {
				ls.loggers = append(ls.loggers, newLoggerWithWriterCreator(loggerType.writerCreator))
				continue
			}
		}
	}

	{
		var logger ILogger
		if len(ls.loggers) != 0 {
			logger = ls.loggers[0]
		}
		return logger
	}
}

func (ls *Loggers) Reset() {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.loggers = nil
}

func GetWriter(name string, level zerolog.Level) io.Writer {
	return writers.GetWriter(name, level)
}

type Writers struct {
	mux     sync.RWMutex // 讀寫鎖
	writers []io.Writer
}

// 會選擇 index 0 的 logger 來打印
func (ls *Writers) GetWriter(name string, level zerolog.Level) io.Writer {
	ls.mux.RLock()
	isInited := ls.writers != nil
	if isInited {
		var writer io.Writer
		if len(ls.writers) != 0 {
			writer = ls.writers[0]
		}
		ls.mux.RUnlock()
		return writer
	}
	ls.mux.RUnlock()

	ls.mux.Lock()
	defer ls.mux.Unlock()

	if ls.writers == nil {
		ls.writers = make([]io.Writer, 0)
		for _, loggerType := range getLoggerTypes() {
			if loggerType.writer != nil {
				ls.writers = append(ls.writers, loggerType.writer)
				continue
			}
			if loggerType.writerCreator != nil {
				w := loggerType.writerCreator.GetWriter(name, level)
				ls.writers = append(ls.writers, w)
				continue
			}
		}
	}

	{
		var writer io.Writer
		if len(ls.writers) != 0 {
			writer = ls.writers[0]
		}
		return writer
	}
}

func (ls *Writers) Reset() {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	ls.writers = nil
}

type loggerType struct {
	writerCreator IWriterCreator
	writer        io.Writer
}

// 宣告可以打印用的工具，加入的順序就是用來打印的優先權，index 小的優先
func getLoggerTypes() []loggerType {
	return loggerTypes.GetLoggerTypes()
}

type LoggerTypes struct {
	mux            sync.RWMutex // 讀寫鎖
	loggerTypes    []loggerType
	deleteIndexMap map[int]bool
}

// 宣告註冊的 logger
func (ls *LoggerTypes) loadLoggerTypes() {
	ls.loggerTypes = make([]loggerType, 0)

	if cfg, err := bootstrap.Get(); err == nil {
		if logger := NewLokiLogger(cfg); logger != nil {
			index := len(ls.loggerTypes)
			lt := loggerType{
				writerCreator: newHandleErrorWriterCreator(logger, index),
			}
			ls.loggerTypes = append(ls.loggerTypes, lt)
		}

		if logger := NewFileLogger(cfg); logger != nil {
			index := len(ls.loggerTypes)
			lt := loggerType{
				writerCreator: newHandleErrorWriterCreator(logger, index),
			}
			ls.loggerTypes = append(ls.loggerTypes, lt)
		}
	}

	consoleWriter := os.Stdout
	index := len(ls.loggerTypes)
	lt := loggerType{
		writer: newHandleErrorWriter(consoleWriter, index),
	}
	ls.loggerTypes = append(ls.loggerTypes, lt)
}

func (ls *LoggerTypes) getLoggerTypes() []loggerType {
	result := make([]loggerType, 0)
	for index, l := range ls.loggerTypes {
		if ls.deleteIndexMap[index] {
			continue
		}

		result = append(result, l)
	}
	return result
}

// 取得註冊的 logger
// ，不適合頻繁讀取，使用者需快取
func (ls *LoggerTypes) GetLoggerTypes() []loggerType {
	ls.mux.RLock()
	isInited := ls.loggerTypes != nil
	if isInited {
		result := ls.getLoggerTypes()
		ls.mux.RUnlock()
		return result
	}
	ls.mux.RUnlock()

	ls.mux.Lock()
	defer ls.mux.Unlock()

	if ls.loggerTypes == nil {
		ls.loadLoggerTypes()
	}

	return ls.getLoggerTypes()
}

func (ls *LoggerTypes) Delete(index int) {
	ls.mux.Lock()
	defer ls.mux.Unlock()
	if ls.deleteIndexMap == nil {
		ls.deleteIndexMap = make(map[int]bool)
	}
	_, exist := ls.deleteIndexMap[index]
	if exist {
		return
	}

	ls.deleteIndexMap[index] = true
}

func handleLoggerIndexFail(index int) {
	loggerTypes.Delete(index)
	loggers.Reset()
	writers.Reset()
}
