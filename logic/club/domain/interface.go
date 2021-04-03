package domain

import clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"

type ICmdHandler interface {
	ReadParam(jsonBytes []byte) error
	SetSingleParamMode()
	ICmdLogic
}

type ICmdLogic interface {
	Do(text string) error
	Init(ICmdHandlerContext, func(requireRawParamAttr, requireRawParamAttrText string, isInputImmediately bool)) error
	GetSingleParam(attr string) string
	LoadSingleParam(attr, text string) error
	GetInputTemplate(requireRawParamAttr string) interface{}
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsComfirmed() bool
	GetRequireInputCmdText(cmd *TextCmd, attr, attrText string, isInputImmediately bool) (string, error)
	CacheParams() error
	GetCancelSignl() (string, error)
	GetComfirmSignl() (string, error)
	GetCancelInpuingSignl() (string, error)
}
