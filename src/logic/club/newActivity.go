package club

import (
	"fmt"
	courtDomain "heroku-line-bot/src/logic/badminton/court/domain"
	badmintonPlaceLogic "heroku-line-bot/src/logic/badminton/place"
	"heroku-line-bot/src/logic/club/domain"
	commonLogic "heroku-line-bot/src/logic/common"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	linebotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/activity"
	"heroku-line-bot/src/repo/database/database/clubdb/place"
	"strconv"
	"strings"
	"time"
)

type NewActivity struct {
	context domain.ICmdHandlerContext `json:"-"`
	domain.TimePostbackParams
	PlaceID     int                          `json:"place_id"`
	TeamID      int                          `json:"team_id"`
	Description string                       `json:"description"`
	PeopleLimit *int16                       `json:"people_limit"`
	ClubSubsidy int16                        `json:"club_subsidy"`
	Courts      []*courtDomain.ActivityCourt `json:"courts"`
}

func (b *NewActivity) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	nowTime := global.TimeUtilObj.Now()
	*b = NewActivity{
		context: context,
		TimePostbackParams: domain.TimePostbackParams{
			Date: *util.NewDateTimePOf(&nowTime),
		},
		PlaceID:     1,
		Description: "7人出團",
		Courts: []*courtDomain.ActivityCourt{
			{
				FromTime:     commonLogic.GetTime(1, 1, 1, 18),
				ToTime:       commonLogic.GetTime(1, 1, 1, 20, 30),
				Count:        1,
				PricePerHour: 480,
			},
			{
				FromTime:     commonLogic.GetTime(1, 1, 1, 19, 30),
				ToTime:       commonLogic.GetTime(1, 1, 1, 20, 30),
				Count:        1,
				PricePerHour: 480,
			},
		},
		TeamID: clubTeamID,
	}
	totalHours := b.getCourtHours()
	b.PeopleLimit = util.GetInt16P(int16(totalHours.MulFloat(float64(domain.PEOPLE_PER_HOUR)).ToInt()))

	return nil
}

func (b *NewActivity) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *NewActivity) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	case "date":
		valueText = b.Date.Time().Format(util.DATE_FORMAT)
	case "ICmdLogic.place_id":
		attrNameText = "地點"

		if dbDatas, errInfo := badmintonPlaceLogic.Load(b.PlaceID); errInfo == nil || !errInfo.IsError() {
			for _, v := range dbDatas {
				valueText = v.Name
				return
			}
		}
		valueText = "未設置"
	case "ICmdLogic.補助額":
		attrNameText = "數量"
	case "ICmdLogic.people_limit":
		attrNameText = "人數上限"

		if b.PeopleLimit == nil {
			valueText = "未設置"
		} else {
			valueText = strconv.Itoa(int(*b.PeopleLimit))
		}
	case "ICmdLogic.courts":
		attrNameText = "場地"

		valueText = "場數-每場價錢-hh:mm~hh:mm"
	case "ICmdLogic.description":
		valueText = b.Description
	case "ICmdLogic.club_subsidy":
		valueText = strconv.Itoa(int(b.ClubSubsidy))
	}
	return
}

func (b *NewActivity) GetInputTemplate(attr string) (messages interface{}) {
	return
}

func (b *NewActivity) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(util.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.Date = util.DateTime(t)
	case "ICmdLogic.place_id":
		if dbDatas, err := database.Club().Place.Select(
			place.Reqs{
				Name: &text,
			},
			place.COLUMN_ID,
			place.COLUMN_Name,
		); err != nil {
			errInfo := errUtil.NewError(err)
			if resultErrInfo == nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
		} else if len(dbDatas) == 0 {
			errInfo := errUtil.New("未登記的球場")
			if resultErrInfo == nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
			return
		} else {
			for _, v := range dbDatas {
				b.PlaceID = v.ID
			}
		}
	case "ICmdLogic.description":
		b.Description = text
	case "ICmdLogic.people_limit":
		i, err := strconv.Atoi(text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.PeopleLimit = util.GetInt16P(int16(i))
	case "ICmdLogic.club_subsidy":
		i, err := strconv.Atoi(text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.ClubSubsidy = int16(i)
	case "ICmdLogic.courts":
		if isJson := strings.ContainsAny(text, "{"); !isJson {
			if errInfo := b.ParseCourts(text); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				return
			}
		}
	default:
	}

	return nil
}

