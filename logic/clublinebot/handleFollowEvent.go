package clublinebot

import (
	clubLogicDomain "heroku-line-bot/logic/club/domain"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/service/linebot/domain"
	lineBotModel "heroku-line-bot/service/linebot/domain/model"
	lineBotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
)

func (b *ClubLineBot) handleFollowEvent(event *lineBotModel.FollowEvent) error {
	replyToken := event.ReplyToken
	if err := b.tryLine(
		func() error {
			replyMessges := []interface{}{
				linebot.GetFlexMessage(
					"歡迎!",
					linebot.GetFlexMessageBubbleContent(
						linebot.GetFlexMessageBoxComponent(
							domain.VERTICAL_MESSAGE_LAYOUT,
							nil,
							linebot.GetFlexMessageTextComponent(0, "麻煩註冊，註冊後就能使用服務喔"),
							linebot.GetButtonComponent(0,
								linebot.GetMessageAction(string(clubLogicDomain.REGISTER_TEXT_CMD)),
								nil,
							),
						),
						nil,
					),
				),
			}
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
