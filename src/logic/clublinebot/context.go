package clublinebot

import (
	"heroku-line-bot/src/pkg/service/linebot"
	"heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	"heroku-line-bot/src/repo/redis"
)

type Context struct {
	userID     string
	replyToken string
	lineBot    *ClubLineBot
}

func NewContext(
	userID string,
	replyToken string,
	lineBot *ClubLineBot,
) Context {
	return Context{
		userID:     userID,
		replyToken: replyToken,
		lineBot:    lineBot,
	}
}

func (c *Context) GetUserID() string {
	return c.userID
}

func (c *Context) GetUserName() string {
	if profile, err := c.lineBot.GetUserProfile(c.userID); err != nil {
		return ""
	} else {
		return profile.DisplayName
	}
}

func (c *Context) SaveParam(json string) error {
	if errInfo := redis.Badminton().UserUsingStatus.HSet(c.userID, json); errInfo != nil && errInfo.IsError() {
		return errInfo
	}
	return nil
}

func (c *Context) DeleteParam() error {
	if _, errInfo := redis.Badminton().UserUsingStatus.HDel(c.userID); errInfo != nil && errInfo.IsError() {
		return errInfo
	}
	return nil
}

func (c *Context) GetParam() (json *string) {
	json, _ = redis.Badminton().UserUsingStatus.HGet(c.userID)
	return
}

func (c *Context) Reply(replyMessges []interface{}) error {
	return c.lineBot.tryReply(c.replyToken, replyMessges)
}

func (c *Context) PushAdmin(replyMessges []interface{}) error {
	return c.lineBot.tryLine(
		func() error {
			_, err := c.lineBot.PushMessage(&reqs.PushMessage{
				To:       c.lineBot.lineAdminID,
				Messages: replyMessges,
			})
			return err
		},
		c.replyToken,
	)
}

func (c *Context) PushRoom(roomID string, replyMessges []interface{}) error {
	return c.lineBot.tryLine(
		func() error {
			_, err := c.lineBot.PushMessage(&reqs.PushMessage{
				To:       roomID,
				Messages: replyMessges,
			})
			return err
		},
		c.replyToken,
	)
}

func (c *Context) GetBot() *linebot.LineBot {
	return c.lineBot.LineBot
}
