package domain

import (
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	errUtil "heroku-line-bot/util/error"
)

type ICmdHandler interface {
	ReadParam(jsonBytes []byte) (resultErrInfo errUtil.IError)
	SetSingleParamMode()
	ICmdLogic
}

type ICmdLogic interface {
	Do(text string) (resultErrInfo errUtil.IError)
	Init(ICmdHandlerContext) (resultErrInfo errUtil.IError)
	GetSingleParam(attr string) string
	LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError)
	GetInputTemplate(requireRawParamAttr string) interface{}
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsComfirmed() bool
	CacheParams() (resultErrInfo errUtil.IError)
	ICmdHandlerSignal
	SetRequireInputMode(attr, attrText string, isInputImmediately bool)
}

type ICmdHandlerSignal interface {
	GetKeyValueInputMode(pathValueMap map[string]interface{}) ICmdHandlerSignal
	GetCancelMode() ICmdHandlerSignal
	GetComfirmMode() ICmdHandlerSignal
	GetCancelInputMode() ICmdHandlerSignal
	GetRequireInputMode(attr, attrText string, isInputImmediately bool) ICmdHandlerSignal
	GetCmdInputMode(cmdP *TextCmd) ICmdHandlerSignal
	GetDateTimeCmdInputMode(timeCmd DateTimeCmd, attr string) ICmdHandlerSignal
	GetSignal() (string, errUtil.IError)
}
