package club

import (
	"fmt"
	"heroku-line-bot/global"
	courtDomain "heroku-line-bot/logic/badminton/court/domain"
	badmintonPlaceLogic "heroku-line-bot/logic/badminton/place"
	"heroku-line-bot/logic/club/domain"
	clubLineuserLogic "heroku-line-bot/logic/club/lineuser"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	activityDb "heroku-line-bot/storage/database/database/clubdb/table/activity"
	dbReqs "heroku-line-bot/storage/database/domain/model/reqs"
	"heroku-line-bot/util"
	errUtil "heroku-line-bot/util/error"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type NewActivity struct {
	Context     domain.ICmdHandlerContext    `json:"-"`
	Date        commonLogic.DateTime         `json:"date"`
	PlaceID     int                          `json:"place_id"`
	Description string                       `json:"description"`
	PeopleLimit *int16                       `json:"people_limit"`
	ClubSubsidy int16                        `json:"club_subsidy"`
	IsComplete  bool                         `json:"is_complete"`
	Courts      []*courtDomain.ActivityCourt `json:"courts"`
}

func (b *NewActivity) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	nowTime := global.TimeUtilObj.Now()
	*b = NewActivity{
		Context:     context,
		Date:        commonLogic.DateTime(util.DateOf(nowTime)),
		PlaceID:     1,
		Description: "7人出團",
		IsComplete:  false,
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
	}
	totalHours := b.getCourtHours()
	b.PeopleLimit = util.GetInt16P(int16(totalHours.MulFloat(float64(domain.PEOPLE_PER_HOUR)).ToInt()))

	return nil
}

func (b *NewActivity) GetSingleParam(attr string) string {
	switch attr {
	case "date":
		return b.Date.Time().Format(commonLogicDomain.DATE_FORMAT)
	case "ICmdLogic.place_id":
		if dbDatas, errInfo := badmintonPlaceLogic.Load(b.PlaceID); errInfo == nil || !errInfo.IsError() {
			for _, v := range dbDatas {
				return v.Name
			}
		}
		return "未設置"
	case "ICmdLogic.description":
		return b.Description
	case "ICmdLogic.people_limit":
		if b.PeopleLimit == nil {
			return "未設置"
		} else {
			return strconv.Itoa(int(*b.PeopleLimit))
		}
	case "ICmdLogic.club_subsidy":
		return strconv.Itoa(int(b.ClubSubsidy))
	case "ICmdLogic.courts":
		return "場數-每場價錢-hh:mm~hh:mm"
	default:
		return ""
	}
}

func (b *NewActivity) LoadSingleParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.Date = commonLogic.DateTime(t)
	case "ICmdLogic.place_id":
		if dbDatas, err := database.Club.Place.IDName(dbReqs.Place{
			Name: &text,
		}); err != nil {
			errInfo := errUtil.NewError(err)
			if resultErrInfo == nil {
				resultErrInfo = errInfo
			} else {
				resultErrInfo = resultErrInfo.Append(errInfo)
			}
		} else if len(dbDatas) == 0 {
			errInfo := errUtil.New("未登記的球場")
			if resultErrInfo == nil {
				resultErrInfo = errInfo
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
				resultErrInfo = errInfo
				return
			}
		}
	default:
	}

	return nil
}

func (b *NewActivity) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *NewActivity) Do(text string) (resultErrInfo errUtil.IError) {
	if u, err := clubLineuserLogic.Get(b.Context.GetUserID()); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE &&
			u.Role != domain.CADRE_CLUB_ROLE {
			resultErrInfo = errUtil.NewError(domain.NO_AUTH_ERROR)
			return
		}
	}

	if b.Context.IsComfirmed() {
		transaction := database.Club.Begin()
		if err := transaction.Error; err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		defer database.CommitTransaction(transaction, resultErrInfo)

		if resultErrInfo = b.InsertActivity(transaction); resultErrInfo != nil {
			return
		}

		if err := b.Context.DeleteParam(); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if err := b.Context.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}

		return
	}

	if errInfo := b.Context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	contents := []interface{}{}
	actions := domain.NewActivityLineTemplate{}

	if js, errInfo := b.Context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
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

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.place_id", "地點", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		actions.PlaceAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.club_subsidy", "補助額", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		actions.ClubSubsidyAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.people_limit", "人數上限", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		actions.PeopleLimitAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.courts", "場地", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		actions.CourtAction = linebot.GetPostBackAction(
			"修改場地",
			js,
		)
	}

	lineContents := b.getLineComponents(actions)
	contents = append(contents, lineContents...)

	cancelSignlJs, errInfo := b.Context.
		GetCancelMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errInfo
		return
	}
	comfirmSignlJs, errInfo := b.Context.
		GetComfirmMode().
		GetSignal()
	if errInfo != nil {
		resultErrInfo = errInfo
		return
	}
	contents = append(contents,
		GetComfirmComponent(
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
	if err := b.Context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *NewActivity) InsertActivity(transaction *gorm.DB) (resultErrInfo errUtil.IError) {
	courtsStr := b.getCourtsStr()
	if transaction == nil {
		transaction = database.Club.Begin()
		if err := transaction.Error; err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		defer database.CommitTransaction(transaction, resultErrInfo)
	}

	data := &activityDb.ActivityTable{
		Date:          b.Date.Time(),
		PlaceID:       b.PlaceID,
		CourtsAndTime: courtsStr,
		ClubSubsidy:   b.ClubSubsidy,
		Description:   b.Description,
		PeopleLimit:   b.PeopleLimit,
		IsComplete:    b.IsComplete,
	}
	if err := database.Club.Activity.Insert(transaction, data); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
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
		b.Date.Time().Format(commonLogicDomain.DATE_FORMAT),
		commonLogic.WeekDayName(b.Date.Time().Weekday()),
		minTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		maxTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
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
	valueText := fmt.Sprintf("%s(%s)", b.Date.Time().Format(commonLogicDomain.DATE_FORMAT), commonLogic.WeekDayName(b.Date.Time().Weekday()))
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
	totalFee := util.ToFloat(0)
	for _, court := range b.Courts {
		cost := court.Cost()
		totalFee = totalFee.Plus(cost)
	}
	return totalFee
}

func (b *NewActivity) getCourtHours() util.Float {
	totalHours := util.ToFloat(0)
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
			court.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
			court.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
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
			errInfo = errInfo.Trace()
			resultErrInfo = errInfo
			return
		}
		fromTimeStr := times[0]
		toTimeStr := times[1]
		if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, fromTimeStr); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		} else {
			court.FromTime = t
		}
		if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
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
