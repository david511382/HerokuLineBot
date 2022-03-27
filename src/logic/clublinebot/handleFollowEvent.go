package clublinebot

import (
	"fmt"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/logic/club"
	clubLogicDomain "heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	"heroku-line-bot/src/pkg/service/linebot/domain"
	lineBotModel "heroku-line-bot/src/pkg/service/linebot/domain/model"
	lineBotReqs "heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	"heroku-line-bot/src/pkg/util"
)

func (b *ClubLineBot) handleFollowEvent(event *lineBotModel.FollowEvent) error {
	replyToken := event.ReplyToken
	if err := b.tryLine(
		func() error {
			userLineID := event.Source.UserID
			c := NewContext(userLineID, replyToken, b)
			name := c.GetUserName()
			registerMember := club.NewRegisterMember(name, util.GetStringP(userLineID))
			registerMember.Init(&c)

			replyMessges := make([]interface{}, 0)

			isExist := false
			if _, dbName, errInfo := registerMember.LoadMemberID(); errInfo != nil {
				if errInfo.IsError() {
					return errInfo
				}

				logger.Log("LINE_BOT", errInfo)
			} else if isExist = dbName != nil; isExist {
				replyMessges = append(replyMessges,
					linebot.GetTextMessage(fmt.Sprintf("您已註冊過，暱稱使用以前您使用的暱稱: %s，可在選單中 修改會員資料 修改", *dbName)),
				)
			} else {
				replyMessges = append(replyMessges,
					linebot.GetTextMessage(fmt.Sprintf("暱稱預設使用您的 Line 名稱: %s，可在選單中 修改會員資料 修改", name)),
				)
			}

			if errInfo := registerMember.Registe(nil); errInfo != nil {
				if errInfo.IsError() {
					return errInfo
				}

				logger.Log("LINE_BOT", errInfo)
			}

			if !isExist {
				if errInfo := registerMember.NotifyAdmin(); errInfo != nil {
					logger.Log("LINE_BOT", errInfo)
				}
			}

			replyMessges = append(replyMessges, linebot.GetFlexMessage(
				"歡迎!",
				linebot.GetFlexMessageBubbleContent(
					linebot.GetFlexMessageBoxComponent(
						domain.VERTICAL_MESSAGE_LAYOUT,
						nil,
						linebot.GetFlexMessageTextComponent("若是公司成員，麻煩登記資料", nil),
						linebot.GetClassButtonComponent(
							linebot.GetMessageAction(string(clubLogicDomain.REGISTE_COMPANY_TEXT_CMD)),
						),
					),
					nil,
				),
			))

			replyReqs := &lineBotReqs.ReplyMessage{
				ReplyToken: replyToken,
				Messages:   replyMessges,
			}
			if _, err := b.ReplyMessage(replyReqs); err != nil {
				return err
			}

			return nil
		},
		replyToken,
	); err != nil {
		return err
	}

	return nil
}
