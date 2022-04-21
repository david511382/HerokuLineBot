package activity

import (
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"time"
)

func GetUnfinishedActiviysSqlReqs(
	fromDate, toDate *util.DateTime,
	teamIDs,
	placeIDs []int,
	everyWeekdays []time.Weekday,
) (
	resultArgs []*activity.Reqs,
	resultErrInfo errUtil.IError,
) {
	resultArgs = make([]*activity.Reqs, 0)

	args := make([]*activity.Reqs, 0)
	arg := &activity.Reqs{
		PlaceIDs: placeIDs,
	}
	if fromDate != nil {
		arg.FromDate = fromDate.TimeP()
	}
	if toDate != nil {
		arg.ToDate = toDate.TimeP()
	}
	if isNotSpecifyingTeam := len(teamIDs) == 0; isNotSpecifyingTeam {
		args = append(args, arg)
	}
	for _, teamID := range teamIDs {
		copyArg := *arg
		copyArg.TeamID = util.GetIntP(teamID)
		args = append(args, &copyArg)
	}

	for _, arg := range args {
		if len(everyWeekdays) > 0 {
			startDate := fromDate
			endtDate := toDate
			if startDate == nil || endtDate == nil {
				noDateArg := *arg
				noDateArg.FromDate = nil
				noDateArg.ToDate = nil
				maxDate, minTime, err := database.Club().Activity.MinMaxDate(noDateArg)
				if err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}

				if startDate == nil {
					startDate = util.NewDateTimePOf(&minTime)
				}
				if endtDate == nil {
					endtDate = util.NewDateTimePOf(&maxDate)
				}
			}

			arg.FromDate = nil
			arg.ToDate = nil
			arg.Dates = make([]*time.Time, 0)
			dates := util.GetDatesInWeekdays(*startDate, *endtDate, everyWeekdays...)
			for _, v := range dates {
				arg.Dates = append(arg.Dates, v.TimeP())
			}

			if isEmpty := len(arg.Dates) == 0; isEmpty {
				continue
			}
		}

		resultArgs = append(resultArgs, arg)
	}

	return
}
