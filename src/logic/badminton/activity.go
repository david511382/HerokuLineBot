package badminton

import (
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"time"
)

type IBadmintonActivityLogic interface {
	GetUnfinishedActiviysSqlReqs(
		fromDate, toDate *util.DateTime,
		teamIDs,
		placeIDs []uint,
		everyWeekdays []time.Weekday,
	) (
		resultArgs []*activity.Reqs,
		resultErrInfo errUtil.IError,
	)

	GetActivityDetail(
		activityReqs *activity.Reqs,
	) (
		activityID_detailMap map[uint]*CourtDetail,
		resultErrInfo errUtil.IError,
	)
}

type BadmintonActivityLogic struct {
	clubDb *clubdb.Database
}

func NewBadmintonActivityLogic(clubDb *clubdb.Database) *BadmintonActivityLogic {
	return &BadmintonActivityLogic{
		clubDb: clubDb,
	}
}

func (l *BadmintonActivityLogic) GetUnfinishedActiviysSqlReqs(
	fromDate, toDate *util.DateTime,
	teamIDs,
	placeIDs []uint,
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
		copyArg.TeamID = util.PointerOf(teamID)
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
				maxDate, minTime, err := l.clubDb.Activity.MinMaxDate(noDateArg)
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

func (l *BadmintonActivityLogic) GetActivityDetail(
	activityReqs *activity.Reqs,
) (
	activityID_detailMap map[uint]*CourtDetail,
	resultErrInfo errUtil.IError,
) {
	activityID_detailMap = make(map[uint]*CourtDetail)

	dbDatas, err := l.clubDb.JoinActivityDetail(clubdb.ReqsClubJoinActivityDetail{
		Activity: activityReqs,
	})
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	for _, v := range dbDatas {
		activityID := v.ActivityID

		startTime, err := commonLogic.HourMinTime(v.RentalCourtDetailStartTime).Time()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		endTime, err := commonLogic.HourMinTime(v.RentalCourtDetailEndTime).Time()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		activityID_detailMap[activityID] = &CourtDetail{
			TimeRange: util.TimeRange{
				From: startTime,
				To:   endTime,
			},
			Count: uint8(v.RentalCourtDetailCount),
		}
	}

	return
}
