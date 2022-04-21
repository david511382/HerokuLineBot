package api

import (
	badmintonActivityLogic "heroku-line-bot/src/logic/badminton/activity"
	badmintonPlaceLogic "heroku-line-bot/src/logic/badminton/place"
	badmintonTeamLogic "heroku-line-bot/src/logic/badminton/team"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/server/domain/resp"
	"sort"
	"time"
)

var MockGetActivitys func(
	fromDate,
	toDate *util.DateTime,
	pageIndex,
	pageSize uint,
	placeIDs,
	teamIDs []int,
	everyWeekdays []time.Weekday,
) (
	result resp.GetActivitys,
	resultErrInfo errUtil.IError,
)

// pageIndex: 1開始
func GetActivitys(
	fromDate,
	toDate *util.DateTime,
	pageIndex,
	pageSize uint,
	placeIDs,
	teamIDs []int,
	everyWeekdays []time.Weekday,
) (
	result resp.GetActivitys,
	resultErrInfo errUtil.IError,
) {
	if MockGetActivitys != nil {
		return MockGetActivitys(fromDate, toDate, pageIndex, pageSize, placeIDs, teamIDs, everyWeekdays)
	}

	result.Activitys = make([]*resp.GetActivitysActivity, 0)

	activityMap := make(map[int]*resp.GetActivitysActivity)
	activityIDs := make([]int, 0)
	idPlaceMap := make(map[int]string)
	idTeamMap := make(map[int]string)
	{
		args, errInfo := badmintonActivityLogic.GetUnfinishedActiviysSqlReqs(
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
				dbDatas, err := database.Club().Activity.Select(
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
					activityMap[activityID].PeopleLimit = util.GetIntP(int(*v.PeopleLimit))
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
		dbDatas, err := database.Club().JoinActivityDetail(clubdb.ReqsClubJoinActivityDetail{
			Activity: &activity.Reqs{
				IDs: activityIDs,
			},
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
			court := &resp.GetActivitysCourt{
				FromTime: startTime,
				ToTime:   endTime,
				Count:    v.RentalCourtDetailCount,
			}
			activityMap[activityID].Courts = append(activityMap[activityID].Courts, court)
		}
	}

	{
		placeIDs := make([]int, 0)
		for id := range idPlaceMap {
			placeIDs = append(placeIDs, id)
		}

		dbDatas, errInfo := badmintonPlaceLogic.Load(placeIDs...)
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
		teamIDs := make([]int, 0)
		for id := range idTeamMap {
			teamIDs = append(teamIDs, id)
		}

		dbDatas, errInfo := badmintonTeamLogic.Load(teamIDs...)
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

	activityIDmemberIDNameMap := make(map[int]map[int]bool)
	memberIDNameMap := make(map[int]string)
	{
		activityIDs := make([]int, 0)
		for activityID, v := range activityMap {
			if v.IsShowMembers {
				activityIDs = append(activityIDs, activityID)
			}
		}

		if len(activityIDs) > 0 {
			dbDatas, err := database.Club().MemberActivity.Select(
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

			memberIDs := make([]int, 0)
			for _, v := range dbDatas {
				activityID := v.ActivityID
				memberID := v.MemberID

				if activityIDmemberIDNameMap[activityID] == nil {
					activityIDmemberIDNameMap[activityID] = make(map[int]bool)
				}
				activityIDmemberIDNameMap[activityID][memberID] = false

				memberIDs = append(memberIDs, v.MemberID)
			}

			if len(memberIDs) > 0 {
				dbDatas, err := database.Club().Member.Select(
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
