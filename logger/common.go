package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/service/telegram"
	"io"
	"strconv"
)

var (
	telegramLogger *telegramLoggerHandler
	fileLogger     *fileLoggerHandler
	teminalLogger  *teminalLoggerHandler
	PanicWriter    io.Writer

	telegramBot      *telegram.Bot
	notifyTelegramID int
)

func init() {
	telegramLogger = &telegramLoggerHandler{
		limit: 4096,
	}
	fileLogger = &fileLoggerHandler{}
	teminalLogger = &teminalLoggerHandler{}
	PanicWriter = &panicWriter{}
}

func Init(cfg *bootstrap.Config) *errLogic.ErrorInfo {
	channelAccessToken := cfg.TelegramBot.ChannelAccessToken
	if telegramID, err := strconv.Atoi(cfg.TelegramBot.AdminID); err != nil {
		return errLogic.NewError(err)
	} else {
		notifyTelegramID = telegramID
	}
	telegramBot = telegram.NewBot(channelAccessToken)
	return nil
}

func Log(name string, LogErrInfo *errLogic.ErrorInfo) {
	go func() {
		LogRightNow(name, LogErrInfo)
	}()
}

func LogRightNow(name string, LogErrInfo *errLogic.ErrorInfo) {
	if LogErrInfo == nil {
		return
	}

	msg := message(name, LogErrInfo)
	switch LogErrInfo.Level {
	case errLogic.ERROR, errLogic.WARN:
		if telegramBot == nil || notifyTelegramID == 0 {
			logOnFile(name, msg)
		} else {
			logOnTelegram(msg)
		}
	case errLogic.INFO:
		logOnFile(name, msg)
	}
}

func logOnTelegram(msg string) {
	if errInfo := telegramLogger.log(notifyTelegramID, msg); errInfo != nil {
		logOnFile("System", msg)

		errInfo = errInfo.NewParent("log telegram fail")
		errInfo = errInfo.NewParent(msg)
		logOnTeminal(errInfo.Error().Error())
	}
}

func logOnFile(name, msg string) {
	if errInfo := fileLogger.log(name, msg); errInfo != nil {
		logOnTeminal(msg)

		errInfo = errInfo.NewParent("log file fail")
		errInfo = errInfo.NewParent(msg)
		logOnTeminal(errInfo.Error().Error())
	}
}

func logOnTeminal(msg string) {
	if errInfo := teminalLogger.log("", msg); errInfo != nil {
		fmt.Println(msg)

		errInfo = errInfo.NewParent("log teminal fail")
		errInfo = errInfo.NewParent(msg)
		fmt.Println(errInfo.Error().Error())
	}
}

func message(name string, errInfo *errLogic.ErrorInfo) string {
	msg := errInfo.Error().Error()
	switch errInfo.Level {
	case errLogic.ERROR:
		return fmt.Sprintf("%s: ERROR: %s\n", name, msg)
	case errLogic.WARN:
		return fmt.Sprintf("%s: WARN: %s\n", name, msg)
	case errLogic.INFO:
		return fmt.Sprintf("%s: %s\n", name, msg)
	default:
		return fmt.Sprintf("UNDEFIND LEVEL: %s ON NAME: %s\n", errInfo.Level, name)
	}
}
