package club

import (
	"fmt"
	"heroku-line-bot/logic/club/domain"
	commonLogic "heroku-line-bot/logic/common"
	commonLogicDomain "heroku-line-bot/logic/common/domain"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/logic/redis/lineuser"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
	"heroku-line-bot/storage/database"
	logisticDb "heroku-line-bot/storage/database/database/clubdb/table/logistic"
	"heroku-line-bot/util"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type NewLogistic struct {
	Context     domain.ICmdHandlerContext `json:"-"`
	Date        time.Time                 `json:"date"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Amount      int16                     `json:"amount"`
}

func (b *NewLogistic) Init(context domain.ICmdHandlerContext) (resultErrInfo errLogic.IError) {
	nowTime := commonLogic.TimeUtilObj.Now()
	*b = NewLogistic{
		Context:     context,
		Date:        util.DateOf(nowTime),
		Name:        domain.BALL_NAME,
		Description: "買球 https://shopee.tw/product/4013408/4461135276",
		Amount:      180,
	}

	return nil
}

func (b *NewLogistic) GetSingleParam(attr string) string {
	switch attr {
	case "date":
		return b.Date.Format(commonLogicDomain.DATE_FORMAT)
	case "ICmdLogic.name":
		return b.Name
	case "ICmdLogic.description":
		return b.Description
	case "ICmdLogic.amount":
		return strconv.Itoa(int(b.Amount))
	default:
		return ""
	}
}

func (b *NewLogistic) LoadSingleParam(attr, text string) (resultErrInfo errLogic.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		b.Date = t
	case "ICmdLogic.name":
		b.Name = text
	case "ICmdLogic.description":
		b.Description = text
	case "ICmdLogic.amount":
		i, err := strconv.Atoi(text)
		if err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		b.Amount = int16(i)
	default:
	}

	return nil
}

func (b *NewLogistic) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *NewLogistic) Do(text string) (resultErrInfo errLogic.IError) {
	if u, err := lineuser.Get(b.Context.GetUserID()); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE {
			resultErrInfo = errLogic.NewError(domain.NO_AUTH_ERROR)
			return
		}
	}

	if b.Context.IsComfirmed() {
		transaction := database.Club.Begin()
		if err := transaction.Error; err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		defer database.CommitTransaction(transaction, resultErrInfo)

		if resultErrInfo = b.InsertLogistic(transaction); resultErrInfo != nil {
			return
		}

		if err := b.Context.DeleteParam(); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		replyMessges := []interface{}{
			linebot.GetTextMessage("完成"),
		}
		if err := b.Context.Reply(replyMessges); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}

		return
	}

	if errInfo := b.Context.CacheParams(); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	boxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
	)

	if js, errInfo := b.Context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		dateStr := fmt.Sprintf("%s(%s)",
			b.Date.Format(commonLogicDomain.MONTH_DATE_SLASH_FORMAT),
			commonLogic.WeekDayName(b.Date.Weekday()),
		)
		boxComponent.Contents = append(boxComponent.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					dateStr,
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.XL_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: "#FFBF00",
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						linebot.GetTimeAction(
							"修改",
							js,
							"",
							"",
							linebotDomain.DATE_TIME_ACTION_MODE,
						),
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
			),
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.name", "品項", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		boxComponent.Contents = append(boxComponent.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					b.Name,
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.XL_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: "#FFBF00",
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						linebot.GetPostBackAction(
							"修改",
							js,
						),
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
			),
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.amount", "數量", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		amountStr := fmt.Sprintf("%d個, %d打", b.Amount, b.Amount/12)
		boxComponent.Contents = append(boxComponent.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					amountStr,
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.XL_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: "#FFBF00",
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						linebot.GetPostBackAction(
							"修改",
							js,
						),
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
			),
		)
	}

	if js, errInfo := b.Context.
		GetRequireInputMode("ICmdLogic.description", "備註", false).
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		boxComponent.Contents = append(boxComponent.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				nil,
				linebot.GetFlexMessageTextComponent(
					b.Description,
					&linebotModel.FlexMessageTextComponentOption{
						Size:  linebotDomain.XL_FLEX_MESSAGE_SIZE,
						Color: "#555555",
					},
				),
				linebot.GetFlexMessageBoxComponent(
					linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
					&linebotModel.FlexMessageBoxComponentOption{
						BackgroundColor: "#FFBF00",
						CornerRadius:    "12px",
					},
					linebot.GetButtonComponent(
						linebot.GetPostBackAction(
							"修改",
							js,
						),
						&linebotModel.ButtonOption{
							Color:  domain.WHITE_COLOR,
							Height: linebotDomain.SM_FLEX_MESSAGE_SIZE,
						},
					),
				),
			),
		)
	}

	if js, errInfo := b.Context.
		GetComfirmMode().
		GetSignal(); errInfo != nil {
		resultErrInfo = errInfo
		return
	} else {
		boxComponent.Contents = append(boxComponent.Contents,
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
				&linebotModel.FlexMessageBoxComponentOption{
					BackgroundColor: "#1E90FF",
					CornerRadius:    "12px",
				},
				linebot.GetButtonComponent(
					linebot.GetPostBackAction(
						"新增",
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

	replyMessges := []interface{}{
		linebot.GetFlexMessage(
			"新增物品紀錄",
			linebot.GetFlexMessageBubbleContent(
				boxComponent,
				nil,
			),
		),
	}
	if err := b.Context.Reply(replyMessges); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return nil
}

func (b *NewLogistic) InsertLogistic(transaction *gorm.DB) (resultErrInfo errLogic.IError) {
	if transaction == nil {
		transaction = database.Club.Begin()
		if err := transaction.Error; err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		defer database.CommitTransaction(transaction, resultErrInfo)
	}

	data := &logisticDb.LogisticTable{
		Date:        b.Date,
		Name:        b.Name,
		Amount:      b.Amount,
		Description: b.Description,
	}
	if err := database.Club.Logistic.BaseTable.Insert(transaction, data); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	}

	return nil
}
