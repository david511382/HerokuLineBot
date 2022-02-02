package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/service/telegram"
	errUtil "heroku-line-bot/src/util/error"
	"io"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
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
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		e, ok := err.(*errUtil.ErrorInfo)
		if ok {
			return e.TraceMessage()
		}
		return pkgerrors.MarshalStack(err)
	}
	zerolog.ErrorHandler = func(err error) {
		logOnTelegram(fmt.Sprintf("Log Fail:%s", err.Error()))
	}

	telegramLogger = &telegramLoggerHandler{
		limit: 4096,
	}
	fileLogger = &fileLoggerHandler{}
	teminalLogger = &teminalLoggerHandler{}
	PanicWriter = &panicWriter{}
	lokiLogger = NewLokiLog(nil)
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
	if errInfo := telegramLogger.log(notifyTelegramID, msg); errInfo != nil {
		logOnFile("System", msg)

		errInfo = errInfo.NewParent("log telegram fail")
		errInfo = errInfo.NewParent(msg)
		logOnTeminal(errInfo.ErrorWithTrace())
	}
}

func logOnFile(name, msg string) {
	if errInfo := fileLogger.log(name, msg); errInfo != nil {
		logOnTeminal(msg)

		errInfo = errInfo.NewParent("log file fail")
		errInfo = errInfo.NewParent(msg)
		logOnTeminal(errInfo.ErrorWithTrace())
	}
}

func logOnTeminal(msg string) {
	if errInfo := teminalLogger.log("", msg); errInfo != nil {
		fmt.Println(msg)

		errInfo = errInfo.NewParent("log teminal fail")
		errInfo = errInfo.NewParent(msg)
		fmt.Println(errInfo.ErrorWithTrace())
	}
}

func logOnLoki(msg string) {
	errInfo := lokiLogger.log(
		LoggerInfo{
			Message: msg,
		},
	)
	if errInfo != nil {
		logOnTeminal(msg)

		errInfo = errInfo.NewParent("log loki fail")
		errInfo = errInfo.NewParent(msg)
		logOnTeminal(errInfo.ErrorWithTrace())
	}
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
