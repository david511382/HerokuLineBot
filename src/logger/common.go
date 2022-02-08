package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/service/telegram"
	errUtil "heroku-line-bot/src/util/error"
	"io"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	telegramLogger *telegramLoggerHandler
	fileLogger     *fileLoggerHandler
	teminalLogger  *teminalLoggerHandler
	lokiLogger     *lokiLoggerHandler
	PanicWriter    io.Writer

	telegramBot      *telegram.Bot
	notifyTelegramID int
)

func init() {
	zerolog.LevelFieldName = "lvl"
	zerolog.ErrorStackMarshaler = errUtil.ErrorStackMarshaler
	zerolog.ErrorHandler = func(err error) {
		logOnTelegram(fmt.Sprintf("Log Fail:%s", err.Error()))
	}

	PanicWriter = &panicWriter{}
}

func getTelegramLogger() *telegramLoggerHandler {
	if telegramLogger == nil {
		telegramLogger = &telegramLoggerHandler{
			limit: 4096,
		}
	}
	return telegramLogger
}

func getFileLogger() *fileLoggerHandler {
	if fileLogger == nil {
		fileLogger = &fileLoggerHandler{}
	}
	return fileLogger
}

func getTeminalLogger() *teminalLoggerHandler {
	if teminalLogger == nil {
		teminalLogger = &teminalLoggerHandler{}
	}
	return teminalLogger
}

func getLokiLogger() *lokiLoggerHandler {
	if lokiLogger == nil {
		lokiLogger = NewLokiLog(nil)
	}
	return lokiLogger
}

func Init() errUtil.IError {
	cfg, errInfo := bootstrap.Get()
	if errInfo != nil {
		return errInfo
	}

	channelAccessToken := cfg.TelegramBot.ChannelAccessToken
	if telegramID, err := strconv.Atoi(cfg.TelegramBot.AdminID); err != nil {
		return errUtil.NewError(err)
	} else {
		notifyTelegramID = telegramID
	}
	telegramBot = telegram.NewBot(channelAccessToken)
	return nil
}

func Log(name string, logErrInfo errUtil.IError) {
	go func() {
		LogRightNow(name, logErrInfo)
	}()
}

func LogRightNow(name string, logErrInfo errUtil.IError) {
	if logErrInfo == nil {
		return
	}

	msg := message(name, logErrInfo)
	if logErrInfo.IsError() || logErrInfo.IsWarn() {
		if telegramBot == nil || notifyTelegramID == 0 {
			logOnFile(name, msg)
		} else {
			logOnTelegram(msg)
		}
	} else if logErrInfo.IsInfo() {
		logOnFile(name, msg)
	}
}

func logOnTelegram(msg string) {
	if errInfo := getTelegramLogger().log(notifyTelegramID, msg); errInfo != nil {
		logOnFile("System", msg)

		errInfo.AppendMessage("log telegram fail")
		errInfo.AppendMessage(msg)
		logOnTeminal(errInfo.ErrorWithTrace())
	}
}

func logOnFile(name, msg string) {
	if errInfo := getFileLogger().log(name, msg); errInfo != nil {
		logOnTeminal(msg)

		errInfo.AppendMessage("log file fail")
		errInfo.AppendMessage(msg)
		logOnTeminal(errInfo.ErrorWithTrace())
	}
}

func logOnTeminal(msg string) {
	if errInfo := getTeminalLogger().log("", msg); errInfo != nil {
		fmt.Println(msg)

		errInfo.AppendMessage("log teminal fail")
		errInfo.AppendMessage(msg)
		fmt.Println(errInfo.ErrorWithTrace())
	}
}

func logOnLoki(name string, err error) {
	getLokiLogger().log(name, err)
}

func message(name string, errInfo errUtil.IError) string {
	msg := errInfo.ErrorWithTrace()
	if errInfo.IsError() {
		return fmt.Sprintf("%s: ERROR: %s\n", name, msg)
	} else if errInfo.IsWarn() {
		return fmt.Sprintf("%s: WARN: %s\n", name, msg)
	} else if errInfo.IsInfo() {
		return fmt.Sprintf("%s: %s\n", name, msg)
	} else {
		return fmt.Sprintf("UNDEFIND ERROR: %s ON NAME: %s\n", msg, name)
	}
}
