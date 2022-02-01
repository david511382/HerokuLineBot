package common

func IsRedisError(err error) bool {
	if err == nil ||
		err.Error() == ERROR_MSG_NOT_CHANGE ||
		err.Error() == ERROR_MSG_NOT_EXIST {
		return false
	}

	return true
}
