package errorcode

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"net/http"
)

type ErrorMsg string

const (
	ERROR_MSG_SUCCESS     ErrorMsg = "完成"
	ERROR_MSG_NO_ACTIVITY ErrorMsg = "活動不存在"

	ERROR_MSG_REQUEST          ErrorMsg = "參數錯誤"
	ERROR_MSG_WRONG_PAY        ErrorMsg = "金額不符"
	ERROR_MSG_NO_DATES         ErrorMsg = "沒日期"
	ERROR_MSG_NO_DESPOSIT_DATE ErrorMsg = "沒訂金日期"
	ERROR_MSG_NO_BALANCE_DATE  ErrorMsg = "沒尾款日期"
	ERROR_MSG_WRONG_PLACE      ErrorMsg = "錯誤地點"
	ERROR_MSG_WRONG_TEAM       ErrorMsg = "錯誤隊伍"

	ERROR_MSG_AUTH ErrorMsg = "未登入"

	ERROR_MSG_FORBIDDEN ErrorMsg = "沒權限"

	ERROR_MSG_ERROR ErrorMsg = "發生錯誤"
)

func (em ErrorMsg) New(logErrInfos ...errUtil.IError) IErrorcode {
	var errInfo errUtil.IError
	if len(logErrInfos) > 0 {
		errInfo = logErrInfos[0]
	}

	code := http.StatusOK
	switch em {
	case ERROR_MSG_REQUEST,
		ERROR_MSG_NO_DATES,
		ERROR_MSG_NO_DESPOSIT_DATE,
		ERROR_MSG_NO_BALANCE_DATE,
		ERROR_MSG_WRONG_PLACE,
		ERROR_MSG_WRONG_TEAM:
		code = http.StatusBadRequest
	case ERROR_MSG_AUTH:
		code = http.StatusUnauthorized
	case ERROR_MSG_FORBIDDEN:
		code = http.StatusForbidden
	case ERROR_MSG_ERROR:
		code = http.StatusInternalServerError
	}

	return newErrorcode(em, errInfo, code)
}

func (em ErrorMsg) Error() string {
	return string(em)
}
