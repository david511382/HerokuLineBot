package club

import (
	"fmt"
	"heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	linebotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	"heroku-line-bot/src/pkg/util"
	errUtil "heroku-line-bot/src/pkg/util/error"
	"heroku-line-bot/src/repo/database"
	"heroku-line-bot/src/repo/database/database/clubdb"
	"heroku-line-bot/src/repo/database/database/clubdb/logistic"
	"strconv"
	"time"
)

type NewLogistic struct {
	context domain.ICmdHandlerContext `json:"-"`
	domain.TimePostbackParams
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      int16  `json:"amount"`
	TeamID      int    `json:"team_id"`
}

func (b *NewLogistic) Init(context domain.ICmdHandlerContext) (resultErrInfo errUtil.IError) {
	nowTime := global.TimeUtilObj.Now()
	*b = NewLogistic{
		context: context,
		TimePostbackParams: domain.TimePostbackParams{
			Date: *util.NewDateTimePOf(&nowTime),
		},
		Name:        domain.BALL_NAME,
		Description: "買球 https://shopee.tw/product/4013408/4461135276",
		Amount:      180,
		TeamID:      clubTeamID,
	}

	return nil
}

func (b *NewLogistic) GetRequireAttr() (requireAttr string, warnMessage interface{}, resultErrInfo errUtil.IError) {
	return
}

func (b *NewLogistic) GetRequireAttrInfo(rawAttr string) (attrNameText string, valueText string, isNotRequireChecking bool) {
	switch rawAttr {
	case "date":
		valueText = b.Date.Time().Format(util.DATE_FORMAT)
	case "ICmdLogic.name":
		attrNameText = "品項"
		valueText = b.Name
	case "ICmdLogic.amount":
		attrNameText = "數量"
		valueText = strconv.Itoa(int(b.Amount))
	case "ICmdLogic.description":
		attrNameText = "備註"
		valueText = b.Description
	}
	return
}

func (b *NewLogistic) GetInputTemplate(attr string) (messages interface{}) {
	return
}

func (b *NewLogistic) LoadRequireInputTextParam(attr, text string) (resultErrInfo errUtil.IError) {
	switch attr {
	case "date":
		t, err := time.Parse(util.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.Date = *util.NewDateTimePOf(&t)
	case "ICmdLogic.name":
		b.Name = text
	case "ICmdLogic.description":
		b.Description = text
	case "ICmdLogic.amount":
		i, err := strconv.Atoi(text)
		if err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		b.Amount = int16(i)
	default:
	}

	return nil
}

func (b *NewLogistic) Do(text string) (resultErrInfo errUtil.IError) {
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
	} else if user.Role != domain.ADMIN_CLUB_ROLE {
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
			}
		}()

		if resultErrInfo = b.InsertLogistic(&db); resultErrInfo != nil {
			return
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

		return
	}

	if errInfo := b.context.CacheParams(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	}

	boxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
	)

	if js, errInfo := NewSignal().
		GetBasePath("ICmdLogic").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		return
	} else {
		dateStr := fmt.Sprintf("%s(%s)",
			b.Date.Time().Format(util.MONTH_DATE_SLASH_FORMAT),
			util.GetWeekDayName(b.Date.Time().Weekday()),
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

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.name").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.amount").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

	if js, errInfo := NewSignal().
		GetRequireInputMode("ICmdLogic.description").
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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

	if js, errInfo := NewSignal().
		GetConfirmMode().
		GetSignal(); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
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
	if err := b.context.Reply(replyMessges); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}

func (b *NewLogistic) InsertLogistic(db *clubdb.Database) (resultErrInfo errUtil.IError) {
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

	data := &logistic.Model{
		Date:        b.Date.Time(),
		Name:        b.Name,
		Amount:      b.Amount,
		Description: b.Description,
		TeamID:      b.TeamID,
	}
	if err := db.Logistic.Insert(data); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	}

	return nil
}
