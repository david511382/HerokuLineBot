package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
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

func (b *GetActivities) Init(context domain.ICmdHandlerContext) error {
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

func (b *GetActivities) LoadSingleParam(attr, text string) error {
	switch attr {
	default:
	}

	return nil
}

func (b *GetActivities) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *GetActivities) init() error {
	context := b.context
	arg := dbReqs.Activity{
		IsComplete: util.GetBoolP(false),
	}
	if dbDatas, err := database.Club.Activity.IDDatePlaceCourtsSubsidyDescriptionPeopleLimit(arg); err != nil {
		return err
	} else {
		activityIDs := []int{}
		for _, v := range dbDatas {
			activityIDs = append(activityIDs, v.ID)
		}

		activityIDMap := make(map[int][]*getActivitiesActivityJoinedMembers)
		memberActivityArg := dbReqs.MemberActivity{
			ActivityIDs: activityIDs,
		}
		if dbDatas, err := database.Club.MemberActivity.IDMemberIDActivityID(memberActivityArg); err != nil {
			return err
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
				return err
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
					Place:       v.Place,
					Description: v.Description,
					PeopleLimit: v.PeopleLimit,
					ClubSubsidy: v.ClubSubsidy,
					IsComplete:  false,
				},
				JoinedMembers: activityIDMap[v.ID],
				ActivityID:    v.ID,
			}
			if err := activity.ParseCourts(v.CourtsAndTime); err != nil {
				return err
			}
			b.activities = append(b.activities, activity)
		}
	}
	return nil
}

