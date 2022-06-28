package badminton

import (
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/pkg/util/flow"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"time"
)

type IBadmintonActivityLogic interface {
	GetUnfinishedActiviysSqlReqs(
		fromDate, toDate *time.Time,
		teamIDs,
		placeIDs []uint,
		everyWeekdays []time.Weekday,
	) (
		resultArgs []*activity.Reqs,
		resultErrInfo errUtil.IError,
	)

	GetActivityDetail(
		activityReqs *activity.Reqs,
		respActivityID_detailsMap map[uint][]*CourtDetail,
	) flow.IStep
}

type BadmintonActivityLogic struct {
	clubDb *clubdb.Database
}

func NewBadmintonActivityLogic(
	clubDb *clubdb.Database,
) *BadmintonActivityLogic {
	return &BadmintonActivityLogic{
		clubDb: clubDb,
	}
}

func (l *BadmintonActivityLogic) GetUnfinishedActiviysSqlReqs(
	fromDate, toDate *time.Time,
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
			var startDate, endtDate *util.DefinedTime[util.DateInt]
			if fromDate != nil {
				startDate = util.PointerOf(util.Date().Of(*fromDate))
			}
			if toDate != nil {
				endtDate = util.PointerOf(util.Date().Of(*toDate))
			}
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
					startDate = util.PointerOf(util.Date().Of(minTime))
				}
				if endtDate == nil {
					endtDate = util.PointerOf(util.Date().Of(maxDate))
				}
			}

			arg.FromDate = nil
			arg.ToDate = nil
			arg.Dates = make([]*time.Time, 0)
			dates := util.GetDatesInWeekdays(*startDate, *endtDate, everyWeekdays...)
			for _, v := range dates {
				arg.Dates = append(arg.Dates, util.PointerOf(v.Time()))
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
	respActivityID_detailsMap map[uint][]*CourtDetail,
) flow.IStep {
	var (
		activity_courtsMap = make(map[uint]ActivityCourts)
	)
	return flow.Flow("GetActivityDetail",
		flow.Step{
			StepName: "讀資料庫資料",
			Fun: func() (resultErrInfo errUtil.IError) {
				dbReqs := activity.Reqs{}
				if activityReqs != nil {
					dbReqs = *activityReqs
				}
				dbDatas, err := l.clubDb.Activity.Select(
					dbReqs,
					activity.COLUMN_ID,
					activity.COLUMN_CourtsAndTime,
				)
				if err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				} else if len(dbDatas) == 0 {
					return flow.ErrorBreak
				}

				for _, v := range dbDatas {
					courts, errInfo := DbActivityCourtsStr(v.CourtsAndTime).ParseCourts()
					if errInfo != nil {
						resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
						if resultErrInfo.IsError() {
							return
						}
					}
					activity_courtsMap[v.ID] = courts
				}
				return
			},
		},
		flow.Step{
			StepName: "整合返回資料",
			Fun: func() (resultErrInfo errUtil.IError) {
				for activityID, courts := range activity_courtsMap {
					for _, court := range courts {
						respActivityID_detailsMap[activityID] = append(respActivityID_detailsMap[activityID], &CourtDetail{
							TimeRange: util.TimeRange{
								From: court.FromTime,
								To:   court.ToTime,
							},
							Count: court.Count,
						})
					}
				}
				return
			},
		},
	)
}
