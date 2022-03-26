package club

import (
	"encoding/json"
	"errors"
	"fmt"
	"heroku-line-bot/src/logic/club/domain"
	clublinebotDomain "heroku-line-bot/src/logic/clublinebot/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type CmdHandler struct {
	*domain.CmdBase
	*domain.TimePostbackParams
	clublinebotDomain.IContext `json:"-"`
	domain.ICmdLogic
	InputParam
	pathValueMap map[string]interface{}
}

func (b *CmdHandler) ReadParam(textJsonResult gjson.Result) errUtil.IError {
	js, err := json.Marshal(b)
	if err != nil {
		return errUtil.NewError(err)
	}

	jsStr := string(js)
	for path, jr := range textJsonResult.Map() {
		s, err := sjson.Set(jsStr, path, jr.Value())
		if err != nil {
			return errUtil.NewError(err)
		}
		jsStr = s
	}

	if err := json.Unmarshal([]byte(jsStr), b); err != nil {
		return errUtil.NewError(err)
	}

	if errInfo := b.CacheParams(); errInfo != nil {
		return errInfo
	}
	return nil
}

func (b *CmdHandler) CacheParams() (resultErrInfo errUtil.IError) {
	if b.IsNotCache {
		return
	}

	if jsBytes, err := json.Marshal(b); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		js := string(jsBytes)
		if err := b.SaveParam(js); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
	}
	return nil
}

func (b *CmdHandler) IsConfirmed() bool {
	return b.IsConfirm
}

func (b *CmdHandler) Do(text string) (resultErrInfo errUtil.IError) {
	if b.InputParam.IsCancel {
		if err := b.DeleteParam(); err != nil {
			return errUtil.NewError(err)
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("取消"),
		}
		if err := b.IContext.Reply(replyMessges); err != nil {
			return errUtil.NewError(err)
		}

		return nil
	}

	{
		paramHandler, isUpdateRequireAttr, errInfo := b.InputParam.GetHandler()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		if paramHandler.IsReading() {
			// 是否正在讀取輸入資料
			if text != "" {
				// 讀取輸入的資料到指定欄位
				if errInfo := paramHandler.Read(text); errInfo != nil {
					msg := fmt.Sprintf("參數格式錯誤:%s", errInfo.Error())
					replyMessges := []interface{}{
						linebot.GetTextMessage(msg),
					}
					if err := b.Reply(replyMessges); err != nil {
						errInfo := errUtil.NewError(err)
						resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
						return
					}
					return
				}

				paramHandler, _, errInfo = b.InputParam.GetHandler()
				if errInfo != nil {
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					if resultErrInfo.IsError() {
						return
					}
				}

				if errInfo := b.CacheParams(); errInfo != nil {
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					if errInfo.IsError() {
						return
					}
				}
			}

			// 判斷是否需要顯示輸入資料的介面
			if showMessge := paramHandler.GetInputTemplate(); showMessge != nil {
				replyMessges := []interface{}{
					showMessge,
				}
				if err := b.IContext.Reply(replyMessges); err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}

				return
			}
		} else if 不會讀取輸入資料並儲存資料 := text == ""; isUpdateRequireAttr && 不會讀取輸入資料並儲存資料 {
			if errInfo := b.CacheParams(); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if errInfo.IsError() {
					return
				}
			}
		}
	}

	if errInfo := b.ICmdLogic.Do(text); errors.Is(errInfo, domain.USER_NOT_REGISTERED) ||
		errors.Is(errInfo, domain.NO_AUTH_ERROR) {
		replyMessges := []interface{}{
			linebot.GetTextMessage(errInfo.Error()),
		}
		if err := b.IContext.Reply(replyMessges); err != nil {
			return errUtil.NewError(err)
		}
	} else {
		return errInfo
	}

	return nil
}
