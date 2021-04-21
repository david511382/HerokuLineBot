package clublinebot

import (
	"encoding/json"
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
			groupID := event.Source.GroupID
			pushMessages := []interface{}{
				linebot.GetTextMessage("member join group : " + groupID),
			}

			userInfoMsgs := []string{}
			for _, source := range event.Joined.Members {
				userID := source.UserID
				userInfo, err := b.GetUserProfile(userID)
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("%s : %s", userInfo.DisplayName, userID)
				userInfoMsgs = append(userInfoMsgs, msg)
			}
			if len(userInfoMsgs) > 0 {
				userInfoMsg := strings.Join(userInfoMsgs, "\n")
				pushMessages = append(pushMessages, linebot.GetTextMessage(userInfoMsg))
			} else {
				bs, err := json.Marshal(event)
				if err != nil {
					return err
				}
				msg := fmt.Sprintf("event:%s", string(bs))
				pushMessages = append(pushMessages, linebot.GetTextMessage(msg))
			}

			pushReqs := &lineBotReqs.PushMessage{
				To:       b.lineAdminID,
				Messages: pushMessages,
			}
			if _, err := b.PushMessage(pushReqs); err != nil {
				return err
			}

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

			return nil
		},
		replyToken,
	); err != nil {
		return err
	}

	return nil
}