func (b *NewActivity) Do(text string) (resultErrInfo errUtil.IError) {
	if user, isAutoRegiste, errInfo := autoRegiste(b.context); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	} else if isAutoRegiste {
		replyMessges := autoRegisteMessage()
		if err := b.context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
	} else if user.Role != domain.ADMIN_CLUB_ROLE &&
		user.Role != domain.CADRE_CLUB_ROLE {
		resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
		return
	}

	if b.context.IsConfirmed() {
		db, transaction, err := database.Club().Begin()
		if err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		defer func() {
			if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
				if resultErrInfo.IsError() {
					return
				}
			}
		}()

		if errInfo := b.InsertActivity(&db); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}

		if err := b.context.DeleteParam(); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if err := b.context.Reply(replyMessges); err != nil {
			errInfo := errUtil.NewError(err)
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}

		return
	}

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if resultErrInfo.IsError() {
			return
		}
	}

	contents := []interface{}{}
	actions := domain.NewActivityLineTemplate{}

	if js, errInfo := NewSignal().
		GetBasePath("ICmdLogic").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		actions.DateAction = linebot.GetTimeAction(
			"修改",
			js,
			"",
			"",
			linebotDomain.DATE_TIME_ACTION_MODE,
		)
	}

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.place_id").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		actions.PlaceAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.club_subsidy").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		actions.ClubSubsidyAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.people_limit").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		actions.PeopleLimitAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.courts").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		actions.CourtAction = linebot.GetPostBackAction(
			"修改場地",
			js,
		)
	}

	lineContents := b.getLineComponents(actions)
	contents = append(contents, lineContents...)

	cancelSignlJs, errInfo := NewSignal().
		GetCancelMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	comfirmSignlJs, errInfo := NewSignal().
		GetConfirmMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}
	contents = append(contents,
		GetConfirmComponent(
			linebot.GetPostBackAction(
				"取消",
				cancelSignlJs,
			),
			linebot.GetPostBackAction(
				"新增",
				comfirmSignlJs,
			),
		),
	)

	replyMessges := []interface{}{
		linebot.GetFlexMessage(
			"新增活動",
			linebot.GetFlexMessageBubbleContent(
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					contents...,
				),
				nil,
			),
		),
	}
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *NewActivity) InsertActivity(db *clubdb.Database) (resultErrInfo errUtil.IError) {
	courtsStr := b.getCourtsStr()
	if db == nil {
		dbConn, transaction, err := database.Club().Begin()
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		defer func() {
			if errInfo := database.CommitTransaction(transaction, resultErrInfo); errInfo != nil {
				resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			}
		}()
		db = &dbConn
	}

	data := &activity.Model{
		Date:          b.Date.Time(),
		PlaceID:       b.PlaceID,
		CourtsAndTime: courtsStr,
		ClubSubsidy:   b.ClubSubsidy,
		Description:   b.Description,
		PeopleLimit:   b.PeopleLimit,
		TeamID:        b.TeamID,
	}
	if err := db.Activity.Insert(data); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return
}

func (b *NewActivity) getPlaceTimeTemplate() (result []interface{}) {
	result = []interface{}{}

	place := "無球館"
	if dbDatas, errInfo := badmintonPlaceLogic.Load(b.PlaceID); errInfo == nil || !errInfo.IsError() {
		for _, v := range dbDatas {
			place = v.Name
		}
	}
	result = append(result,
		linebot.GetFlexMessageTextComponent(
			place,
			&linebotModel.FlexMessageTextComponentOption{
				Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				Size:   linebotDomain.XXL_FLEX_MESSAGE_SIZE,
				Margin: linebotDomain.MD_FLEX_MESSAGE_SIZE,
			},
		),
	)

	minTime, maxTime := b.getCourtTimeRange()
	valueText := fmt.Sprintf("%s(%s) %s~%s",
		b.Date.Time().Format(util.DATE_FORMAT),
		util.GetWeekDayName(b.Date.Time().Weekday()),
		minTime.Format(util.TIME_HOUR_MIN_FORMAT),
		maxTime.Format(util.TIME_HOUR_MIN_FORMAT),
	)
	result = append(result,
		linebot.GetFlexMessageTextComponent(
			valueText,
			&linebotModel.FlexMessageTextComponentOption{
				Size:  linebotDomain.XS_FLEX_MESSAGE_SIZE,
				Color: "#aaaaaa",
				Wrap:  true,
			},
		),
	)

	return
}

