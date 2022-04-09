package logger

import (
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/pkg/service/telegram"
	telegramDomain "heroku-line-bot/src/pkg/service/telegram/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"io"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	telegramLogger *telegramLoggerHandler
)

func GetTelegram() *telegramLoggerHandler {
	if telegramLogger == nil {
		if cfg, err := bootstrap.Get(); err == nil {
			telegramLogger = NewTelegramLogger(cfg)
		}
	}
	return telegramLogger
}

func GetTelegramLogger() ILogger {
	writerCreator := getTelegramWriterCreator(0)
	if writerCreator == nil {
		return nil
	}
	return newLoggerWithWriterCreator(writerCreator)
}

func GetTelegramWriter(name string, level zerolog.Level) io.Writer {
	writerCreator := getTelegramWriterCreator(0)
	if writerCreator == nil {
		return nil
	}
	return writerCreator.GetWriter(name, level)
}

func LogTelegram(name string, msg string, a ...interface{}) {
	logger := GetTelegramLogger()
	if logger == nil {
		return
	}

	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}

	logger.Log(name, errUtil.New(msg, zerolog.InfoLevel))
}

func getTelegramWriterCreator(loggerIndex int) IWriterCreator {
	logger := GetTelegram()
	if logger == nil {
		return nil
	}

	// 使用 loggers 打印錯誤
	return newHandleErrorWriterCreator(logger, loggerIndex)
}

type telegramLoggerHandler struct {
	limit       int
	telegramBot *telegram.Bot

	notifyIDs []int
}

func NewTelegramLogger(cfg *bootstrap.Config) *telegramLoggerHandler {
	channelAccessToken := cfg.TelegramBot.ChannelAccessToken
	bot := telegram.NewBot(channelAccessToken)
	if bot == nil {
		return nil
	}

	result := &telegramLoggerHandler{
		limit:       4096,
		telegramBot: bot,
		notifyIDs:   make([]int, 0),
	}

	telegramID, err := strconv.Atoi(cfg.TelegramBot.AdminID)
	if err != nil {
		return nil
	}
	if telegramID != 0 {
		result.notifyIDs = append(result.notifyIDs, telegramID)
	}
	return result
}

// 當 error 等級是錯誤以上時，通知 ProjectOwnerID
func (lh telegramLoggerHandler) GetWriter(name string, level zerolog.Level) io.Writer {
	return lh
}

func (lh telegramLoggerHandler) Write(p []byte) (n int, resultErr error) {
	return lh.sendMessage(p, lh.notifyIDs...)
}

func (lh telegramLoggerHandler) sendMessage(p []byte, chatIDs ...int) (n int, resultErr error) {
	n = len(p)

	if len(chatIDs) == 0 {
		return
	}

	msgs := make([][]byte, 0)
	for i := 0; i < len(p); i += lh.limit {
		to := i + lh.limit
		if to > len(p) {
			to = len(p)
		}
		s := p[i:to]
		msgs = append(msgs, s)
	}

	for _, chatID := range chatIDs {
		for _, s := range msgs {
			if err := lh.telegramBot.SendMessage(
				telegramDomain.ReqsSendMessage{
					ChatID: chatID,
					Text:   string(s),
				}, nil,
			); err != nil {
				resultErr = errUtil.NewError(err)
				return
			}
		}
	}
	return
}
