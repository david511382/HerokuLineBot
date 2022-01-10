package club

import (
	"fmt"
	badmintonPlaceLogic "heroku-line-bot/logic/badminton/place"
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogic "heroku-line-bot/logic/club/lineuser"
	clubLineuserLogicDomain "heroku-line-bot/logic/club/lineuser/domain"
	dbModel "heroku-line-bot/model/database"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	linebotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/activity"
	"heroku-line-bot/storage/database/database/clubdb/member"
	"heroku-line-bot/storage/database/database/clubdb/memberactivity"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"sort"
	"strconv"
	"time"
)

type GetActivities struct {
	context               domain.ICmdHandlerContext `json:"-"`
	activities            []*getActivitiesActivity
	JoinActivityID        int                           `json:"join_activity_id"`
	TeamID                int                           `json:"team_id"`
	LeaveActivityID       int                           `json:"leave_activity_id"`
	currentUser           clubLineuserLogicDomain.Model `json:"-"`
	ListMembersActivityID int                           `json:"list_members_activity_id"`
}

type getActivitiesActivity struct {
	NewActivity
	JoinedMembers []*getActivitiesActivityJoinedMembers `json:"joined_members"`
	ActivityID    int                                   `json:"activity_id"`
}

type getActivitiesActivityJoinedMembers struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b *GetActivities) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = GetActivities{
		context:    context,
		activities: make([]*getActivitiesActivity, 0),
		TeamID:     clubTeamID,
	}

	return nil
}

