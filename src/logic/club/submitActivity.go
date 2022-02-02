package club

import (
	"fmt"
	"heroku-line-bot/src/logic/club/domain"
	clubLineuserLogic "heroku-line-bot/src/logic/club/lineuser"
	clubLineuserLogicDomain "heroku-line-bot/src/logic/club/lineuser/domain"
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/member"
	"heroku-line-bot/src/repo/database/database/clubdb/memberactivity"
	"heroku-line-bot/src/service/linebot"
	linebotDomain "heroku-line-bot/src/service/linebot/domain"
	linebotModel "heroku-line-bot/src/service/linebot/domain/model"
	"heroku-line-bot/src/util"
	errUtil "heroku-line-bot/src/util/error"
	"sort"
	"strconv"
)

type submitActivity struct {
	context domain.ICmdHandlerContext `json:"-"`
	NewActivity
	JoinedMembers  []*submitActivityJoinedMembers `json:"joined_members"`
	JoinedGuests   []*submitActivityJoinedMembers `json:"joined_guests"`
	ActivityID     int                            `json:"activity_id"`
	CurrentUser    *clubLineuserLogicDomain.Model `json:"current_user"`
	HasLoad        bool                           `json:"has_load"`
	Rsl4Consume    int16                          `json:"rsl4_consume"`
	AttendIndex    *int                           `json:"attend_index,omitempty"`
	PayIndex       *int                           `json:"pay_index,omitempty"`
	IsJoinedMember bool                           `json:"is_joined_member_index"`
}

type submitActivityJoinedMembers struct {
	getActivitiesActivityJoinedMembers
	IsAttend         bool `json:"is_attend"`
	IsPaid           bool `json:"is_paid"`
	MemberActivityID int  `json:"id"`
}

func (b *submitActivity) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	*b = submitActivity{
		context: context,
	}

	return nil
}

func (b *submitActivity) GetSingleParam(attr string) string {
	switch attr {
	case "rsl4_consume":
		return strconv.Itoa(int(b.Rsl4Consume))
	default:
		return ""
	}
}

func (b *submitActivity) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "rsl4_consume":
		i, err := strconv.Atoi(text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.Rsl4Consume = int16(i)
	default:
	}

	return nil
}

func (b *submitActivity) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *submitActivity) init() (resultErrInfo errUtil.IError) {
	if b.HasLoad {
		return nil
	}

	context := b.context
	arg := dbModel.ReqsClubActivity{
		ID: util.GetIntP(b.ActivityID),
	}
	if dbDatas, err := database.Club.Activity.Select(
		arg,
		activity.COLUMN_ID,
		activity.COLUMN_Date,
		activity.COLUMN_PlaceID,
		activity.COLUMN_CourtsAndTime,
		activity.COLUMN_ClubSubsidy,
		activity.COLUMN_Description,
		activity.COLUMN_PeopleLimit,
		activity.COLUMN_TeamID,
	); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		memberJoinDate := v.Date
		b.NewActivity = NewActivity{
			Context:     context,
			Date:        util.DateTime(v.Date),
			PlaceID:     v.PlaceID,
			Description: v.Description,
			PeopleLimit: v.PeopleLimit,
			ClubSubsidy: v.ClubSubsidy,
			TeamID:      v.TeamID,
		}
		if errInfo := b.NewActivity.ParseCourts(v.CourtsAndTime); errInfo != nil {
			resultErrInfo = errInfo
			return
		}

		memberActivityArg := dbModel.ReqsClubMemberActivity{
			ActivityID: util.GetIntP(b.ActivityID),
		}
		if dbDatas, err := database.Club.MemberActivity.Select(
			memberActivityArg,
			memberactivity.COLUMN_ID,
			memberactivity.COLUMN_MemberID,
		); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			type isClubMemberName struct {
				isClubMember bool
				name         string
			}
			memberIDs := []int{}
			for _, v := range dbDatas {
				memberIDs = append(memberIDs, v.MemberID)
			}
			arg := dbModel.ReqsClubMember{
				IDs: memberIDs,
			}
			clubMemberIDMap := make(map[int]isClubMemberName)
			if dbDatas, err := database.Club.Member.Select(
				arg,
				member.COLUMN_ID,
				member.COLUMN_Name,
				member.COLUMN_Department,
				member.COLUMN_JoinDate,
			); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			} else {
				for _, v := range dbDatas {
					isClubMember := false
					if v.JoinDate != nil && !v.JoinDate.After(memberJoinDate) {
						isClubMember = Department(v.Department).IsClubMember()
					}

					clubMemberIDMap[v.ID] = isClubMemberName{
						name:         v.Name,
						isClubMember: isClubMember,
					}
				}
			}

			sort.Slice(dbDatas, func(i, j int) bool {
				return dbDatas[i].ID < dbDatas[j].ID
			})
			peopleLimit, _ := getJoinCount(len(dbDatas), b.PeopleLimit)
			dbDatas = dbDatas[:peopleLimit]
			for _, v := range dbDatas {
				memberID := v.MemberID
				clubMember := clubMemberIDMap[memberID]
				member := &submitActivityJoinedMembers{
					getActivitiesActivityJoinedMembers: getActivitiesActivityJoinedMembers{
						ID:   v.MemberID,
						Name: clubMember.name,
					},
					MemberActivityID: v.ID,
				}
				if clubMember.isClubMember {
					b.JoinedMembers = append(b.JoinedMembers, member)
				} else {
					b.JoinedGuests = append(b.JoinedGuests, member)
				}
			}
		}
	}

	b.HasLoad = true

	return nil
}

