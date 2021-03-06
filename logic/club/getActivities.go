package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
	rdsBadmintonplaceLogic "heroku-line-bot/logic/redis/badmintonplace"
	lineUserLogic "heroku-line-bot/logic/redis/lineuser"
	lineUserLogicDomain "heroku-line-bot/logic/redis/lineuser/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	linebotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/memberactivity"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"sort"
	"strconv"
	"time"
)

type GetActivities struct {
	context               domain.ICmdHandlerContext `json:"-"`
	activities            []*getActivitiesActivity
	JoinActivityID        int                       `json:"join_activity_id"`
	LeaveActivityID       int                       `json:"leave_activity_id"`
	currentUser           lineUserLogicDomain.Model `json:"-"`
	ListMembersActivityID int                       `json:"list_members_activity_id"`
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

func (b *GetActivities) Init(context domain.ICmdHandlerContext) (resultErrInfo errLogic.IError) {
	*b = GetActivities{
		context:    context,
		activities: make([]*getActivitiesActivity, 0),
	}

	return nil
}

func (b *GetActivities) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *GetActivities) LoadSingleParam(attr, text string) (resultErrInfo errLogic.IError) {
	switch attr {
	default:
	}

	return nil
}

func (b *GetActivities) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *GetActivities) init() (resultErrInfo errLogic.IError) {
	context := b.context
	arg := dbReqs.Activity{
		IsComplete: util.GetBoolP(false),
	}
	if dbDatas, err := database.Club.Activity.IDDatePlaceIDCourtsSubsidyDescriptionPeopleLimit(arg); err != nil {
		errInfo := errLogic.NewError(err)
		if resultErrInfo == nil {
			resultErrInfo = errInfo
		} else {
			resultErrInfo = resultErrInfo.Append(errInfo)
		}
		return
	} else if len(dbDatas) > 0 {
		activityIDs := []int{}
		idPlaceMap := make(map[int]string)
		for _, v := range dbDatas {
			activityIDs = append(activityIDs, v.ID)
			idPlaceMap[v.PlaceID] = ""
		}

		placeIDs := make([]int, 0)
		for id := range idPlaceMap {
			placeIDs = append(placeIDs, id)
		}
		if dbDatas, errInfo := rdsBadmintonplaceLogic.Load(placeIDs...); errInfo != nil {
			errInfo := errLogic.NewError(err)
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
		} else {
			for id, v := range dbDatas {
				idPlaceMap[id] = v.Name
			}
		}

		activityIDMap := make(map[int][]*getActivitiesActivityJoinedMembers)
		memberActivityArg := dbReqs.MemberActivity{
			ActivityIDs: activityIDs,
		}
		if dbDatas, err := database.Club.MemberActivity.IDMemberIDActivityID(memberActivityArg); err != nil {
			errInfo := errLogic.NewError(err)
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
			memberArg := dbReqs.Member{
				IDs: memberIDs,
			}
			if dbDatas, err := database.Club.Member.IDNameLineID(memberArg); err != nil {
				resultErrInfo = errLogic.NewError(err)
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

		sort.Slice(dbDatas, func(i, j int) bool {
			return dbDatas[i].Date.Before(dbDatas[j].Date)
		})
		for _, v := range dbDatas {
			activity := &getActivitiesActivity{
				NewActivity: NewActivity{
					Context:     context,
					Date:        v.Date,
					PlaceID:     v.PlaceID,
					Description: v.Description,
					PeopleLimit: v.PeopleLimit,
					ClubSubsidy: v.ClubSubsidy,
					IsComplete:  false,
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
	}

	return
}

func (b *GetActivities) listMembers() (resultErrInfo errLogic.IError) {
	var date time.Time
	var place string
	var peopleLimit *int16
	arg := dbReqs.Activity{
		ID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.Activity.DatePlaceIDPeopleLimit(arg); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		replyMessges := []interface{}{
			linebot.GetTextMessage("????????????"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		return nil
	} else {
		v := dbDatas[0]
		date = v.Date
		peopleLimit = v.PeopleLimit

		if dbDatas, errInfo := rdsBadmintonplaceLogic.Load(v.PlaceID); errInfo != nil {
			errInfo := errLogic.NewError(err)
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
	memberActivityArg := dbReqs.MemberActivity{
		ActivityID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.MemberActivity.IDMemberID(memberActivityArg); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else {
		memberIDs := make([]int, 0)
		for _, v := range dbDatas {
			memberIDs = append(memberIDs, v.MemberID)
		}
		memberIDNameMap := make(map[int]string)
		memberArg := dbReqs.Member{
			IDs: memberIDs,
		}
		if dbDatas, err := database.Club.Member.IDName(memberArg); err != nil {
			resultErrInfo = errLogic.NewError(err)
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
			"??????",
			date.Format(commonLogicDomain.DATE_FORMAT),
			nil,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"??????",
			place,
			nil,
		),
	)

	isUsingPeopleLimit := peopleLimit != nil
	joinedCount, _ := getJoinCount(len(memberComponents), peopleLimit)
	contents = append(contents, linebot.GetFlexMessageTextComponent("????????????:", nil))
	contents = append(contents, memberComponents[:joinedCount]...)
	if isUsingPeopleLimit {
		contents = append(contents, linebot.GetFlexMessageTextComponent("????????????:", nil))
		contents = append(contents, memberComponents[joinedCount:]...)
	}

	replyMessge := linebot.GetFlexMessage(
		"????????????",
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
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return nil
}

func (b *GetActivities) joinActivity() (resultErrInfo errLogic.IError) {
	userData := b.currentUser
	activityID := b.JoinActivityID
	uID := userData.ID

	insertData := &memberactivity.MemberActivityTable{
		ActivityID: activityID,
		MemberID:   uID,
		IsAttend:   false,
	}
	if err := database.Club.MemberActivity.Insert(nil, insertData); err != nil && !database.IsUniqErr(err) {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return nil
}

func (b *GetActivities) leaveActivity() (resultErrInfo errLogic.IError) {
	userData := b.currentUser

	var peopleLimit *int16
	activityPlace := ""
	var activityDate *time.Time
	activityArg := dbReqs.Activity{
		ID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.Activity.DatePlaceIDPeopleLimit(activityArg); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		peopleLimit = v.PeopleLimit
		activityDate = &v.Date

		if dbDatas, errInfo := rdsBadmintonplaceLogic.Load(v.PlaceID); errInfo != nil {
			errInfo := errLogic.NewError(err)
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
	arg := dbReqs.MemberActivity{
		ActivityID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.MemberActivity.IDMemberID(arg); err != nil {
		resultErrInfo = errLogic.NewError(err)
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

	arg = dbReqs.MemberActivity{
		ID: &deleteMemberActivityID,
	}
	if err := database.Club.MemberActivity.Delete(nil, arg); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	if isNotifyWaitingPerson := notifyWaitingMemberID != nil; isNotifyWaitingPerson {
		memberArg := dbReqs.Member{
			ID: notifyWaitingMemberID,
		}
		if dbDatas, err := database.Club.Member.LineID(memberArg); err != nil {
			resultErrInfo = errLogic.NewError(err)
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
						"??????????????????",
						linebot.GetFlexMessageBubbleContent(
							linebot.GetFlexMessageBoxComponent(
								linebotDomain.VERTICAL_MESSAGE_LAYOUT,
								nil,
								linebot.GetTextMessage("????????????????????????!!"),
								linebot.GetTextMessage(
									fmt.Sprintf("?????? %s(%s) %s",
										activityDate.Format(commonLogicDomain.MONTH_DATE_SLASH_FORMAT),
										commonLogic.WeekDayName(activityDate.Weekday()),
										activityPlace,
									),
								),
								linebot.GetTextMessage("???????????????????????????????????????~"),
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
					resultErrInfo = errLogic.NewError(err)
					return
				}
			}
		}
	}

	return nil
}

func (b *GetActivities) loadCurrentUserID() (resultErrInfo errLogic.IError) {
	lineID := b.context.GetUserID()
	userData, err := lineUserLogic.Get(lineID)
	if err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else if userData == nil {
		resultErrInfo = errLogic.NewError(domain.USER_NOT_REGISTERED)
		return
	}

	b.currentUser = *userData

	return nil
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

func (b *GetActivities) Do(text string) (resultErrInfo errLogic.IError) {
	if errInfo := b.loadCurrentUserID(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	if isListMembers := b.ListMembersActivityID > 0; isListMembers {
		if errInfo := b.listMembers(); errInfo != nil {
			replyMessges := []interface{}{
				linebot.GetTextMessage("?????????????????????????????????????????????"),
			}
			if replyErr := b.context.Reply(replyMessges); replyErr != nil {
				resultErrInfo = errLogic.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
				return
			}

			resultErrInfo = errInfo
			return
		}
		return nil
	} else if isJoin := b.JoinActivityID > 0; isJoin {
		activityID := b.JoinActivityID
		arg := dbReqs.Activity{
			ID:         &activityID,
			IsComplete: util.GetBoolP(false),
		}
		if count, err := database.Club.Activity.Count(arg); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		} else if isActivityOpen := count > 0; isActivityOpen {
			if errInfo := b.joinActivity(); errInfo != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("???????????????????????????????????????"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					resultErrInfo = errLogic.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
					return
				}

				resultErrInfo = errInfo
				return
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("???????????????"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			}
			return nil
		}
	} else if isLeave := b.LeaveActivityID > 0; isLeave {
		activityID := b.LeaveActivityID
		arg := dbReqs.Activity{
			ID:         &activityID,
			IsComplete: util.GetBoolP(false),
		}
		if count, err := database.Club.Activity.Count(arg); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		} else if isActivityOpen := count > 0; isActivityOpen {
			if errInfo := b.leaveActivity(); errInfo != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("???????????????????????????????????????"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					resultErrInfo = errLogic.Newf("%s---replyErr:%s", errInfo.Error(), replyErr.Error())
					return
				}

				resultErrInfo = errInfo
				return
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("???????????????"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			}
			return nil
		}
	}

	replyMessge, err := b.GetActivitiesMessage("????????????", true, true)
	if err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errLogic.NewError(err)
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
		replyMessge = linebot.GetTextMessage("????????????")
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
	resultErrInfo errLogic.IError,
) {
	contents := []interface{}{
		linebot.GetFlexMessageTextComponent(
			"????????????",
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
					"????????????",
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
					"??????",
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
					"??????",
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
					"??????",
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
			"??????",
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
			"?????????",
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
				"????????????",
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
				"????????????",
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#555555",
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("%s???", totalBallConsume.ToString(-1)),
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
				"????????????",
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
			"????????????",
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
					"??????????????????",
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
					"??????????????????",
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
		_, clubMemberPay, guestPay := calculateActivityPay(people, totalBallConsume, courtFee, util.ToFloat(float64(activity.ClubSubsidy)))
		estimateBox.Contents = append(
			estimateBox.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"??????????????????",
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
					"??????????????????",
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
