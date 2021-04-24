package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	"heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	activityDb "heroku-line-bot/storage/database/database/clubdb/table/activity"
	"heroku-line-bot/util"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type NewActivity struct {
	Context     domain.ICmdHandlerContext `json:"-"`
	Date        time.Time                 `json:"date"`
	Place       string                    `json:"place"`
	Description string                    `json:"description"`
	PeopleLimit *int16                    `json:"people_limit"`
	ClubSubsidy int16                     `json:"club_subsidy"`
	IsComplete  bool                      `json:"is_complete"`
	Courts      []*ActivityCourt          `json:"courts"`
}

type ActivityCourt struct {
	FromTime     time.Time `json:"from_time"`
	ToTime       time.Time `json:"to_time"`
	Count        int16     `json:"count"`
	PricePerHour float64   `json:"price_per_hour"`
}

func (b *ActivityCourt) cost() float64 {
	return b.TotalHours() * b.PricePerHour
}

func (b *ActivityCourt) hours() float64 {
	return b.ToTime.Sub(b.FromTime).Hours()
}

func (b *ActivityCourt) TotalHours() float64 {
	return b.hours() * float64(b.Count)
}

func (b *ActivityCourt) time() string {
	return fmt.Sprintf(
		"%s~%s",
		b.FromTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
		b.ToTime.Format(commonLogicDomain.TIME_HOUR_MIN_FORMAT),
	)
}

func (b *NewActivity) Init(context domain.ICmdHandlerContext) error {
	nowTime := commonLogic.TimeUtilObj.Now()
	*b = NewActivity{
		Context:     context,
		Date:        util.DateOf(nowTime),
		Place:       "大墩羽球館",
		Description: "7人出團",
		IsComplete:  false,
		Courts: []*ActivityCourt{
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
		totalHours = commonLogic.FloatPlus(totalHours, court.TotalHours())
	}
	b.PeopleLimit = util.GetInt16P(int16(totalHours * float64(domain.PEOPLE_PER_HOUR)))

	return nil
}

func (b *NewActivity) GetSingleParam(attr string) string {
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

func (b *NewActivity) LoadSingleParam(attr, text string) error {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			return err
		}
		b.Date = t
	case "ICmdLogic.place":
		b.Place = text
	case "ICmdLogic.description":
		b.Description = text
	case "ICmdLogic.people_limit":
		i, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		b.PeopleLimit = util.GetInt16P(int16(i))
	case "ICmdLogic.club_subsidy":
		i, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		b.ClubSubsidy = int16(i)
	case "ICmdLogic.courts":
		if isJson := strings.ContainsAny(text, "{"); !isJson {
			if err := b.ParseCourts(text); err != nil {
				return err
			}
		}
	default:
	}

	return nil
}

func (b *NewActivity) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *NewActivity) Do(text string) (resultErr error) {
	if u, err := lineuser.Get(b.Context.GetUserID()); err != nil {
		return err
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE &&
			u.Role != domain.CADRE_CLUB_ROLE {
			return domain.NO_AUTH_ERROR
		}
	}

	if b.Context.IsComfirmed() {
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
		if resultErr = b.InsertActivity(transaction); resultErr != nil {
			return
		}

		if resultErr = b.Context.DeleteParam(); resultErr != nil {
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if resultErr = b.Context.Reply(replyMessges); resultErr != nil {
			return resultErr
		}

		return nil
	}

	if err := b.Context.CacheParams(); err != nil {
		return err
	}

	contents := []interface{}{}
	actions := domain.NewActivityLineTemplate{}

	if js, err := b.Context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); err != nil {
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

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.place", "地點", false).
		GetSignal(); err != nil {
		return err
	} else {
		actions.PlaceAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.club_subsidy", "補助額", false).
		GetSignal(); err != nil {
		return err
	} else {
		actions.ClubSubsidyAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.people_limit", "人數上限", false).
		GetSignal(); err != nil {
		return err
	} else {
		actions.PeopleLimitAction = linebot.GetPostBackAction(
			"修改",
			js,
		)
	}

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.courts", "場地", false).
		GetSignal(); err != nil {
		return err
	} else {
		actions.CourtAction = linebot.GetPostBackAction(
			"修改場地",
			js,
		)
	}

	lineContents := b.getLineComponents(actions)
	contents = append(contents, lineContents...)

	cancelSignlJs, err := b.Context.
		GetCancelMode().
		GetSignal()
	if err != nil {
		return err
	}
	comfirmSignlJs, err := b.Context.
		GetComfirmMode().
		GetSignal()
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (b *NewActivity) InsertActivity(transaction *gorm.DB) (resultErr error) {
	courtsStr := b.getCourtsStr()
	if transaction == nil {
		transaction = database.Club.Begin()
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
	}

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

	return nil
}

func (b *NewActivity) getLineComponents(actions domain.NewActivityLineTemplate) (result []interface{}) {
	result = []interface{}{}
	valueText := fmt.Sprintf("%s(%s)", b.Date.Format(commonLogicDomain.DATE_FORMAT), commonLogic.WeekDayName(b.Date.Weekday()))
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

	result = append(result,
		GetKeyValueEditComponent(
			"地點",
			b.Place,
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

func (b *NewActivity) getCourtFee() float64 {
	totalFee := 0.0
	for _, court := range b.Courts {
		cost := court.cost()
		totalFee = commonLogic.FloatPlus(totalFee, cost)
	}
	return totalFee
}

func (b *NewActivity) getCourtHours() float64 {
	totalHours := 0.0
	for _, court := range b.Courts {
		hours := court.hours()
		totalHours = commonLogic.FloatPlus(totalHours, hours)
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

func (b *NewActivity) ParseCourts(courtsStr string) error {
	b.Courts = make([]*ActivityCourt, 0)
	courtsStrs := strings.Split(courtsStr, ",")
	for _, courtsStr := range courtsStrs {
		court := &ActivityCourt{}
		timeStr := ""
		if _, err := fmt.Sscanf(
			courtsStr,
			"%d-%f-%s",
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
		0,
		"",
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

	placeFee := 0.0
	mdSize := linebotDomain.MD_FLEX_MESSAGE_SIZE
	keyValueEditComponentOption := &domain.KeyValueEditComponentOption{
		SizeP: &mdSize,
	}
	for index, court := range b.Courts {
		cost := court.cost()
		placeFee = commonLogic.FloatPlus(placeFee, cost)

		components = append(components, GetKeyValueEditComponent(
			"時間",
			court.time(),
			keyValueEditComponentOption,
		))

		courtBoxComponent := GetDoubleKeyValueComponent(
			"場地數",
			strconv.Itoa(int(court.Count)),
			"價錢",
			strconv.FormatFloat(cost, 'f', 0, 64),
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
		strconv.FormatFloat(courtFee, 'f', -1, 64),
		keyValueEditComponentOption,
	)
	components = append(components, courtFeeComponent)

	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		nil,
		components...,
	)
}