func (b *submitActivity) getJoinedMembersCount() int {
	people := 0
	for _, member := range b.JoinedMembers {
		if member.IsAttend {
			people++
		}
	}
	return people
}

func (b *submitActivity) getJoinedGuestsCount() int {
	people := 0
	for _, member := range b.JoinedGuests {
		if member.IsAttend {
			people++
		}
	}
	return people
}

func (b *submitActivity) loadCurrentUserID() (resultErrInfo errUtil.IError) {
	if b.CurrentUser != nil {
		return nil
	}

	lineID := b.context.GetUserID()
	userData, err := clubLineuserLogic.Get(lineID)
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if userData == nil {
		resultErrInfo = errUtil.NewError(domain.USER_NOT_REGISTERED)
		return
	}

	b.CurrentUser = userData

	return nil
}

func (b *submitActivity) Do(text string) (resultErrInfo errUtil.IError) {
	if errInfo := b.loadCurrentUserID(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	if b.CurrentUser.Role != domain.CADRE_CLUB_ROLE &&
		b.CurrentUser.Role != domain.ADMIN_CLUB_ROLE {
		resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if errInfo := b.init(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	if !b.HasLoad {
		replyMessges := []interface{}{
			linebot.GetTextMessage("活動不存在"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
	}

	if b.context.IsComfirmed() {
		var currentActivity *dbModel.ClubActivity
		{
			dbDatas, err := database.Club.Activity.Select(dbModel.ReqsClubActivity{
				ID: &b.ActivityID,
			})
			if err != nil {
				errInfo := errUtil.NewError(err)
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			} else if len(dbDatas) == 0 {
				replyMessges := []interface{}{
					linebot.GetTextMessage("活動不存在"),
				}
				if err := b.context.Reply(replyMessges); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				}

				if err := b.context.DeleteParam(); err != nil {
					resultErrInfo = errUtil.NewError(err)
					return
				}

				return
			}

			currentActivity = dbDatas[0]
		}

		memberActivityIDs := make([]int, 0)
		finishedActivity := &dbModel.ClubActivityFinished{
			ID:            currentActivity.ID,
			TeamID:        currentActivity.TeamID,
			Date:          currentActivity.Date,
			PlaceID:       currentActivity.PlaceID,
			CourtsAndTime: currentActivity.CourtsAndTime,
			ClubSubsidy:   currentActivity.ClubSubsidy,
			Description:   currentActivity.Description,
			PeopleLimit:   currentActivity.PeopleLimit,
		}
		{
			memberCount := 0
			for _, member := range b.JoinedMembers {
				if member.IsAttend {
					memberCount++
					memberActivityIDs = append(memberActivityIDs, member.MemberActivityID)
				}
			}
			guestCount := 0
			for _, member := range b.JoinedGuests {
				if member.IsAttend {
					guestCount++
					memberActivityIDs = append(memberActivityIDs, member.MemberActivityID)
				}
			}
			peopleCount := memberCount + guestCount
			courtFee := b.getCourtFee()
			_, memberFee, guestFee := calculateActivityPay(
				peopleCount,
				util.NewFloat(float64(b.Rsl4Consume)),
				courtFee,
				util.NewFloat(float64(b.ClubSubsidy)),
			)

			finishedActivity.MemberCount = int16(memberCount)
			finishedActivity.GuestCount = int16(guestCount)
			finishedActivity.MemberFee = int16(memberFee)
			finishedActivity.GuestFee = int16(guestFee)
		}

		db, transaction, err := database.Club.Begin()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		defer func() {
			if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			}
		}()

		isConsumeBall := b.Rsl4Consume > 0
		logisticID := 0
		if isConsumeBall {
			logisticData := &dbModel.ClubLogistic{
				Date:        b.Date.Time(),
				Name:        domain.BALL_NAME,
				Amount:      -b.Rsl4Consume,
				Description: "打球",
				TeamID:      b.TeamID,
			}
			if err := db.Logistic.Insert(logisticData); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			}
			logisticID = logisticData.ID
		}
		if isConsumeBall {
			finishedActivity.LogisticID = &logisticID
		}

		if err := db.Activity.Delete(dbModel.ReqsClubActivity{
			ID: &currentActivity.ID,
		}); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		if err := db.ActivityFinished.Insert(finishedActivity); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		if len(memberActivityIDs) > 0 {
			arg := dbModel.ReqsClubMemberActivity{
				IDs: memberActivityIDs,
			}
			fields := map[string]interface{}{
				"is_attend": true,
			}
			if err := db.MemberActivity.Update(arg, fields); err != nil && !database.IsUniqErr(err) {
				resultErrInfo = errUtil.NewError(err)
				return
			}
		}

		if err := b.context.DeleteParam(); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		return nil
	}

	if b.AttendIndex != nil {
		attendIndex := *b.AttendIndex
		if b.IsJoinedMember {
			b.JoinedMembers[attendIndex].IsAttend = !b.JoinedMembers[attendIndex].IsAttend
		} else {
			b.JoinedGuests[attendIndex].IsAttend = !b.JoinedGuests[attendIndex].IsAttend
		}
	} else if b.PayIndex != nil {
		payIndex := *b.PayIndex
		if b.IsJoinedMember {
			b.JoinedMembers[payIndex].IsPaid = !b.JoinedMembers[payIndex].IsPaid
		} else {
			b.JoinedGuests[payIndex].IsPaid = !b.JoinedGuests[payIndex].IsPaid
		}
	}
	b.AttendIndex = nil
	b.PayIndex = nil

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	bodyBox := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		nil,
		linebot.GetFlexMessageTextComponent(
			"結算資訊",
			&linebotModel.FlexMessageTextComponentOption{
				Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
				Color:  "#1DB446",
			},
		),
	)

	bodyBox.Contents = append(bodyBox.Contents, b.getPlaceTimeTemplate()...)

	boxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
	)
	boxComponent.Contents = append(boxComponent.Contents, b.getCourtsContents()...)
	boxComponent.Contents = append(boxComponent.Contents, b.getAttendInfoContents()...)
	boxComponent.Contents = append(boxComponent.Contents, b.getFeeContents()...)

	b.IsJoinedMember = true
	attendComponents, err := b.getAttendComponent("社員", b.JoinedMembers)
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}
	boxComponent.Contents = append(boxComponent.Contents, attendComponents...)

	b.IsJoinedMember = false
	attendComponents, err = b.getAttendComponent("自費人員", b.JoinedGuests)
	if err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}
	boxComponent.Contents = append(boxComponent.Contents, attendComponents...)
	bodyBox.Contents = append(bodyBox.Contents, boxComponent)

	footerBox := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			JustifyContent: linebotDomain.FLEX_END_JUSTIFY_CONTENT,
			Spacing:        linebotDomain.MD_FLEX_MESSAGE_SIZE,
		},
	)

	if js, errInfo := b.context.
		GetRequireInputMode("rsl4_consume", "使用羽球數", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		footerBox.Contents = append(footerBox.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					BackgroundColor: "#FFBF00",
					CornerRadius:    "12px",
				},
				linebot.GetButtonComponent(
					linebot.GetPostBackAction(
						"使用羽球數",
						js,
					),
					&linebotModel.ButtonOption{
						Color:  domain.WHITE_COLOR,
						Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
					},
				),
			),
		)
	}

	if js, errInfo := b.context.
		GetComfirmMode().
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		footerBox.Contents = append(footerBox.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					BackgroundColor: "#1E90FF",
					CornerRadius:    "12px",
				},
				linebot.GetButtonComponent(
					linebot.GetPostBackAction(
						"提交",
						js,
					),
					&linebotModel.ButtonOption{
						Color:  domain.WHITE_COLOR,
						Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
					},
				),
			),
		)
	}

	replyMessage := linebot.GetFlexMessage(
		"查看活動",
		linebot.GetFlexMessageBubbleContent(
			bodyBox,
			&linebotModel.FlexMessagBubbleComponentOption{
				Footer: footerBox,
			},
		),
	)
	replyMessages := []interface{}{
		replyMessage,
	}
	if err := b.context.Reply(replyMessages); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *submitActivity) getAttendComponent(text string, members []*submitActivityJoinedMembers) ([]interface{}, error) {
	memberCount := 0
	memberBoxs := make([]interface{}, 0)
	for id, member := range members {
		var attendAction *linebotModel.PostBackAction
		attendButtonColor := domain.RED_COLOR
		if member.IsAttend {
			memberCount++
			attendButtonColor = domain.BLUE_GREEN_COLOR
		}
		pathValueMap := map[string]interface{}{
			"ICmdLogic.attend_index":           id,
			"ICmdLogic.is_joined_member_index": b.IsJoinedMember,
		}
		if js, errInfo := b.context.
			GetKeyValueInputMode(pathValueMap).
			GetSignal(); errInfo != nil {
			return nil, errInfo
		} else {
			attendAction = linebot.GetPostBackAction(
				"簽到",
				js,
			)
		}

		var payAction *linebotModel.PostBackAction
		payButtonColor := domain.RED_COLOR
		if member.IsPaid {
			payButtonColor = domain.BLUE_GREEN_COLOR
		}
		pathValueMap = map[string]interface{}{
			"ICmdLogic.is_joined_member_index": b.IsJoinedMember,
			"ICmdLogic.pay_index":              id,
		}
		if js, errInfo := b.context.
			GetKeyValueInputMode(pathValueMap).
			GetSignal(); errInfo != nil {
			return nil, errInfo
		} else {
			payAction = linebot.GetPostBackAction(
				"收費",
				js,
			)
		}

		contents := []interface{}{
			linebot.GetFlexMessageTextComponent(
				member.Name,
				&linebotModel.FlexMessageTextComponentOption{
					Size: linebotDomain.XS_FLEX_MESSAGE_SIZE,
				},
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					JustifyContent: linebotDomain.FLEX_END_JUSTIFY_CONTENT,
					Spacing:        linebotDomain.XS_FLEX_MESSAGE_SIZE,
				},
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: attendButtonColor,
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						attendAction,
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: payButtonColor,
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						payAction,
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
			),
		}

		memberBoxs = append(memberBoxs,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					Margin:     linebotDomain.MD_FLEX_MESSAGE_SIZE,
					AlignItems: linebotDomain.CENTER_ALIGN_ITEMS,
				},
				contents...,
			),
		)
	}

	result := []interface{}{
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageTextComponent(
				text,
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("%d人", memberCount),
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
					Align:  linebotDomain.END_Align,
				},
			),
		),
		linebot.GetSeparatorComponent(&linebotModel.FlexMessageSeparatorComponentOption{
			Margin: linebotDomain.XS_FLEX_MESSAGE_SIZE,
		}),
	}

	if len(memberBoxs) > 0 {
		result = append(result, linebot.GetFlexMessageBoxComponent(
			linebotDomain.VERTICAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
				Spacing: linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
			memberBoxs...,
		))
	}

	return result, nil
}

