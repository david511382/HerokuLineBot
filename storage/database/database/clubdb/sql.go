package clubdb

import (
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/storage/database/database/clubdb/activity"
)

var MockJoinActivityDetail func(arg dbModel.ReqsClubJoinActivityDetail) (
	response []*dbModel.RespClubJoinActivityDetail,
	resultErr error,
)

func (d *Database) JoinActivityDetail(arg dbModel.ReqsClubJoinActivityDetail) (
	response []*dbModel.RespClubJoinActivityDetail,
	resultErr error,
) {
	if MockJoinActivityDetail != nil {
		return MockJoinActivityDetail(arg)
	}

	response = make([]*dbModel.RespClubJoinActivityDetail, 0)

	dp := d.GetSlave()

	activityModel := activity.Activity{}
	m := activityModel.GetTable()
	dp = dp.Model(m)

	if arg.ReqsClubActivity != nil {
		dp = activityModel.WhereArg(dp, *arg.ReqsClubActivity)
	}

	dp = dp.Select(
		`activity.id AS activity_id,
		rental_court_detail.start_time AS rental_court_detail_start_time,
		rental_court_detail.end_time AS rental_court_detail_end_time,
		rental_court_detail.count AS rental_court_detail_count`)
	dp = dp.Joins(
		`JOIN rental_court
		ON rental_court.date = activity.date
			AND rental_court.place_id = activity.place_id`)
	dp = dp.Joins(
		`JOIN rental_court_ledger_court
		ON rental_court_ledger_court.rental_court_id = rental_court.id
			AND rental_court_ledger_court.team_id = activity.team_id`)
	dp = dp.Joins(
		`JOIN rental_court_ledger
		ON rental_court_ledger.id = rental_court_ledger_court.rental_court_ledger_id`)
	dp = dp.Joins(
		`JOIN  rental_court_detail
		ON rental_court_detail.id = rental_court_ledger.rental_court_detail_id`)

	if err := dp.Scan(&response).Error; err != nil {
		resultErr = err
		return
	}

	return
}
