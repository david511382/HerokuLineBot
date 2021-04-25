package club

import (
	"heroku-line-bot/logic/club/domain"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	lineUserLogic "heroku-line-bot/logic/redis/lineuser"
	lineUserLogicDomain "heroku-line-bot/logic/redis/lineuser/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/storage/database"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	"sort"
	"strconv"
)

type submitActivity struct {
	context domain.ICmdHandlerContext `json:"-"`
	NewActivity
	JoinedMembers  []*submitActivityJoinedMembers `json:"joined_members"`
	JoinedGuests   []*submitActivityJoinedMembers `json:"joined_guests"`
	ActivityID     int                            `json:"activity_id"`
	CurrentUser    *lineUserLogicDomain.Model     `json:"current_user"`
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

func (b *submitActivity) Init(context domain.ICmdHandlerContext) error {
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

func (b *submitActivity) LoadSingleParam(attr, text string) error {
	switch attr {
	case "rsl4_consume":
		i, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		b.Rsl4Consume = int16(i)
	default:
	}

	return nil
}

func (b *submitActivity) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *submitActivity) init() error {
	if b.HasLoad {
		return nil
	}

	context := b.context
	arg := dbReqs.Activity{
		ID: util.GetIntP(b.ActivityID),
	}
	if dbDatas, err := database.Club.Activity.IDDatePlaceCourtsSubsidyDescriptionPeopleLimit(arg); err != nil {
		return err
	} else if len(dbDatas) == 0 {
		return nil
	} else {
		v := dbDatas[0]
		b.NewActivity = NewActivity{
			Context:     context,
			Date:        v.Date,
			Place:       v.Place,
			Description: v.Description,
			PeopleLimit: v.PeopleLimit,
			ClubSubsidy: v.ClubSubsidy,
			IsComplete:  false,
		}
		if err := b.NewActivity.ParseCourts(v.CourtsAndTime); err != nil {
			return err
		}

		memberActivityArg := dbReqs.MemberActivity{
			ActivityID: util.GetIntP(b.ActivityID),
		}
		if dbDatas, err := database.Club.MemberActivity.IDMemberIDMemberName(memberActivityArg); err != nil {
			return err
		} else {
			memberIDs := []int{}
			for _, v := range dbDatas {
				memberIDs = append(memberIDs, v.MemberID)
			}
			arg := dbReqs.Member{
				IDs: memberIDs,
			}
			memberIDMap := make(map[int]bool)
			if dbDatas, err := database.Club.Member.IDDepartment(arg); err != nil {
				return err
			} else {
				for _, v := range dbDatas {
					memberIDMap[v.ID] = Department(v.Department).IsClubMember()
				}
			}

			sort.Slice(dbDatas, func(i, j int) bool {
				return dbDatas[i].ID < dbDatas[j].ID
			})
			peopleLimit, _ := getJoinCount(len(dbDatas), b.PeopleLimit)
			dbDatas = dbDatas[:peopleLimit]
			for _, v := range dbDatas {
				memberID := v.MemberID
				member := &submitActivityJoinedMembers{
					getActivitiesActivityJoinedMembers: getActivitiesActivityJoinedMembers{
						ID:   v.MemberID,
						Name: v.MemberName,
					},
					MemberActivityID: v.ID,
				}
				if isClubMember := memberIDMap[memberID]; isClubMember {
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

func (b *submitActivity) loadCurrentUserID() error {
	if b.CurrentUser != nil {
		return nil
	}

	lineID := b.context.GetUserID()
	userData, err := lineUserLogic.Get(lineID)
	if err != nil {
		return err
	} else if userData == nil {
		return domain.USER_NOT_REGISTERED
	}

	b.CurrentUser = userData

	return nil
}

func (b *submitActivity) Do(text string) (resultErr error) {
	if err := b.loadCurrentUserID(); err != nil {
		return err
	}

	if b.CurrentUser.Role != domain.CADRE_CLUB_ROLE &&
		b.CurrentUser.Role != domain.ADMIN_CLUB_ROLE {
		return domain.NO_AUTH_ERROR
	}

	if err := b.init(); err != nil {
		return err
	}

	if !b.HasLoad {
		replyMessges := []interface{}{
			linebot.GetTextMessage("活動不存在"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			return err
		}
	}

	if b.context.IsComfirmed() {
		transaction := database.Club.Begin()
		if err := transaction.Error; err != nil {
			return err
		}
		defer func() {
			if resultErr == nil {
				if resultErr = transaction.Commit().Error; resultErr != nil {
					return
				}
			}

			if err := transaction.Rollback().Error; err != nil {
				if resultErr == nil {
					resultErr = err
				}
				return
			}
		}()
		arg := dbReqs.Activity{
			ID: &b.ActivityID,
		}
		courtFee := b.getCourtFee()
		memberActivityIDs := make([]int, 0)
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
		_, memberFee, guestFee := calculateActivityPay(peopleCount, float64(b.Rsl4Consume), courtFee, float64(b.ClubSubsidy))
		updateFields := map[string]interface{}{
			"is_complete":  true,
			"rsl4_count":   b.Rsl4Consume,
			"member_count": memberCount,
			"guest_count":  guestCount,
			"member_fee":   memberFee,
			"guest_fee":    guestFee,
		}
		if resultErr = database.Club.Activity.Update(transaction, arg, updateFields); resultErr != nil {
			return
		}

		if len(memberActivityIDs) > 0 {
			arg := dbReqs.MemberActivity{
				IDs: memberActivityIDs,
			}
			fields := map[string]interface{}{
				"is_attend": true,
			}
			if err := database.Club.MemberActivity.Update(transaction, arg, fields); err != nil && !database.IsUniqErr(err) {
				return err
			}
		}

		if resultErr = b.context.DeleteParam(); resultErr != nil {
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if resultErr = b.context.Reply(replyMessges); resultErr != nil {
			return resultErr
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

	if err := b.context.CacheParams(); err != nil {
		return err
	}

	mdSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	smallKeyValueEditComponentOption := &domain.KeyValueEditComponentOption{
		SizeP: &mdSize,
	}
	contents := []interface{}{}

	contents = append(contents,
		GetKeyValueEditComponent(
			"日期",
			b.Date.Format(commonLogicDomain.DATE_FORMAT),
			nil,
		),
	)

	contents = append(contents,
		GetKeyValueEditComponent(
			"地點",
			b.Place,
			nil,
		),
	)

	contents = append(contents, b.getCourtsBoxComponent(nil))

	if js, err := b.context.
		GetRequireInputMode("rsl4_consume", "使用羽球數", false).
		GetSignal(); err != nil {
		return err
	} else {
		action := linebot.GetPostBackAction(
			"使用羽球數",
			js,
		)
		contents = append(contents,
			GetKeyValueEditComponent(
				"使用羽球數",
				strconv.Itoa(int(b.Rsl4Consume)),
				&domain.KeyValueEditComponentOption{
					Action: action,
				},
			),
		)
	}

	courtFee := b.getCourtFee()

	activityFee, ballFee := calculateActivity(float64(b.Rsl4Consume), courtFee)
	ballFeeComponent := GetKeyValueEditComponent(
		"羽球費用",
		strconv.FormatFloat(ballFee, 'f', -1, 64),
		nil,
	)
	contents = append(contents, ballFeeComponent)
	activityFeeComponent := GetKeyValueEditComponent(
		"活動費用",
		strconv.FormatFloat(activityFee, 'f', -1, 64),
		nil,
	)
	contents = append(contents, activityFeeComponent)

	contents = append(contents,
		GetKeyValueEditComponent(
			"補助額",
			strconv.Itoa(int(b.ClubSubsidy)),
			nil,
		),
	)

	clubMemberPeople := b.getJoinedMembersCount()
	guestPeople := b.getJoinedGuestsCount()
	people := clubMemberPeople + guestPeople
	contents = append(contents,
		GetKeyValueEditComponent(
			"參加人數",
			strconv.Itoa(people),
			nil,
		),
	)

	component := GetDoubleKeyValueComponent(
		"社員人數",
		strconv.Itoa(clubMemberPeople),
		"自費人數",
		strconv.Itoa(guestPeople),
		nil,
		smallKeyValueEditComponentOption,
	)
	contents = append(contents, component)

	if people > 0 {
		clubMemberPay, guestPay := calculatePay(people, activityFee, float64(b.ClubSubsidy))
		clubMemberFeeComponent := GetKeyValueEditComponent(
			"社員費用",
			strconv.Itoa(clubMemberPay),
			nil,
		)
		contents = append(contents, clubMemberFeeComponent)
		guestFeeComponent := GetKeyValueEditComponent(
			"自費費用",
			strconv.Itoa(guestPay),
			nil,
		)
		contents = append(contents, guestFeeComponent)
	}

	contents = append(contents, linebot.GetFlexMessageTextComponent(0, "社員:"))
	b.IsJoinedMember = true
	attendComponents, err := b.getAttendComponent(b.JoinedMembers)
	if err != nil {
		return err
	}
	contents = append(contents, attendComponents...)
	contents = append(contents, linebot.GetFlexMessageTextComponent(0, "自費人員:"))
	b.IsJoinedMember = false
	attendComponents, err = b.getAttendComponent(b.JoinedGuests)
	if err != nil {
		return err
	}
	contents = append(contents, attendComponents...)

	if js, err := b.context.
		GetComfirmMode().
		GetSignal(); err != nil {
		return err
	} else {
		contents = append(contents,
			linebot.GetButtonComponent(0, linebot.GetPostBackAction(
				"提交",
				js,
			), nil),
		)
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

func (b *submitActivity) getAttendComponent(members []*submitActivityJoinedMembers) ([]interface{}, error) {
	memberComponents := []interface{}{}
	if len(members) == 0 {
		memberComponents = append(memberComponents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFillerComponent(),
				linebot.GetFlexMessageTextComponent(0, "無"),
			),
		)
	}
	for id, member := range members {
		contents := []interface{}{
			linebot.GetFillerComponent(),
			linebot.GetFlexMessageTextComponent(0, member.Name),
		}
		attendButtonOption := domain.DarkButtonOption
		if member.IsAttend {
			attendButtonOption = domain.NormalButtonOption
		}
		pathValueMap := map[string]interface{}{
			"ICmdLogic.attend_index":           id,
			"ICmdLogic.is_joined_member_index": b.IsJoinedMember,
		}
		if js, err := b.context.
			GetKeyValueInputMode(pathValueMap).
			GetSignal(); err != nil {
			return nil, err
		} else {
			action := linebot.GetPostBackAction(
				"簽到",
				js,
			)
			contents = append(contents,
				linebot.GetButtonComponent(
					0,
					action,
					&attendButtonOption,
				),
			)
		}
		payButtonOption := domain.DarkButtonOption
		if member.IsPaid {
			payButtonOption = domain.NormalButtonOption
		}
		pathValueMap = map[string]interface{}{
			"ICmdLogic.is_joined_member_index": b.IsJoinedMember,
			"ICmdLogic.pay_index":              id,
		}
		if js, err := b.context.
			GetKeyValueInputMode(pathValueMap).
			GetSignal(); err != nil {
			return nil, err
		} else {
			action := linebot.GetPostBackAction(
				"收費",
				js,
			)
			contents = append(contents,
				linebot.GetButtonComponent(
					0,
					action,
					&payButtonOption,
				),
			)
		}
		memberComponents = append(memberComponents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				contents...,
			),
		)
	}

	return memberComponents, nil
}
