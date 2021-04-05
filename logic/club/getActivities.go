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
	"heroku-line-bot/storage/database"
	"heroku-line-bot/storage/database/database/clubdb/table/memberactivity"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"math"
	"sort"
	"strconv"
	"time"
)

type getActivities struct {
	context               domain.ICmdHandlerContext `json:"-"`
	activities            []*getActivitiesActivity
	JoinActivityID        int                       `json:"join_activity_id"`
	LeaveActivityID       int                       `json:"leave_activity_id"`
	currentUser           lineUserLogicDomain.Model `json:"-"`
	ListMembersActivityID int                       `json:"list_members_activity_id"`
}

type getActivitiesActivity struct {
	newActivity
	JoinedMembers []*getActivitiesActivityJoinedMembers `json:"joined_members"`
	ActivityID    int                                   `json:"activity_id"`
}

type getActivitiesActivityJoinedMembers struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (b *getActivities) Init(context domain.ICmdHandlerContext, initCmdBaseF func(requireRawParamAttr, requireRawParamAttrText string, isInputImmediately bool)) error {
	*b = getActivities{
		context:    context,
		activities: make([]*getActivitiesActivity, 0),
	}

	return nil
}

func (b *getActivities) GetSingleParam(attr string) string {
	switch attr {
	default:
		return ""
	}
}

func (b *getActivities) LoadSingleParam(attr, text string) error {
	switch attr {
	default:
	}

	return nil
}

func (b *getActivities) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *getActivities) init() error {
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
		if dbDatas, err := database.Club.MemberActivity.IDMemberIDActivityIDMemberName(memberActivityArg); err != nil {
			return err
		} else {
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
					Name: v.MemberName,
				})
			}
		}

		for _, v := range dbDatas {
			activity := &getActivitiesActivity{
				newActivity: newActivity{
					context:     context,
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
			if err := activity.parseCourts(v.CourtsAndTime); err != nil {
				return err
			}
			b.activities = append(b.activities, activity)
		}
	}
	return nil
}

