package domain

import clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"

type ICmdHandler interface {
	ReadParam(jsonBytes []byte) error
	SetSingleParamMode()
	ICmdLogic
}

type ICmdLogic interface {
	Do(text string) error
	Init(ICmdHandlerContext)
	GetSingleParam(attr string) string
	LoadSingleParam(attr, text string) (resultValue interface{}, resultErr error)
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsComfirmed() bool
	GetRequireInputCmdText(cmd *TextCmd, attr, attrText string, isInputImmediately bool) (string, error)
	CacheParams() error
	GetCancelSignl() string
	GetComfirmSignl() string
}
