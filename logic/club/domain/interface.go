package domain

import clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"

type ICmdHandler interface {
	ReadParam(jsonBytes []byte) error
	SetSingleParamMode()
	ICmdLogic
}

type ICmdLogic interface {
	Do(text string) error
	Init(ICmdHandlerContext) error
	GetSingleParam(attr string) string
	LoadSingleParam(attr, text string) error
	GetInputTemplate(requireRawParamAttr string) interface{}
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsComfirmed() bool
	CacheParams() error
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
	GetSignal() (string, error)
}