func (b *getActivities) listMembers() error {
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
	if dbDatas, err := database.Club.MemberActivity.IDMemberIDMemberName(memberActivityArg); err != nil {
		return err
	} else {
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
					v.MemberName,
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
	joinedCount, _ := b.getJoinCount(len(memberComponents), peopleLimit)
	contents = append(contents, linebot.GetFlexMessageTextComponent(0, "參加人員:"))
	contents = append(contents, memberComponents[:joinedCount]...)
	if isUsingPeopleLimit {
		contents = append(contents, linebot.GetFlexMessageTextComponent(0, "候補人員:"))
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

func (b *getActivities) joinActivity() error {
	userData := b.currentUser

	insertData := &memberactivity.MemberActivityTable{
		ActivityID: b.JoinActivityID,
		MemberID:   userData.ID,
		MemberName: userData.Name,
	}
	if err := database.Club.MemberActivity.Insert(nil, insertData); err != nil && !database.IsUniqErr(err) {
		return err
	}

	return nil
}

func (b *getActivities) leaveActivity() error {
	userData := b.currentUser

	deleteData := &memberactivity.MemberActivityTable{}
	arg := dbReqs.MemberActivity{
		MemberID:   util.GetIntP(userData.ID),
		ActivityID: util.GetIntP(b.LeaveActivityID),
	}
	if dbDatas, err := database.Club.MemberActivity.ID(arg); err != nil {
		return err
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		deleteData.ID = v.ID
	}

	if err := database.Club.MemberActivity.Delete(nil, arg); err != nil {
		return err
	}

	return nil
}

func (b *getActivities) loadCurrentUserID() error {
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

func (b *getActivities) isJoined(activity *getActivitiesActivity) bool {
	for _, v := range activity.JoinedMembers {
		mID := v.ID
		if b.currentUser.ID == mID {
			return true
		}
	}
	return false
}

func (b *getActivities) getJoinCount(totalCount int, limit *int16) (joinedCount, waitingCount int) {
	joinedCount = totalCount
	peopleLimit := 0
	if limit != nil {
		peopleLimit = int(*limit)
		if joinedCount > peopleLimit {
			waitingCount = joinedCount - peopleLimit
			joinedCount = peopleLimit
		}
	}
	return
}

func (b *getActivities) Do(text string) (resultErr error) {
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
		if err := b.joinActivity(); err != nil {
			replyMessges := []interface{}{
				linebot.GetTextMessage("參加發生錯誤，已通知管理員"),
			}
			if replyErr := b.context.Reply(replyMessges); replyErr != nil {
				err = fmt.Errorf("%s---replyErr:%s", err.Error(), replyErr.Error())
			}

			return err
		}
	} else if isLeave := b.LeaveActivityID > 0; isLeave {
		if err := b.leaveActivity(); err != nil {
			replyMessges := []interface{}{
				linebot.GetTextMessage("退出發生錯誤，已通知管理員"),
			}
			if replyErr := b.context.Reply(replyMessges); replyErr != nil {
				err = fmt.Errorf("%s---replyErr:%s", err.Error(), replyErr.Error())
			}

			return err
		}
	}

	if err := b.init(); err != nil {
		return err
	}

	size := linebotDomain.MD_FLEX_MESSAGE_SIZE
	valueSize := linebotDomain.SM_FLEX_MESSAGE_SIZE

	carouselContents := []*linebotModel.FlexMessagBubbleComponent{}
	for _, activity := range b.activities {
		contents := []interface{}{}
		activityContents := activity.getLineComponents(domain.NewActivityLineTemplate{})
		contents = append(contents, activityContents...)

		minTime, maxTime := activity.getCourtTimeRange()
		valueText := fmt.Sprintf("%s~%s", minTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT), maxTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT))
		valueSize = linebotDomain.MD_FLEX_MESSAGE_SIZE
		timeContent := GetKeyValueEditComponent(
			"時間",
			valueText,
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
			},
		)
		contents = util.InsertAtIndex(contents, 1, timeContent)

		estimateTitleComponent := linebot.GetFlexMessageTextComponent(
			0,
			"",
			linebot.GetFlexMessageTextComponentSpan(
				"預估費用",
				linebotDomain.XL_FLEX_MESSAGE_SIZE,
				linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
			),
		)
		contents = append(contents, estimateTitleComponent)

		size = linebotDomain.MD_FLEX_MESSAGE_SIZE
		valueSize = linebotDomain.SM_FLEX_MESSAGE_SIZE
		courtHours := activity.getCourtHours()
		totalBallConsume := domain.ESTIMATE_BALL_CONSUME_PER_HOUR * courtHours
		ballComponent := GetKeyValueEditComponent(
			"預估羽球消耗",
			strconv.FormatFloat(totalBallConsume, 'f', -1, 64)+" 顆",
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, ballComponent)

		estimateBallFee := totalBallConsume * domain.PRICE_PER_BALL
		ballFeeComponent := GetKeyValueEditComponent(
			"預估羽球費用",
			strconv.FormatFloat(estimateBallFee, 'f', -1, 64),
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, ballFeeComponent)

		courtFee := activity.getCourtFee()
		courtFeeComponent := GetKeyValueEditComponent(
			"場地費用",
			strconv.FormatFloat(courtFee, 'f', -1, 64),
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, courtFeeComponent)

		estimateActivityFee := commonLogic.FloatPlus(estimateBallFee, courtFee)
		activityFeeComponent := GetKeyValueEditComponent(
			"活動費用",
			strconv.FormatFloat(estimateActivityFee, 'f', -1, 64),
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, activityFeeComponent)

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
		joinedCountComponent := GetKeyValueEditComponent(
			"目前參加人數",
			strconv.Itoa(joinedCount),
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, joinedCountComponent)
		waitingCountComponent := GetKeyValueEditComponent(
			"目前候補人數",
			strconv.Itoa(waitingCount),
			&domain.KeyValueEditComponentOption{
				ValueSizeP: &valueSize,
				SizeP:      &size,
			},
		)
		contents = append(contents, waitingCountComponent)

		if activity.PeopleLimit != nil {
			people := int(*activity.PeopleLimit)
			shareMoney := commonLogic.FloatMinus(estimateActivityFee, float64(activity.ClubSubsidy))

			p := people * domain.MONEY_UNIT
			clubMemberPay := math.Ceil(shareMoney/float64(p)) * domain.MONEY_UNIT
			clubMemberFeeComponent := GetKeyValueEditComponent(
				"預估滿人社員費用",
				strconv.FormatFloat(clubMemberPay, 'f', -1, 64),
				&domain.KeyValueEditComponentOption{
					ValueSizeP: &valueSize,
					SizeP:      &size,
				},
			)
			contents = append(contents, clubMemberFeeComponent)

			guestPay := math.Ceil(estimateActivityFee/float64(p)) * domain.MONEY_UNIT
			guestFeeComponent := GetKeyValueEditComponent(
				"預估滿人自費費用",
				strconv.FormatFloat(guestPay, 'f', -1, 64),
				&domain.KeyValueEditComponentOption{
					ValueSizeP: &valueSize,
					SizeP:      &size,
				},
			)
			contents = append(contents, guestFeeComponent)
		}

		if len(activity.JoinedMembers) > 0 {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.list_members_activity_id": activity.ActivityID,
			}
			if js, err := b.context.GetInputSignl(pathValueMap); err != nil {
				return err
			} else {
				action := linebot.GetPostBackAction(
					"查看人員",
					js,
				)
				buttonComponent := linebot.GetButtonComponent(0, action, &domain.NormalButtonOption)
				contents = append(contents, buttonComponent)
			}
		}

		if b.isJoined(activity) {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.leave_activity_id": activity.ActivityID,
			}
			if js, err := b.context.GetInputSignl(pathValueMap); err != nil {
				return err
			} else {
				action := linebot.GetPostBackAction(
					"退出",
					js,
				)
				leaveButtonComponent := linebot.GetButtonComponent(0, action, &domain.AlertButtonOption)
				contents = append(contents, leaveButtonComponent)
			}
		} else {
			pathValueMap := map[string]interface{}{
				"ICmdLogic.join_activity_id": activity.ActivityID,
			}
			if js, err := b.context.GetInputSignl(pathValueMap); err != nil {
				return err
			} else {
				action := linebot.GetPostBackAction(
					"參加",
					js,
				)
				joinButtonComponent := linebot.GetButtonComponent(0, action, &domain.NormalButtonOption)
				contents = append(contents, joinButtonComponent)
			}
		}

		carouselContents = append(
			carouselContents,
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					linebot.GetFlexMessageBoxComponent(
						linebotDomain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						contents...,
					),
				),
				nil,
			),
		)
	}

	var replyMessge interface{}
	if len(carouselContents) == 0 {
		replyMessge = linebot.GetTextMessage("沒有活動")
	} else {
		replyMessge = linebot.GetFlexMessage(
			"查看活動",
			linebot.GetFlexMessageCarouselContent(carouselContents...),
		)
	}
	replyMessges := []interface{}{
		replyMessge,
	}
	if err := b.context.Reply(replyMessges); err != nil {
		return err
	}

	return nil
}
