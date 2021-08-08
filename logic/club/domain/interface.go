package domain

import (
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	errLogic "heroku-line-bot/logic/error"
)

type ICmdHandler interface {
	ReadParam(jsonBytes []byte) (resultErrInfo errLogic.IError)
	SetSingleParamMode()
	ICmdLogic
}

type ICmdLogic interface {
	Do(text string) (resultErrInfo errLogic.IError)
	Init(ICmdHandlerContext) (resultErrInfo errLogic.IError)
	GetSingleParam(attr string) string
	LoadSingleParam(attr, text string) (resultErrInfo errLogic.IError)
	GetInputTemplate(requireRawParamAttr string) interface{}
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsComfirmed() bool
	CacheParams() (resultErrInfo errLogic.IError)
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
	GetSignal() (string, errLogic.IError)
}
