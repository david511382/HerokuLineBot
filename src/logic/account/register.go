package account

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
)

func Registe(
	db *clubdb.Database,
	data *member.Model,
) (resultErrInfo errUtil.IError) {
	if db == nil {
		db = database.Club()
	}

	if err := db.Member.Insert(data); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return
}
