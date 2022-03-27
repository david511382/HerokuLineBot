package domain

import (
	clublinebotDomain "heroku-line-bot/src/logic/clublinebot/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/tidwall/gjson"
)

type ICmdHandler interface {
	ReadParam(jr gjson.Result) (resultErrInfo errUtil.IError)
	ICmdLogic
	CacheParams() (resultErrInfo errUtil.IError)
}

type ICmdLogic interface {
	IParamTextValue
	Do(text string) (resultErrInfo errUtil.IError)
	Init(ICmdHandlerContext) (resultErrInfo errUtil.IError)
}

type IParamTextValue interface {
	// 正在確認輸入資料時不會呼叫
	GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool)
	// 正在確認輸入資料時不會呼叫
	GetInputTemplate(attr string) (messages interface{})
	// 正在確認輸入資料時不會呼叫
	GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError)
	LoadRequireInputTextParam(rawAttr, textValue string) (resultErrInfo errUtil.IError)
}

type ICmdHandlerContext interface {
	clublinebotDomain.IContext
	IsConfirmed() bool
	CacheParams() (resultErrInfo errUtil.IError)
}
