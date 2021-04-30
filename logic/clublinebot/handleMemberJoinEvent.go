package clublinebot

import (
	"fmt"
	"heroku-line-bot/service/linebot"
	lineBotModel "heroku-line-bot/service/linebot/domain/model"
	lineBotReqs "heroku-line-bot/service/linebot/domain/model/reqs"
	"strings"
)

func (b *ClubLineBot) handleMemberJoinedEvent(event *lineBotModel.MemberJoinEvent) error {
	replyToken := event.ReplyToken
	if err := b.tryLine(
		func() error {
			replyMessges := []interface{}{
				linebot.GetTextMessage("歡迎加入，跟我加入好友可以獲取更多社團的資訊喔!"),
			}
			replyReqs := &lineBotReqs.ReplyMessage{
				ReplyToken: replyToken,
				Messages:   replyMessges,
			}
			if _, err := b.ReplyMessage(replyReqs); err != nil {
				return err
			}

			groupID := event.Source.GroupID
			adminMessages := []string{
				"member join group : " + groupID,
			}
			userInfoMsgs := []string{}
			for _, source := range event.Joined.Members {
				userID := source.UserID
				userInfo, err := b.GetUserProfile(userID)
				if err != nil {
					adminMessages = append(adminMessages, userID)
					adminMessages = append(adminMessages, err.Error())
					continue
				}
				msg := fmt.Sprintf("%s : %s", userInfo.DisplayName, userID)
				userInfoMsgs = append(userInfoMsgs, msg)
			}
			if len(userInfoMsgs) > 0 {
				userInfoMsg := strings.Join(userInfoMsgs, "\n")
				adminMessages = append(adminMessages, userInfoMsg)
			}

			if err := b.pushMessageToAdmin(adminMessages...); err != nil {
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