func (b *NewActivity) getLineComponents(actions domain.NewActivityLineTemplate) (result []interface{}) {
	result = []interface{}{}
	valueText := fmt.Sprintf("%s(%s)", b.Date.Time().Format(util.DATE_FORMAT), util.GetWeekDayName(b.Date.Time().Weekday()))
	valueTextSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	result = append(result,
		GetKeyValueEditComponent(
			"日期",
			valueText,
			&domain.KeyValueEditComponentOption{
				Action:     actions.DateAction,
				ValueSizeP: &valueTextSize,
			},
		),
	)

	place := "無球館"
	if dbDatas, errInfo := badmintonPlaceLogic.Load(b.PlaceID); errInfo == nil || !errInfo.IsError() {
		for _, v := range dbDatas {
			place = v.Name
		}
	}
	result = append(result,
		GetKeyValueEditComponent(
			"地點",
			place,
			&domain.KeyValueEditComponentOption{
				Action: actions.PlaceAction,
			},
		),
	)

	result = append(result,
		GetKeyValueEditComponent(
			"補助額",
			strconv.Itoa(int(b.ClubSubsidy)),
			&domain.KeyValueEditComponentOption{
				Action: actions.ClubSubsidyAction,
			},
		),
	)

	if b.PeopleLimit != nil {
		result = append(result,
			GetKeyValueEditComponent(
				"人數上限",
				strconv.Itoa(int(*b.PeopleLimit)),
				&domain.KeyValueEditComponentOption{
					Action: actions.PeopleLimitAction,
				},
			),
		)
	}

	result = append(result, b.getCourtsBoxComponent(actions.CourtAction))

	return
}

func (b *NewActivity) getCourtFee() util.Float {
	totalFee := util.NewFloat(0)
	for _, court := range b.Courts {
		cost := court.Cost()
		totalFee = totalFee.Plus(cost)
	}
	return totalFee
}

func (b *NewActivity) getCourtHours() util.Float {
	totalHours := util.NewFloat(0)
	for _, court := range b.Courts {
		hours := court.TotalHours()
		totalHours = totalHours.Plus(hours)
	}
	return totalHours
}

func (b *NewActivity) getCourtsStr() string {
	courtStrs := []string{}
	for _, court := range b.Courts {
		courtStr := fmt.Sprintf(
			"%d-%.1f-%s~%s",
			court.Count,
			court.PricePerHour,
			court.FromTime.Format(util.TIME_HOUR_MIN_FORMAT),
			court.ToTime.Format(util.TIME_HOUR_MIN_FORMAT),
		)
		courtStrs = append(courtStrs, courtStr)
	}
	return strings.Join(courtStrs, ",")
}