func (b *submitActivity) getAttendInfoContents() []interface{} {
	clubMemberPeople := b.getJoinedMembersCount()
	guestPeople := b.getJoinedGuestsCount()
	people := clubMemberPeople + guestPeople

	return []interface{}{
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageTextComponent(
				"參加人數",
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
					Color:  "#555555",
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("%d人", people),
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
					Align:  linebotDomain.END_Align,
				},
			),
		),
		linebot.GetSeparatorComponent(&linebotModel.FlexMessageSeparatorComponentOption{
			Margin: linebotDomain.XS_FLEX_MESSAGE_SIZE,
		}),
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.VERTICAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
				Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"社員",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("%d人", clubMemberPeople),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.END_Align,
						Color: "#111111",
					},
				),
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"自費人員",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("%d人", guestPeople),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.END_Align,
						Color: "#111111",
					},
				),
			),
		),
	}
}

func (b *submitActivity) getFeeContents() []interface{} {
	courtFee := b.getCourtFee()
	activityFee, ballFee := calculateActivity(
		util.NewFloat(float64(b.Rsl4Consume)),
		courtFee,
	)
	clubMemberPeople := b.getJoinedMembersCount()
	guestPeople := b.getJoinedGuestsCount()
	people := clubMemberPeople + guestPeople
	clubSubsidy := util.NewFloat(float64(b.ClubSubsidy))
	shareMoney := activityFee.Minus(clubSubsidy)

	result := []interface{}{
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageTextComponent(
				"費用",
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", shareMoney.ToString(0)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
					Align:  linebotDomain.END_Align,
				},
			),
		),
		linebot.GetSeparatorComponent(&linebotModel.FlexMessageSeparatorComponentOption{
			Margin: linebotDomain.XS_FLEX_MESSAGE_SIZE,
		}),
		linebot.GetFlexMessageBoxComponent(
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
					fmt.Sprintf("%d顆", int(b.Rsl4Consume)),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.CENTER_Align,
						Color: "#111111",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("$%s", ballFee.ToString(-1)),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.END_Align,
						Color: "#111111",
					},
				),
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"場地費用",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("$%s", courtFee.ToString(-1)),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.END_Align,
						Color: "#111111",
					},
				),
			),
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					"補助額",
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageTextComponent(
					fmt.Sprintf("$-%d", b.ClubSubsidy),
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
						Align: linebotDomain.END_Align,
						Color: "#111111",
					},
				),
			),
		),
	}

	if people > 0 {
		clubMemberPay, guestPay := calculatePay(people, activityFee, clubSubsidy)
		result = append(
			result,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
					Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
				},
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					nil,
					linebot.GetFlexMessageTextComponent(
						"社員費用",
						&linebotModel.FlexMessageTextComponentOption{
							Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
							Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
							Color:  "#555555",
						},
					),
					linebot.GetFlexMessageTextComponent(
						fmt.Sprintf("$%d/人", clubMemberPay),
						&linebotModel.FlexMessageTextComponentOption{
							Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
							Align: linebotDomain.END_Align,
							Color: "#111111",
						},
					),
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					nil,
					linebot.GetFlexMessageTextComponent(
						"自費費用",
						&linebotModel.FlexMessageTextComponentOption{
							Size:   linebotDomain.SM_FLEX_MESSAGE_SIZE,
							Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
							Color:  "#555555",
						},
					),
					linebot.GetFlexMessageTextComponent(
						fmt.Sprintf("$%d/人", guestPay),
						&linebotModel.FlexMessageTextComponentOption{
							Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
							Align: linebotDomain.END_Align,
							Color: "#111111",
						},
					),
				),
			),
		)
	}

	return result
}