func (b *GetActivities) listMembers() error {
	var date time.Time
	var place string
	var peopleLimit *int16
	arg := dbReqs.Activity{
		ID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.Activity.DatePlacePeopleLimit(arg); err != nil {
		return err
	} else if len(dbDatas) == 0 {
		replyMessges := []interface{}{
			linebot.GetTextMessage("查無活動"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			return err
		}
		return nil
	} else {
		v := dbDatas[0]
		date = v.Date
		peopleLimit = v.PeopleLimit
		place = v.Place
	}

	memberComponents := []interface{}{}
	memberActivityArg := dbReqs.MemberActivity{
		ActivityID: &b.ListMembersActivityID,
	}
	if dbDatas, err := database.Club.MemberActivity.IDMemberID(memberActivityArg); err != nil {
		return err
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
			return err
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
			date.Format(commonLogicDomain.DATE_FORMAT),
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
		return err
	}

	return nil
}

func (b *GetActivities) joinActivity() error {
	userData := b.currentUser
	activityID := b.JoinActivityID
	uID := userData.ID

	insertData := &memberactivity.MemberActivityTable{
		ActivityID: activityID,
		MemberID:   uID,
		IsAttend:   false,
	}
	if err := database.Club.MemberActivity.Insert(nil, insertData); err != nil && !database.IsUniqErr(err) {
		return err
	}

	return nil
}

func (b *GetActivities) leaveActivity() error {
	userData := b.currentUser

	var peopleLimit *int16
	activityPlace := ""
	var activityDate *time.Time
	activityArg := dbReqs.Activity{
		ID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.Activity.DatePlacePeopleLimit(activityArg); err != nil {
		return err
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		peopleLimit = v.PeopleLimit
		activityPlace = v.Place
		activityDate = &v.Date
	}

	var notifyWaitingMemberID *int
	deleteMemberActivityID := 0
	arg := dbReqs.MemberActivity{
		ActivityID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.MemberActivity.IDMemberID(arg); err != nil {
		return err
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
		return err
	}

	if isNotifyWaitingPerson := notifyWaitingMemberID != nil; isNotifyWaitingPerson {
		memberArg := dbReqs.Member{
			ID: notifyWaitingMemberID,
		}
		if dbDatas, err := database.Club.Member.LineID(memberArg); err != nil {
			return err
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
										activityDate.Format(commonLogicDomain.MONTH_DATE_SLASH_FORMAT),
										commonLogic.WeekDayName(activityDate.Weekday()),
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
					return err
				}
			}
		}
	}

	return nil
}

func (b *GetActivities) loadCurrentUserID() error {
	lineID := b.context.GetUserID()
	userData, err := lineUserLogic.Get(lineID)
	if err != nil {
		return err
	} else if userData == nil {
		return domain.USER_NOT_REGISTERED
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

func (b *GetActivities) Do(text string) (resultErr error) {
	if err := b.loadCurrentUserID(); err != nil {
		return err
	}

	if isListMembers := b.ListMembersActivityID > 0; isListMembers {
		if err := b.listMembers(); err != nil {
			replyMessges := []interface{}{
				linebot.GetTextMessage("查看人員發生錯誤，已通知管理員"),
			}
			if replyErr := b.context.Reply(replyMessges); replyErr != nil {
				err = fmt.Errorf("%s---replyErr:%s", err.Error(), replyErr.Error())
			}

			return err
		}
		return nil
	} else if isJoin := b.JoinActivityID > 0; isJoin {
		activityID := b.JoinActivityID
		arg := dbReqs.Activity{
			ID:         &activityID,
			IsComplete: util.GetBoolP(false),
		}
		if count, err := database.Club.Activity.Count(arg); err != nil {
			return err
		} else if isActivityOpen := count > 0; isActivityOpen {
			if err := b.joinActivity(); err != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("參加發生錯誤，已通知管理員"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					err = fmt.Errorf("%s---replyErr:%s", err.Error(), replyErr.Error())
				}

				return err
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("活動已關閉"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				return err
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
			return err
		} else if isActivityOpen := count > 0; isActivityOpen {
			if err := b.leaveActivity(); err != nil {
				replyMessges := []interface{}{
					linebot.GetTextMessage("退出發生錯誤，已通知管理員"),
				}
				if replyErr := b.context.Reply(replyMessges); replyErr != nil {
					err = fmt.Errorf("%s---replyErr:%s", err.Error(), replyErr.Error())
				}

				return err
			}
		} else {
			replyMessges := []interface{}{
				linebot.GetTextMessage("活動已關閉"),
			}
			if err := b.context.Reply(replyMessges); err != nil {
				return err
			}
			return nil
		}
	}

	replyMessge, err := b.GetActivitiesMessage("查看活動", true, true)
	if err != nil {
		return err
	}

	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		return err
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
	err error,
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
			if js, err := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); err != nil {
				return nil, err
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
								Color: "#ffffff",
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
			if js, err := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); err != nil {
				return nil, err
			} else {
				action := linebot.GetPostBackAction(
					"退出",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: "#FF6347",
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: "#ffffff",
							},
						),
					),
				)
			}
		} else {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.join_activity_id": activity.ActivityID,
			}
			if js, err := b.context.
				GetCmdInputMode(nil).
				GetKeyValueInputMode(pathValueMap).
				GetSignal(); err != nil {
				return nil, err
			} else {
				action := linebot.GetPostBackAction(
					"參加",
					js,
				)
				footerContents = append(footerContents,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
						&linebotModel.FlexMessageBoxComponentOption{
							BackgroundColor: "#00cc99",
							CornerRadius:    "12px",
						},
						linebot.GetButtonComponent(
							action,
							&linebotModel.ButtonOption{
								Color: "#ffffff",
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
			if js, err := getCmd(cmd, pathValueMap); err != nil {
				return nil, err
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
								Color: "#ffffff",
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
	totalBallConsume := domain.ESTIMATE_BALL_CONSUME_PER_HOUR * courtHours
	estimateBallFee := totalBallConsume * domain.PRICE_PER_BALL
	courtFee := activity.getCourtFee()
	estimateActivityFee := commonLogic.FloatPlus(estimateBallFee, courtFee)
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
				fmt.Sprintf("%s顆", strconv.FormatFloat(totalBallConsume, 'f', -1, 64)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.CENTER_Align,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", strconv.FormatFloat(estimateBallFee, 'f', -1, 64)),
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
				fmt.Sprintf("$%s", strconv.FormatFloat(estimateActivityFee, 'f', -1, 64)),
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
		_, clubMemberPay, guestPay := calculateActivityPay(people, totalBallConsume, courtFee, float64(activity.ClubSubsidy))
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
