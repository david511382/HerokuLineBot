package wsjob

import (
	"encoding/json"
	"fmt"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/server/domain/resp"
	"heroku-line-bot/src/server/ws"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
)

type IScheduleWsConnJobHandler interface {
	RunJob() (resp.Base, errUtil.IError)
	UpdateReqs(reqsBs []byte) (resultErrInfo errUtil.IError)
}

type BaseScheduleWsConnJob struct {
	c       *gin.Context
	handler IScheduleWsConnJobHandler
	sync.Mutex
	running bool
}

func NewBaseScheduleWsConnJob(
	c *gin.Context,
	handler IScheduleWsConnJobHandler,
) *BaseScheduleWsConnJob {
	return &BaseScheduleWsConnJob{
		c:       c,
		handler: handler,
	}
}

func (w *BaseScheduleWsConnJob) Run(wsSender ws.IWsSender) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	result, errInfo := w.handler.RunJob()
	if errInfo != nil {
		if errInfo.IsError() {
			w.Error(wsSender, errInfo)
			wsSender.Close()
			return
		}

		logger.LogError(logger.NAME_API, errInfo)
	}

	resultBs, err := json.Marshal(result)
	if err != nil {
		errInfo := errUtil.NewError(err)
		w.Error(wsSender, errInfo)
		wsSender.Close()
		return
	}

	if err := wsSender.Send(websocket.TextMessage, resultBs); err != nil {
		errInfo := errUtil.NewError(err)
		w.Error(wsSender, errInfo)
		wsSender.Close()
		return
	}
}

func (w *BaseScheduleWsConnJob) Listen(wsSender ws.IWsSender, wsMsg *ws.WsConnReadMessage) {
	if wsMsg.Err != nil {
		errInfo := errUtil.NewError(wsMsg.Err)
		w.Error(wsSender, errInfo)
		wsSender.Close()
		return
	}

	if wsMsg.MessageType != websocket.TextMessage {
		return
	}

	if errInfo := w.UpdateReqs(wsSender, wsMsg.P); errInfo != nil {
		w.Error(wsSender, errInfo)
		wsSender.Close()
		return
	}

	w.Run(wsSender)
}

func (w *BaseScheduleWsConnJob) Error(wsSender ws.IWsSender, errInfo errUtil.IError) {
	result := resp.Base{
		Message: errInfo.Error(),
	}
	resultBs, err := json.Marshal(result)
	if err != nil {
		logger.LogError(logger.NAME_API, errUtil.New(result.Message))
		logger.LogError(logger.NAME_API, errUtil.NewError(err))
		return
	}

	if err := wsSender.Send(websocket.TextMessage, resultBs); err != nil {
		logger.LogError(logger.NAME_API, errUtil.NewError(err))
		return
	}
}

func (w *BaseScheduleWsConnJob) UpdateReqs(wsSender ws.IWsSender, p []byte) (resultErrInfo errUtil.IError) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	if errInfo := w.handler.UpdateReqs(p); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	return
}

// reqsP : pointer of reqs
func (w *BaseScheduleWsConnJob) parseJson(reqsBs []byte, reqsP interface{}) error {
	if err := json.Unmarshal(reqsBs, reqsP); err != nil {
		errInfo := errUtil.NewError(err)
		return errInfo
	}
	if err := w.validate(reqsP); err != nil {
		errInfo := errUtil.NewError(err)
		return errInfo
	}
	locationConverter := util.NewLocationConverter(global.TimeUtilObj.GetLocation(), false)
	locationConverter.Convert(reqsP)
	return nil
}

func (w *BaseScheduleWsConnJob) validate(obj interface{}) error {
	if binding.Validator == nil {
		return fmt.Errorf("No Validator")
	}
	return binding.Validator.ValidateStruct(obj)
}
