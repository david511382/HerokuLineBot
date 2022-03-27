package club

import (
	"heroku-line-bot/src/logic/club/domain"
	errUtil "heroku-line-bot/src/pkg/util/error"
)

type InputParam struct {
	RequireRawAttr string                 `json:"require_raw_attr"`
	IsCancel       bool                   `json:"is_cancel"`
	IsConfirm      bool                   `json:"is_comfirm,omitempty"`
	IsNotCache     bool                   `json:"is_not_cache"`
	param          domain.IParamTextValue `json:"-"`
}

func NewInputParam(param domain.IParamTextValue) *InputParam {
	return &InputParam{
		param: param,
	}
}

func (b *InputParam) GetHandler() (handler ParamHandler, isUpdateRequireAttr bool, warnMessage interface{}, resultErrInfo errUtil.IError) {
	originRequireAttr := b.RequireRawAttr
	requireAttr, message, errInfo := b.getRequireAttr()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if errInfo.IsError() {
			return
		}
	}

	warnMessage = message
	isUpdateRequireAttr = originRequireAttr != requireAttr
	handler = *NewParamHandler(requireAttr, b.param)
	return
}

func (b *InputParam) getRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	if b.RequireRawAttr == "" {
		attr, message, errInfo := b.param.GetRequireAttr()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
		warnMessage = message
		requireAttr = attr

		if _, _, isNotRequireChecking := b.param.GetRequireAttrInfo(attr); !isNotRequireChecking {
			b.RequireRawAttr = attr
		}
	} else {
		requireAttr = b.RequireRawAttr
	}
	return
}
