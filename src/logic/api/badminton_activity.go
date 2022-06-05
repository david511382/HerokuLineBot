package api

import (
	badmintonLogic "heroku-line-bot/src/logic/badminton"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/repo/redis/db/badminton"
	"heroku-line-bot/src/server/domain/resp"
	"sort"
	"time"
)

type IBadmintonActivityApiLogic interface {
	GetActivitys(
		fromDate,
		toDate *time.Time,
		pageIndex,
		pageSize uint,
		placeIDs,
		teamIDs []uint,
		everyWeekdays []time.Weekday,
	) (
		result resp.GetActivitys,
		resultErrInfo errUtil.IError,
	)
}

type BadmintonActivityApiLogic struct {
	clubDb                 *clubdb.Database
	badmintonTeamLogic     badmintonLogic.IBadmintonTeamLogic
	badmintonActivityLogic badmintonLogic.IBadmintonActivityLogic
	badmintonPlaceLogic    badmintonLogic.IBadmintonPlaceLogic
}

func NewBadmintonActivityApiLogic(
	clubDb *clubdb.Database,
	badmintonRds *badminton.Database,
	badmintonTeamLogic badmintonLogic.IBadmintonTeamLogic,
	badmintonActivityLogic badmintonLogic.IBadmintonActivityLogic,
	badmintonPlaceLogic badmintonLogic.IBadmintonPlaceLogic,
) *BadmintonActivityApiLogic {
	return &BadmintonActivityApiLogic{
		clubDb:                 clubDb,
		badmintonTeamLogic:     badmintonTeamLogic,
		badmintonActivityLogic: badmintonActivityLogic,
		badmintonPlaceLogic:    badmintonPlaceLogic,
	}
}