func (b *NewActivity) ParseCourts(courtsStr string) (resultErrInfo errUtil.IError) {
	b.Courts = make([]*courtDomain.ActivityCourt, 0)
	courtsStrs := strings.Split(courtsStr, ",")
	for _, courtsStr := range courtsStrs {
		court := &courtDomain.ActivityCourt{}
		timeStr := ""
		if _, err := fmt.Sscanf(
			courtsStr,
			"%d-%f-%s",
			&court.Count,
			&court.PricePerHour,
			&timeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		times := strings.Split(timeStr, "~")
		if len(times) != 2 {
			errInfo := errUtil.New("時間格式錯誤")
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			return
		}
		fromTimeStr := times[0]
		toTimeStr := times[1]
		if t, err := time.Parse(util.TIME_HOUR_MIN_FORMAT, fromTimeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			court.FromTime = t
		}
		if t, err := time.Parse(util.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			court.ToTime = t
		}

		b.Courts = append(b.Courts, court)
	}

	return nil
}

func (b *NewActivity) getCourtTimeRange() (minTime, maxTime *time.Time) {
	for _, court := range b.Courts {
		if minTime == nil || court.FromTime.Before(*minTime) {
			minTime = &court.FromTime
		}
		if maxTime == nil || court.ToTime.After(*maxTime) {
			maxTime = &court.ToTime
		}
	}
	return
}

func (b *NewActivity) getCourtsBoxComponent(buttonAction *linebotModel.PostBackAction) *linebotModel.FlexMessageBoxComponent {
	components := []interface{}{}

	headComponents := []interface{}{}
	titleComponent := linebot.GetFlexMessageTextComponent(
		"",
		&linebotModel.FlexMessageTextComponentOption{
			Contents: []*linebotModel.FlexMessageTextComponentSpan{
				linebot.GetFlexMessageTextComponentSpan(
					"場地",
					linebotDomain.XL_FLEX_MESSAGE_SIZE,
					linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				),
			},
			AdjustMode: linebotDomain.SHRINK_TO_FIT_ADJUST_MODE,
			Align:      linebotDomain.START_Align,
		},
	)
	headComponents = append(headComponents, titleComponent)
	if buttonAction != nil {
		editButtonComponent := linebot.GetButtonComponent(
			buttonAction,
			&domain.NormalButtonOption,
		)
		headComponents = append(headComponents, editButtonComponent)
	}
	headBoxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		headComponents...,
	)
	components = append(components, headBoxComponent)

	mdSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	keyValueEditComponentOption := &domain.KeyValueEditComponentOption{
		SizeP: &mdSize,
	}
	for index, court := range b.Courts {
		cost := court.Cost()

		components = append(components, GetKeyValueEditComponent(
			"時間",
			court.Time(),
			keyValueEditComponentOption,
		))

		courtBoxComponent := GetDoubleKeyValueComponent(
			"場地數",
			strconv.Itoa(int(court.Count)),
			"價錢",
			cost.ToString(0),
			nil,
			keyValueEditComponentOption,
		)
		components = append(components, courtBoxComponent)

		if index < len(b.Courts)-1 {
			components = append(components, linebot.GetSeparatorComponent(nil))
		}
	}

	courtFee := b.getCourtFee()
	courtFeeComponent := GetKeyValueEditComponent(
		"場地費用總計",
		courtFee.ToString(-1),
		keyValueEditComponentOption,
	)
	components = append(components, courtFeeComponent)

	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		nil,
		components...,
	)
}

func (b *NewActivity) getCourtsContents() []interface{} {
	courtFee := b.getCourtFee()
	contents := []interface{}{
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin: linebotDomain.XXL_FLEX_MESSAGE_SIZE,
			},
			linebot.GetFlexMessageTextComponent(
				"場地",
				&linebotModel.FlexMessageTextComponentOption{
					Size:   linebotDomain.MD_FLEX_MESSAGE_SIZE,
					Weight: linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", courtFee.ToString(-1)),
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

	courtContents := make([]interface{}, 0)
	for _, court := range b.Courts {
		cost := court.Cost()

		courtsComponent := linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			nil,
			linebot.GetFlexMessageTextComponent(
				court.Time(),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#555555",
					Flex:  0,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("%d場", court.Count),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.SM_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.CENTER_Align,
				},
			),
			linebot.GetFlexMessageTextComponent(
				fmt.Sprintf("$%s", cost.ToString(0)),
				&linebotModel.FlexMessageTextComponentOption{
					Size:  linebotDomain.XS_FLEX_MESSAGE_SIZE,
					Color: "#111111",
					Align: linebotDomain.END_Align,
				},
			),
		)
		courtContents = append(courtContents, courtsComponent)
	}

	contents = append(
		contents,
		linebot.GetFlexMessageBoxComponent(
			linebotDomain.VERTICAL_MESSAGE_LAYOUT,
			&linebotModel.FlexMessageBoxComponentOption{
				Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
				Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
			},
			courtContents...,
		),
	)

	return contents
}
