package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	activityDb "heroku-line-bot/storage/database/database/clubdb/table/activity"
	"heroku-line-bot/util"
	"strconv"
	"strings"
	"time"
)

type newActivity struct {
	context     domain.ICmdHandlerContext `json:"-"`
	Date        time.Time                 `json:"date"`
	Place       string                    `json:"place"`
	Description string                    `json:"description"`
	PeopleLimit *int16                    `json:"people_limit"`
	ClubSubsidy int16                     `json:"club_subsidy"`
	IsComplete  bool                      `json:"is_complete"`
	Courts      []*newActivityCourt       `json:"courts"`
}

type newActivityCourt struct {
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
	Count        int16     `json:"count"`
	PricePerHour int16     `json:"price_per_hour"`
}

func (b *newActivityCourt) cost() float64 {
	return float64(b.Count) * b.hours() * float64(b.PricePerHour)
}

func (b *newActivityCourt) hours() float64 {
	return b.ToTime.Sub(b.FromTime).Hours()
}

func (b *newActivityCourt) time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
	)
}

func (b *newActivity) Init(context domain.ICmdHandlerContext) error {
	nowTime := commonLogic.TimeUtilObj.Now()
	const PEOPLE_PER_HOUR = 4
	*b = newActivity{
		context:     context,
		Date:        util.DateOf(nowTime),
		Place:       "大墩羽球館",
		Description: "7人出團",
		IsComplete:  false,
		Courts: []*newActivityCourt{
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
	totalHours := 0.0
	for _, court := range b.Courts {
		totalHours = commonLogic.FloatPlus(totalHours, court.hours()*float64(court.Count))
	}
	b.PeopleLimit = util.GetInt16P(int16(totalHours * float64(PEOPLE_PER_HOUR)))

	return nil
}

func (b *newActivity) GetSingleParam(attr string) string {
	switch attr {
	case "date":
		return b.Date.Format(commonLogicDomain.DATE_FORMAT)
	case "ICmdLogic.place":
		return b.Place
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

func (b *newActivity) LoadSingleParam(attr, text string) (resultValue interface{}, resultErr error) {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			return nil, err
		}
		b.Date = t
		resultValue = b.Date
	case "ICmdLogic.place":
		b.Place = text
		resultValue = b.Place
	case "ICmdLogic.description":
		b.Description = text
		resultValue = b.Description
	case "ICmdLogic.people_limit":
		i, err := strconv.Atoi(text)
		if err != nil {
			return nil, err
		}
		b.PeopleLimit = util.GetInt16P(int16(i))
		resultValue = b.PeopleLimit
	case "ICmdLogic.club_subsidy":
		i, err := strconv.Atoi(text)
		if err != nil {
			return nil, err
		}
		b.ClubSubsidy = int16(i)
		resultValue = b.ClubSubsidy
	case "ICmdLogic.courts":
		if isJson := strings.ContainsAny(text, "{"); !isJson {
			if err := b.parseCourts(text); err != nil {
				return nil, err
			}
		}
		resultValue = b.Courts
	default:
	}

	return
}

func (b *newActivity) Do(text string) (resultErr error) {
	courtsStr := b.getCourtsStr()
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

			if resultErr = transaction.Rollback().Error; resultErr != nil {
				return
			}
		}()
		data := &activityDb.ActivityTable{
			Date:          b.Date,
			Place:         b.Place,
			CourtsAndTime: courtsStr,
			ClubSubsidy:   b.ClubSubsidy,
			Description:   b.Description,
			PeopleLimit:   b.PeopleLimit,
			IsComplete:    b.IsComplete,
		}
		if resultErr = database.Club.Activity.BaseTable.Insert(transaction, data); resultErr != nil {
			return
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

	if err := b.context.CacheParams(); err != nil {
		return err
	}

	contents := []interface{}{}
	actions := domain.NewActivityLineTemplate{}

	cmd := domain.DATE_POSTBACK_CMD
	if js, err := b.context.GetRequireInputCmdText(&cmd, "date", "日期", true); err != nil {
		return err
	} else {
		actions.DateAction = linebot.GetTimeAction(
			"修改",
			js,
			"",
			"",
			linebotDomain.DATE_TIME_ACTION_MODE,
		)
	}

	if js, err := b.context.GetRequireInputCmdText(nil, "ICmdLogic.place", "地點", false); err != nil {
		return err
	} else {
		actions.PlaceAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.context.GetRequireInputCmdText(nil, "ICmdLogic.club_subsidy", "補助額", false); err != nil {
		return err
	} else {
		actions.ClubSubsidyAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.context.GetRequireInputCmdText(nil, "ICmdLogic.people_limit", "人數上限", false); err != nil {
		return err
	} else {
		actions.PeopleLimitAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.context.GetRequireInputCmdText(nil, "ICmdLogic.courts", "場地", false); err != nil {
		return err
	} else {
		actions.CourtAction = linebot.GetPostBackAction(
			"修改場地",
			js,
		)
	}

	lineContents := b.getLineComponents(actions)
	contents = append(contents, lineContents...)

	contents = append(contents,
		linebot.GetComfirmComponent(
			linebot.GetPostBackAction(
				"取消",
				b.context.GetCancelSignl(),
			),
			linebot.GetPostBackAction(
				"新增",
				b.context.GetComfirmSignl(),
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
		return err
	}

	return nil
}

func (b *newActivity) getLineComponents(actions domain.NewActivityLineTemplate) (result []interface{}) {
	result = []interface{}{}
	valueText := fmt.Sprintf("%s(%s)", b.Date.Format(commonLogicDomain.DATE_FORMAT), commonLogic.WeekDayName(b.Date.Weekday()))
	valueTextSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	result = append(result,
		linebot.GetKeyValueEditComponent(
			"日期",
			valueText,
			actions.DateAction,
			nil, &valueTextSize,
		),
	)

	result = append(result,
		linebot.GetKeyValueEditComponent(
			"地點",
			b.Place,
			actions.PlaceAction,
			nil, nil,
		),
	)

	result = append(result,
		linebot.GetKeyValueEditComponent(
			"補助額",
			strconv.Itoa(int(b.ClubSubsidy)),
			actions.ClubSubsidyAction,
			nil, nil,
		),
	)

	if b.PeopleLimit != nil {
		result = append(result,
			linebot.GetKeyValueEditComponent(
				"人數上限",
				strconv.Itoa(int(*b.PeopleLimit)),
				actions.PeopleLimitAction,
				nil, nil,
			),
		)
	}

	result = append(result, b.getCourtsBoxComponent(actions.CourtAction))

	return
}

func (b *newActivity) getCourtFee() float64 {
	totalFee := 0.0
	for _, court := range b.Courts {
		cost := court.cost()
		totalFee = commonLogic.FloatPlus(totalFee, cost)
	}
	return totalFee
}

func (b *newActivity) getCourtHours() float64 {
	totalHours := 0.0
	for _, court := range b.Courts {
		hours := court.hours()
		totalHours = commonLogic.FloatPlus(totalHours, hours)
	}
	return totalHours
}

func (b *newActivity) getCourtsStr() string {
	courtStrs := []string{}
	for _, court := range b.Courts {
		courtStr := fmt.Sprintf(
			"%d-%d-%s~%s",
			court.Count,
			court.PricePerHour,
			court.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
			court.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		)
		courtStrs = append(courtStrs, courtStr)
	}
	return strings.Join(courtStrs, ",")
}

func (b *newActivity) parseCourts(courtsStr string) error {
	b.Courts = make([]*newActivityCourt, 0)
	courtsStrs := strings.Split(courtsStr, ",")
	for _, courtsStr := range courtsStrs {
		court := &newActivityCourt{}
		timeStr := ""
		if _, err := fmt.Sscanf(
			courtsStr,
			"%d-%d-%s",
			&court.Count,
			&court.PricePerHour,
			&timeStr); err != nil {
			return err
		}
		times := strings.Split(timeStr, "~")
		if len(times) != 2 {
			return fmt.Errorf("時間格式錯誤")
		}
		fromTimeStr := times[0]
		toTimeStr := times[1]
		if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, fromTimeStr); err != nil {
			return err
		} else {
			court.FromTime = t
		}
		if t, err := time.Parse(commonLogicDomain.TIME_HOUR_MIN_FORMAT, toTimeStr); err != nil {
			return err
		} else {
			court.ToTime = t
		}

		b.Courts = append(b.Courts, court)
	}

	return nil
}

func (b *newActivity) getCourtTimeRange() (minTime, maxTime *time.Time) {
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

func (b *newActivity) getCourtsBoxComponent(buttonAction *linebotModel.PostBackAction) *linebotModel.FlexMessageBoxComponent {
	components := []interface{}{}

	headComponents := []interface{}{}
	titleComponent := linebot.GetFlexMessageTextComponent(
		0,
		linebot.GetFlexMessageTextComponentSpan(
			"場地",
			linebotDomain.XL_FLEX_MESSAGE_SIZE,
			linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT,
		),
	)
	headComponents = append(headComponents, titleComponent)
	if buttonAction != nil {
		editButtonComponent := linebot.GetButtonComponent(
			0,
			buttonAction,
		)
		headComponents = append(headComponents, editButtonComponent)
	}
	headBoxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		headComponents...,
	)
	components = append(components, headBoxComponent)

	placeFee := 0.0
	mdSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	for index, court := range b.Courts {
		cost := court.cost()
		placeFee = commonLogic.FloatPlus(placeFee, cost)

		components = append(components, linebot.GetKeyValueEditComponent(
			"時間",
			court.time(),
			nil,
			&mdSize,
			nil,
		))

		courtComponents := []interface{}{}
		courtComponents = append(courtComponents, linebot.GetKeyValueEditComponent(
			"場地數",
			strconv.Itoa(int(court.Count)),
			nil,
			&mdSize,
			nil,
		))
		courtComponents = append(courtComponents, linebot.GetKeyValueEditComponent(
			"價錢",
			strconv.FormatFloat(cost, 'f', 0, 64),
			nil,
			&mdSize,
			nil,
		))
		courtBoxComponent := linebot.GetFlexMessageBoxComponent(
			linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
			nil,
			courtComponents...,
		)

		components = append(components, courtBoxComponent)

		if index < len(b.Courts)-1 {
			components = append(components, linebot.GetSeparatorComponent(nil))
		}
	}

	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		nil,
		components...,
	)
}
