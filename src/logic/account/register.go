package account

import (
	dbModel "heroku-line-bot/src/model/database"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
)

func Registe(
	db *clubdb.Database,
	data *dbModel.ClubMember,
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