// pageIndex: 1開始
func (l *BadmintonActivityApiLogic) GetActivitys(
	fromDate,
	toDate *time.Time,
	pageIndex,
	pageSize uint,
	placeIDs,
	teamIDs []uint,
	everyWeekdays []time.Weekday,
) (
	result resp.GetActivitys,
	resultErrInfo errUtil.IError,
) {
	result.Activitys = make([]*resp.GetActivitysActivity, 0)

	activityMap := make(map[uint]*resp.GetActivitysActivity)
	activityIDs := make([]uint, 0)
	idPlaceMap := make(map[uint]string)
	idTeamMap := make(map[uint]string)
	{
		args, errInfo := l.badmintonActivityLogic.GetUnfinishedActiviysSqlReqs(
			fromDate, toDate,
			teamIDs, placeIDs, everyWeekdays,
		)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		dbActivityDatas := make([]*activity.Model, 0)
		for _, arg := range args {
			{
				dbDatas, err := l.clubDb.Activity.Select(
					*arg,
					activity.COLUMN_ID,
					activity.COLUMN_Date,
					activity.COLUMN_PlaceID,
					activity.COLUMN_TeamID,
					activity.COLUMN_CourtsAndTime,
					activity.COLUMN_ClubSubsidy,
					activity.COLUMN_Description,
					activity.COLUMN_PeopleLimit,
				)
				if err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}

				dbActivityDatas = append(dbActivityDatas, dbDatas...)
			}
		}

		sort.SliceStable(dbActivityDatas, func(i, j int) bool {
			return dbActivityDatas[i].ID < dbActivityDatas[j].ID
		})
		sort.SliceStable(dbActivityDatas, func(i, j int) bool {
			return dbActivityDatas[i].Date.Before(dbActivityDatas[j].Date)
		})
		result.DataCount = len(dbActivityDatas)
		from, before := commonLogic.PageSlice(result.DataCount, pageSize, pageIndex)
		if from > -1 {
			for _, v := range dbActivityDatas[from:before] {
				activityID := v.ID
				placeID := v.PlaceID
				teamID := v.TeamID
				date := v.Date

				activityMap[activityID] = &resp.GetActivitysActivity{
					ActivityID: activityID,
					PlaceID:    placeID,
					TeamID:     teamID,
					Date:       date,

					Courts:        make([]*resp.GetActivitysCourt, 0),
					Members:       make([]*resp.GetActivitysMember, 0),
					IsShowMembers: true,
				}
				if v.Description != "" {
					activityMap[activityID].Description = &v.Description
				}
				if v.PeopleLimit != nil {
					activityMap[activityID].PeopleLimit = util.PointerOf(int(*v.PeopleLimit))
				}

				activityIDs = append(activityIDs, activityID)

				idPlaceMap[placeID] = ""
				idTeamMap[teamID] = ""
			}
		}
	}

	if len(activityIDs) == 0 {
		return
	}

	{
		activityID_detailsMap := make(map[uint][]*badmintonLogic.CourtDetail)
		_, errInfo := l.badmintonActivityLogic.GetActivityDetail(
			&activity.Reqs{
				IDs: activityIDs,
			},
			activityID_detailsMap,
		).Run()
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if errInfo.IsError() {
				return
			}
		}

		for activityID, details := range activityID_detailsMap {
			for _, detail := range details {
				activityMap[activityID].Courts = append(activityMap[activityID].Courts, &resp.GetActivitysCourt{
					FromTime: detail.From,
					ToTime:   detail.To,
					Count:    int(detail.Count),
				})
			}
		}
	}

	{
		placeIDs := make([]uint, 0)
		for id := range idPlaceMap {
			placeIDs = append(placeIDs, id)
		}

		dbDatas, errInfo := l.badmintonPlaceLogic.Load(placeIDs...)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		for id, v := range dbDatas {
			idPlaceMap[id] = v.Name
		}
	}

	{
		teamIDs := make([]uint, 0)
		for id := range idTeamMap {
			teamIDs = append(teamIDs, id)
		}

		dbDatas, errInfo := l.badmintonTeamLogic.Load(teamIDs...)
		if errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		for id, v := range dbDatas {
			idTeamMap[id] = v.Name
		}
	}

	activityIDmemberIDNameMap := make(map[uint]map[uint]bool)
	memberIDNameMap := make(map[uint]string)
	{
		activityIDs := make([]uint, 0)
		for activityID, v := range activityMap {
			if v.IsShowMembers {
				activityIDs = append(activityIDs, activityID)
			}
		}

		if len(activityIDs) > 0 {
			dbDatas, err := l.clubDb.MemberActivity.Select(
				memberactivity.Reqs{
					ActivityIDs: activityIDs,
				},
				memberactivity.COLUMN_MemberID,
				memberactivity.COLUMN_ActivityID,
			)
			if err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}

			memberIDs := make([]uint, 0)
			for _, v := range dbDatas {
				activityID := v.ActivityID
				memberID := v.MemberID

				if activityIDmemberIDNameMap[activityID] == nil {
					activityIDmemberIDNameMap[activityID] = make(map[uint]bool)
				}
				activityIDmemberIDNameMap[activityID][memberID] = false

				memberIDs = append(memberIDs, v.MemberID)
			}

			if len(memberIDs) > 0 {
				dbDatas, err := l.clubDb.Member.Select(
					member.Reqs{
						IDs: memberIDs,
					},
					member.COLUMN_ID,
					member.COLUMN_Name,
				)
				if err != nil {
					errInfo := errUtil.NewError(err)
					resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
					return
				}

				for _, v := range dbDatas {
					memberID := v.ID
					name := v.Name

					memberIDNameMap[memberID] = name
				}
			}
		}
	}
	for _, activityID := range activityIDs {
		v := activityMap[activityID]
		placeID := v.PlaceID
		teamID := v.TeamID

		if v.IsShowMembers {
			for memberID := range activityIDmemberIDNameMap[activityID] {
				name := memberIDNameMap[memberID]
				v.Members = append(v.Members, &resp.GetActivitysMember{
					ID:   memberID,
					Name: name,
				})
			}
		}

		v.PlaceName = idPlaceMap[placeID]
		v.TeamName = idTeamMap[teamID]

		result.Activitys = append(result.Activitys, v)
	}

	return
}