func (b *GetActivities) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *GetActivities) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *GetActivities) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *GetActivities) init() (resultErrInfo errUtil.IError) {
	context := b.context

	activitys := make([]*dbModel.ClubActivity, 0)
	{
		dbDatas, err := database.Club.Activity.Select(
			dbModel.ReqsClubActivity{
				TeamID: &b.TeamID,
			},
			activity.COLUMN_ID,
			activity.COLUMN_Date,
			activity.COLUMN_PlaceID,
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

		activitys = append(activitys, dbDatas...)
	}

	if len(activitys) == 0 {
		return
	}

	activityIDs := []int{}
	idPlaceMap := make(map[int]string)
	for _, v := range activitys {
		activityIDs = append(activityIDs, v.ID)
		idPlaceMap[v.PlaceID] = ""
	}

	placeIDs := make([]int, 0)
	for id := range idPlaceMap {
		placeIDs = append(placeIDs, id)
	}
	if dbDatas, errInfo := badmintonPlaceLogic.Load(placeIDs...); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else {
		for id, v := range dbDatas {
			idPlaceMap[id] = v.Name
		}
	}

	activityIDMap := make(map[int][]*getActivitiesActivityJoinedMembers)
	memberActivityArg := dbModel.ReqsClubMemberActivity{
		ActivityIDs: activityIDs,
	}
	if dbDatas, err := database.Club.MemberActivity.Select(
		memberActivityArg,
		memberactivity.COLUMN_ID,
		memberactivity.COLUMN_MemberID,
		memberactivity.COLUMN_ActivityID,
	); err != nil {
		errInfo := errUtil.NewError(err)
		if resultErrInfo == nil {
			resultErrInfo = errInfo
		} else {
			resultErrInfo = resultErrInfo.Append(errInfo)
		}
		return
	} else if len(dbDatas) > 0 {
		memberIDs := make([]int, 0)
		for _, v := range dbDatas {
			memberIDs = append(memberIDs, v.MemberID)
		}
		type lineName struct {
			lineID *string
			name   string
		}
		memberIDNameMap := make(map[int]lineName)
		memberArg := dbModel.ReqsClubMember{
			IDs: memberIDs,
		}
		if dbDatas, err := database.Club.Member.Select(
			memberArg,
			member.COLUMN_ID,
			member.COLUMN_Name,
			member.COLUMN_LineID,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			for _, v := range dbDatas {
				memberIDNameMap[v.ID] = lineName{
					lineID: v.LineID,
					name:   v.Name,
				}
			}
		}

		sort.Slice(dbDatas, func(i, j int) bool {
			return dbDatas[i].ID < dbDatas[j].ID
		})
		for _, v := range dbDatas {
			activityID := v.ActivityID
			if activityIDMap[activityID] == nil {
				activityIDMap[activityID] = make([]*getActivitiesActivityJoinedMembers, 0)
			}
			activityIDMap[activityID] = append(activityIDMap[activityID], &getActivitiesActivityJoinedMembers{
				ID:   v.MemberID,
				Name: memberIDNameMap[v.MemberID].name,
			})
		}
	}

	sort.Slice(activitys, func(i, j int) bool {
		return activitys[i].Date.Before(activitys[j].Date)
	})
	for _, v := range activitys {
		activity := &getActivitiesActivity{
			NewActivity: NewActivity{
				Context:     context,
				Date:        util.DateTime(v.Date),
				PlaceID:     v.PlaceID,
				Description: v.Description,
				PeopleLimit: v.PeopleLimit,
				ClubSubsidy: v.ClubSubsidy,
			},
			JoinedMembers: activityIDMap[v.ID],
			ActivityID:    v.ID,
		}
		if errInfo := activity.ParseCourts(v.CourtsAndTime); errInfo != nil {
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
			return
		}
		b.activities = append(b.activities, activity)
	}

	return
}

func (b *GetActivities) listMembers() (resultErrInfo errUtil.IError) {
	var date time.Time
	var place string
	var peopleLimit *int16
	arg := dbModel.ReqsClubActivity{
		ID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.Activity.Select(
		arg,
		activity.COLUMN_Date,
		activity.COLUMN_PlaceID,
		activity.COLUMN_PeopleLimit,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		replyMessges := []interface{}{
			linebot.GetTextMessage("查無活動"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		return nil
	} else {
		v := dbDatas[0]
		date = v.Date
		peopleLimit = v.PeopleLimit

		if dbDatas, errInfo := badmintonPlaceLogic.Load(v.PlaceID); errInfo != nil {
			errInfo := errUtil.NewError(err)
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
		} else {
			for _, v := range dbDatas {
				place = v.Name
			}
		}
	}

	memberComponents := []interface{}{}
	memberActivityArg := dbModel.ReqsClubMemberActivity{
		ActivityID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.MemberActivity.Select(
		memberActivityArg,
		memberactivity.COLUMN_ID,
		memberactivity.COLUMN_MemberID,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		memberIDs := make([]int, 0)
		for _, v := range dbDatas {
			memberIDs = append(memberIDs, v.MemberID)
		}
		memberIDNameMap := make(map[int]string)
		memberArg := dbModel.ReqsClubMember{
			IDs: memberIDs,
		}
		if dbDatas, err := database.Club.Member.Select(
			memberArg,
			member.COLUMN_ID,
			member.COLUMN_Name,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			for _, v := range dbDatas {
				memberIDNameMap[v.ID] = v.Name
			}
		}

		keyValueEditComponentOption := &domain.KeyValueEditComponentOption{
			Indent: util.GetIntP(1),
		}
		sort.Slice(dbDatas, func(i, j int) bool {
			return dbDatas[i].ID < dbDatas[j].ID
		})
		for index, v := range dbDatas {
			id := index + 1
			idStr := strconv.Itoa(id)
			memberComponents = append(memberComponents,
				GetKeyValueEditComponent(
					idStr,
					memberIDNameMap[v.MemberID],
					keyValueEditComponentOption,
				),
			)
		}
	}

	contents := []interface{}{}

	contents = append(contents,
		GetKeyValueEditComponent(
			"日期",
			date.Format(util.DATE_FORMAT),
			nil,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"地點",
			place,
			nil,
		),
	)

	isUsingPeopleLimit := peopleLimit != nil
	joinedCount, _ := getJoinCount(len(memberComponents), peopleLimit)
	contents = append(contents, linebot.GetFlexMessageTextComponent("參加人員:", nil))
	contents = append(contents, memberComponents[:joinedCount]...)
	if isUsingPeopleLimit {
		contents = append(contents, linebot.GetFlexMessageTextComponent("候補人員:", nil))
		contents = append(contents, memberComponents[joinedCount:]...)
	}

	replyMessge := linebot.GetFlexMessage(
		"查看活動",
		linebot.GetFlexMessageBubbleContent(
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				nil,
				contents...,
			),
			nil,
		),
	)
	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *GetActivities) joinActivity() (resultErrInfo errUtil.IError) {
	userData := b.currentUser
	activityID := b.JoinActivityID
	uID := userData.ID

	insertData := &dbModel.ClubMemberActivity{
		ActivityID: activityID,
		MemberID:   uID,
		IsAttend:   false,
	}
	if err := database.Club.MemberActivity.Insert(insertData); err != nil && !database.IsUniqErr(err) {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *GetActivities) leaveActivity() (resultErrInfo errUtil.IError) {
	userData := b.currentUser

	var peopleLimit *int16
	activityPlace := ""
	var activityDate *time.Time
	activityArg := dbModel.ReqsClubActivity{
		ID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.Activity.Select(
		activityArg,
		activity.COLUMN_Date,
		activity.COLUMN_PlaceID,
		activity.COLUMN_PeopleLimit,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		peopleLimit = v.PeopleLimit
		activityDate = &v.Date

		if dbDatas, errInfo := badmintonPlaceLogic.Load(v.PlaceID); errInfo != nil {
			errInfo := errUtil.NewError(err)
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
		} else {
			for _, v := range dbDatas {
				activityPlace = v.Name
			}
		}
	}

	var notifyWaitingMemberID *int
	deleteMemberActivityID := 0
	arg := dbModel.ReqsClubMemberActivity{
		ActivityID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.MemberActivity.Select(
		arg,
		memberactivity.COLUMN_ID,
		memberactivity.COLUMN_MemberID,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		sort.Slice(dbDatas, func(i, j int) bool {
			return dbDatas[i].ID < dbDatas[j].ID
		})

		for i, v := range dbDatas {
			if v.MemberID == userData.ID {
				deleteMemberActivityID = v.ID

				if peopleLimit != nil {
					limitCount := int(*peopleLimit)
					if len(dbDatas) > limitCount && i < limitCount {
						notifyWaitingMemberID = &dbDatas[limitCount].MemberID
					}
				}

				break
			}
		}
	}

	arg = dbModel.ReqsClubMemberActivity{
		ID: &deleteMemberActivityID,
	}
	if err := database.Club.MemberActivity.Delete(arg); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	if isNotifyWaitingPerson := notifyWaitingMemberID != nil; isNotifyWaitingPerson {
		memberArg := dbModel.ReqsClubMember{
			ID: notifyWaitingMemberID,
		}
		if dbDatas, err := database.Club.Member.Select(
			memberArg,
			member.COLUMN_LineID,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if len(dbDatas) > 0 {
			lineID := dbDatas[0].LineID
			if lineID == nil {
				return nil
			}

			pushParam := &linebotReqs.PushMessage{
				To: *lineID,
				Messages: []interface{}{
					linebot.GetFlexMessage(
						"活動正取通知",
						linebot.GetFlexMessageBubbleContent(
							linebot.GetFlexMessageBoxComponent(
								linebotDomain.VERTICAL_MESSAGE_LAYOUT,
								nil,
								linebot.GetTextMessage("你已排上活動正取!!"),
								linebot.GetTextMessage(
									fmt.Sprintf("活動 %s(%s) %s",
										activityDate.Format(util.MONTH_DATE_SLASH_FORMAT),
										util.GetWeekDayName(activityDate.Weekday()),
										activityPlace,
									),
								),
								linebot.GetTextMessage("若無法參加麻煩要退出活動喔~"),
							),
							nil,
						),
					),
				},
			}
			if _, err := b.context.GetBot().PushMessage(pushParam); err != nil {
				if err := b.context.PushAdmin(
					[]interface{}{
						linebot.GetTextMessage(fmt.Sprintf("leaveActivity notifyLineID:%s, %s", *lineID, err.Error())),
					},
				); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				}
			}
		}
	}

	return nil
}

func (b *GetActivities) loadCurrentUserID() (replyMsg *string, resultErrInfo errUtil.IError) {
	lineID := b.context.GetUserID()
	userData, err := clubLineuserLogic.Get(lineID)
	if err != nil {
		errInfo := errUtil.NewError(err)
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if userData == nil {
		replyMsg = util.GetStringP(domain.USER_NOT_REGISTERED.Error())
		return
	}
	// TODO get member teamIDs
	b.currentUser = *userData

	return
}

func (b *GetActivities) isJoined(activity *getActivitiesActivity) bool {
	for _, v := range activity.JoinedMembers {
		mID := v.ID
		if b.currentUser.ID == mID {
			return true
		}
	}
	return false
}

func (b *GetActivities) Do(text string) (resultErrInfo errUtil.IError) {
	if replyMsg, errInfo := b.loadCurrentUserID(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else if replyMsg != nil {
		replyMessges := []interface{}{
			linebot.GetTextMessage(*replyMsg),
		}
		if replyErr := b.context.Reply(replyMessges); replyErr != nil {
			errInfo := errUtil.Newf("replyErr:%s", replyErr.Error())
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		return
	}

	if isListMembers := b.ListMembersActivityID > 0; isListMembers {
		if errInfo := b.listMembers(); errInfo != nil {
			replyMessges := []interface{}{
				linebot.GetTextMessage("查看人員發生錯誤，已通知管理員"),
			}
			if replyErr := b.context.Reply(replyMessges); replyErr != nil {
				errInfo := errUtil.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}

			resultErrInfo = errInfo
			return
		}
		return nil
	} else if isJoin := b.JoinActivityID > 0; isJoin {
		activityID := b.JoinActivityID
		arg := dbModel.ReqsClubActivity{
			ID: &activityID,
		}
		if count, err := database.Club.Activity.Count(arg); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if isActivityOpen := count > 0; isActivityOpen {
			if errInfo := b.joinActivity(); errInfo != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("參加發生錯誤，已通知管理員"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					resultErrInfo = errUtil.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
					return
				}

				resultErrInfo = errInfo
				return
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("活動已關閉"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}
			return nil
		}
	} else if isLeave := b.LeaveActivityID > 0; isLeave {
		activityID := b.LeaveActivityID
		arg := dbModel.ReqsClubActivity{
			ID: &activityID,
		}
		if count, err := database.Club.Activity.Count(arg); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else if isActivityOpen := count > 0; isActivityOpen {
			if errInfo := b.leaveActivity(); errInfo != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("退出發生錯誤，已通知管理員"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					resultErrInfo = errUtil.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
					return
				}

				resultErrInfo = errInfo
				return
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("活動已關閉"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}
			return nil
		}
	}

	replyMessge, err := b.GetActivitiesMessage("查看活動", true, true)
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *GetActivities) GetActivitiesMessage(altText string, isShowCurrentMember, isShowActionButton bool) (replyMessge interface{}, err error) {
	if err := b.init(); err != nil {
		return nil, err
	}

	carouselContents := []*linebotModel.FlexMessagBubbleComponent{}
	for _, activity := range b.activities {
		carouselContent, err := b.GetActivitieMessage(activity, isShowCurrentMember, isShowActionButton)
		if err != nil {
			return nil, err
		}

		carouselContents = append(carouselContents, carouselContent)
	}

	if len(carouselContents) == 0 {
		replyMessge = linebot.GetTextMessage("沒有活動")
	} else {
		replyMessge = linebot.GetFlexMessage(
			altText,
			linebot.GetFlexMessageCarouselContent(carouselContents...),
		)
	}

	return
}

func (b *GetActivities) GetActivitieMessage(
	activity *getActivitiesActivity,
	isShowCurrentMember, isShowActionButton bool,
) (
	carouselContent *linebotModel.FlexMessagBubbleComponent,
	resultErrInfo errUtil.IError,
) {
	contents := []interface{}{
		linebot.GetFlexMessageTextComponent(
			"活動資訊",
			&linebotModel.FlexMessageTextComponentOption{
				Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
				Color:  "#1DB446",
			},
		),
	}
	contents = append(contents, activity.getPlaceTimeTemplate()...)

	boxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
	)
	boxComponent.Contents = append(boxComponent.Contents, b.GetActivitieInfoContents(&activity.NewActivity)...)
	boxComponent.Contents = append(boxComponent.Contents, activity.getCourtsContents()...)
	boxComponent.Contents = append(boxComponent.Contents, b.GetActivitieEstimateContents(activity, isShowCurrentMember)...)
	contents = append(contents, boxComponent)

	footerContents := make([]interface{}, 0)
	if isShowActionButton {
		if len(activity.JoinedMembers) > 0 {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.list_members_activity_id": activity.ActivityID,
			}
			if js, errInfo := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); errInfo != nil {
				resultErrInfo = errInfo
				return
			} else {
				action := linebot.GetPostBackAction(
					"查看人員",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: "#A9A9A9",
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: domain.WHITE_COLOR,
							},
						),
					),
				)
			}
		}

		if b.isJoined(activity) {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.leave_activity_id": activity.ActivityID,
			}
			if js, errInfo := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); errInfo != nil {
				resultErrInfo = errInfo
				return
			} else {
				action := linebot.GetPostBackAction(
					"退出",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: domain.RED_COLOR,
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: domain.WHITE_COLOR,
							},
						),
					),
				)
			}
		} else {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.join_activity_id": activity.ActivityID,
			}
			if js, errInfo := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); errInfo != nil {
				resultErrInfo = errInfo
				return
			} else {
				action := linebot.GetPostBackAction(
					"參加",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: domain.BLUE_GREEN_COLOR,
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: domain.WHITE_COLOR,
							},
						),
					),
				)
			}
		}

		if b.currentUser.Role == domain.CADRE_CLUB_ROLE ||
			b.currentUser.Role == domain.ADMIN_CLUB_ROLE {
			cmd := domain.SUBMIT_ACTIVITY_TEXT_CMD
			pathValueMap := make(map[string]interface{})
			pathValueMap["ICmdLogic.activity_id"] = activity.ActivityID
			if js, errInfo := getCmd(cmd, pathValueMap); errInfo != nil {
				resultErrInfo = errInfo
				return
			} else {
				action := linebot.GetPostBackAction(
					"結算",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: "#1E90FF",
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: domain.WHITE_COLOR,
							},
						),
					),
				)
			}
		}
	}

	var footer *linebotModel.FlexMessageBoxComponent
	if len(footerContents) > 0 {
		footer = linebot.GetFlexMessageBoxComponent(
			linebotDomain.VERTICAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Spacing: linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
			footerContents...,
		)
	}

	return linebot.GetFlexMessageBubbleContent(
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.VERTICAL_MESSAGE_LAYOUT,
			nil,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				nil,
				contents...,
			),
		),
		&linebotModel.FlexMessagBubbleComponentOption{
			Footer: footer,
			Styles: &linebotModel.FlexMessagBubbleComponentStyle{
				Footer: &linebotModel.Background{
					Separator: true,
				},
			},
		},
	), nil
}

