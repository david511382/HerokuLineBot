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
	logisticDb "heroku-line-bot/storage/database/database/clubdb/table/logistic"
	"heroku-line-bot/util"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type NewLogistic struct {
	Context     domain.ICmdHandlerContext `json:"-"`
	Date        time.Time                 `json:"date"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Amount      int16                     `json:"amount"`
}

func (b *NewLogistic) Init(context domain.ICmdHandlerContext) error {
	nowTime := commonLogic.TimeUtilObj.Now()
	*b = NewLogistic{
		Context:     context,
		Date:        util.DateOf(nowTime),
		Name:        "RSL4",
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

func (b *NewLogistic) LoadSingleParam(attr, text string) error {
	switch attr {
	case "date":
		t, err := time.Parse(commonLogicDomain.DATE_TIME_RFC3339_FORMAT, text)
		if err != nil {
			return err
		}
		b.Date = t
	case "ICmdLogic.name":
		b.Name = text
	case "ICmdLogic.description":
		b.Description = text
	case "ICmdLogic.amount":
		i, err := strconv.Atoi(text)
		if err != nil {
			return err
		}
		b.Amount = int16(i)
	default:
	}

	return nil
}

func (b *NewLogistic) GetInputTemplate(requireRawParamAttr string) interface{} {
	return nil
}

func (b *NewLogistic) Do(text string) (resultErr error) {
	if u, err := lineuser.Get(b.Context.GetUserID()); err != nil {
		return err
	} else {
		if u.Role != domain.ADMIN_CLUB_ROLE {
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
		if resultErr = b.InsertLogistic(transaction); resultErr != nil {
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

	boxComponent := linebot.GetFlexMessageBoxComponent(
		linebotDomain.VERTICAL_MESSAGE_LAYOUT,
		&linebotModel.FlexMessageBoxComponentOption{
			Margin:  linebotDomain.LG_FLEX_MESSAGE_SIZE,
			Spacing: linebotDomain.SM_FLEX_MESSAGE_SIZE,
		},
	)

	if js, err := b.Context.
		GetDateTimeCmdInputMode(domain.DATE_POSTBACK_DATE_TIME_CMD, "date").
		GetSignal(); err != nil {
		return err
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

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.name", "品項", false).
		GetSignal(); err != nil {
		return err
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

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.amount", "數量", false).
		GetSignal(); err != nil {
		return err
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

	if js, err := b.Context.
		GetRequireInputMode("ICmdLogic.description", "備註", false).
		GetSignal(); err != nil {
		return err
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

	if js, err := b.Context.
		GetComfirmMode().
		GetSignal(); err != nil {
		return err
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
		return err
	}

	return nil
}

func (b *NewLogistic) InsertLogistic(transaction *gorm.DB) (resultErr error) {
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

	data := &logisticDb.LogisticTable{
		Date:        b.Date,
		Name:        b.Name,
		Amount:      b.Amount,
		Description: b.Description,
	}
	if resultErr = database.Club.Logistic.BaseTable.Insert(transaction, data); resultErr != nil {
		return
	}

	return nil
}
