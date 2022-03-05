package errorcode

import errUtil "heroku-line-bot/src/pkg/util/error"

type ErrorMsg string

const (
	ERROR_MSG_EMPTY            ErrorMsg = ""
	ERROR_MSG_WRONG_PAY        ErrorMsg = "金額不符"
	ERROR_MSG_NO_DATES         ErrorMsg = "沒日期"
	ERROR_MSG_NO_DESPOSIT_DATE ErrorMsg = "沒訂金日期"
	ERROR_MSG_NO_BALANCE_DATE  ErrorMsg = "沒尾款日期"
	ERROR_MSG_WRONG_PLACE      ErrorMsg = "錯誤地點"
	ERROR_MSG_WRONG_TEAM       ErrorMsg = "錯誤隊伍"
)

func (em ErrorMsg) New() errUtil.IError {
	return newErrorcode(em)
}

func (ec ErrorMsg) Equal(err error) bool {
	return err.Error() == string(ec)
}

func GetErrorMsg(err error) ErrorMsg {
	if err == nil {
		return ""
	}

	errs := errUtil.Split(err)
	err = errs[len(errs)-1]
	return ErrorMsg(err.Error())
}
