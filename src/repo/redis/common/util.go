package common

import errUtil "heroku-line-bot/src/pkg/util/error"

func IsContainNotChange(errInfo errUtil.IError) bool {
	if errInfo == nil {
		return false
	}

	for _, err := range errUtil.Split(errInfo) {
		errInfo, ok := err.(errUtil.IError)
		if ok {
			if errUtil.Equal(errInfo, NotChangeErrInfo) {
				return true
			}
			continue
		}

		if err.Error() == ERROR_MSG_NOT_CHANGE {
			return true
		}
	}

	return false
}