func (b *GetActivities) GetActivitieInfoContents(activity *NewActivity) (contents []interface{}) {
	contents = []interface{}{
		linebot.GetFlexMessageTextComponent(
			"資訊",
			&linebotModel.FlexMessageTextComponentOption{
				Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
				Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
		),
		linebot.GetSeparatorComponent(&linebotModel.FlexMessageSeparatorComponentOption{
			Margin: linebotDomain.XS_FLEX_MESSAGE_SIZE,
		}),
	}
	component := linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
		linebot.GetFlexMessageTextComponent(
			"補助額",
			&linebotModel.FlexMessageTextComponentOption{
				Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
				Color: "#555555",
			},
		),
		linebot.GetFlexMessageTextComponent(
			fmt.Sprintf("$%d", activity.ClubSubsidy),
			&linebotModel.FlexMessageTextComponentOption{
				Size:  linebotDomain.XS_FLEX_MESSAGE_SIZE,
				Color: "#111111",
				Align: linebotDomain.END_Align,
			},
		),
	)
	contents = append(contents, component)

	if activity.PeopleLimit != nil {
		peopleLimitComponent := linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin: linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageTextComponent(
				"人數上限",
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#555555",
				},
			),
			linebot.GetFlexMessageTextComponent(
				strconv.Itoa(int(*activity.PeopleLimit)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.XS_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.END_Align,
				},
			),
		)
		contents = append(contents, peopleLimitComponent)
	}
	return
}

