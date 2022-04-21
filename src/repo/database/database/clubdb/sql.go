package clubdb

import (
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
)

type ReqsClubJoinActivityDetail struct {
	Activity *activity.Reqs
}

type RespClubJoinActivityDetail struct {
	ActivityID                 int
	RentalCourtDetailStartTime string
	RentalCourtDetailEndTime   string
	RentalCourtDetailCount     int
}

var MockJoinActivityDetail func(arg ReqsClubJoinActivityDetail) (
	response []*RespClubJoinActivityDetail,
	resultErr error,
)

func (d *Database) JoinActivityDetail(arg ReqsClubJoinActivityDetail) (
	response []*RespClubJoinActivityDetail,
	resultErr error,
) {
	if MockJoinActivityDetail != nil {
		return MockJoinActivityDetail(arg)
	}

	response = make([]*RespClubJoinActivityDetail, 0)

	dp, err := d.GetSlave()
	if err != nil {
		resultErr = err
		return
	}

	dp = dp.Model(new(activity.Model))

	if arg.Activity != nil {
		dp = arg.Activity.WhereArg(dp)
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
