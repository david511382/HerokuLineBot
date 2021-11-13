package common

import "heroku-line-bot/storage/redis/domain"

func IsRedisError(err error) bool {
	if err == nil ||
		err.Error() == domain.NOT_CHANGE.Error() ||
		err.Error() == domain.NOT_EXIST.Error() {
		return false
	}

	return true
}