func (b *GetActivities) GetActivitieEstimateContents(activity *getActivitiesActivity, isShowCurrentMember bool) (contents []interface{}) {
	courtHours := activity.getCourtHours()
	totalBallConsume := courtHours.MulFloat(float64(domain.ESTIMATE_BALL_CONSUME_PER_HOUR))
	estimateBallFee := totalBallConsume.MulFloat(float64(domain.PRICE_PER_BALL))
	courtFee := activity.getCourtFee()
	estimateActivityFee := estimateBallFee.Plus(courtFee)
	joinedCount := len(activity.JoinedMembers)
	peopleLimit := 0
	waitingCount := 0
	if activity.PeopleLimit != nil {
		peopleLimit = int(*activity.PeopleLimit)
		if joinedCount > peopleLimit {
			waitingCount = joinedCount - peopleLimit
			joinedCount = peopleLimit
		}
	}

	estimateBox := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			nil,
			linebot.GetFlexMessageTextComponent(
				"羽球費用",
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#555555",
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("%s顆", totalBallConsume.ToString(-1)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.CENTER_Align,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", estimateBallFee.ToString(-1)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.END_Align,
				},
			),
		),
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			nil,
			linebot.GetFlexMessageTextComponent(
				"活動費用",
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#555555",
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", estimateActivityFee.ToString(-1)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.END_Align,
				},
			),
		),
	)
	contents = []interface{}{
		linebot.GetFlexMessageTextComponent(
			"預估費用",
			&linebotModel.FlexMessageTextComponentOption{
				Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
				Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
		),
		linebot.GetSeparatorComponent(&linebotModel.FlexMessageSeparatorComponentOption{
			Margin: linebotDomain.XS_FLEX_MESSAGE_SIZE,
		}),
		estimateBox,
	}

	if isShowCurrentMember {
		estimateBox.Contents = append(
			estimateBox.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"目前參加人數",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					strconv.Itoa(joinedCount),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#111111",
						Align: linebotDomain.END_Align,
					},
				),
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"目前候補人數",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					strconv.Itoa(waitingCount),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#111111",
						Align: linebotDomain.END_Align,
					},
				),
			),
		)
	}

	if activity.PeopleLimit != nil {
		people := int(*activity.PeopleLimit)
		_, clubMemberPay, guestPay := calculateActivityPay(people, totalBallConsume, courtFee, util.NewFloat(float64(activity.ClubSubsidy)))
		estimateBox.Contents = append(
			estimateBox.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"人滿社員費用",
					&linebotModel.FlexMessageTextComponentOption{
						Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
						Color:  "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("$%d", clubMemberPay),
					&linebotModel.FlexMessageTextComponentOption{
						Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color:  "#111111",
						Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
						Align:  linebotDomain.END_Align,
					},
				),
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"人滿自費費用",
					&linebotModel.FlexMessageTextComponentOption{
						Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color:  "#555555",
						Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("$%d", guestPay),
					&linebotModel.FlexMessageTextComponentOption{
						Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color:  "#111111",
						Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
						Align:  linebotDomain.END_Align,
					},
				),
			),
		)
	}

	return
}
