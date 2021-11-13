package clublinebot

import (
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/storage/redis"
	commonRedis "heroku-line-bot/storage/redis/common"
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
	if err := redis.UserUsingStatus.HSet(c.userID, json); commonRedis.IsRedisError(err) {
		return err
	}
	return nil
}

func (c *Context) DeleteParam() error {
	if _, err := redis.UserUsingStatus.HDel(c.userID); commonRedis.IsRedisError(err) {
		return err
	}
	return nil
}

func (c *Context) GetParam() (json string) {
	json, _ = redis.UserUsingStatus.HGet(c.userID)
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

func (c *Context) PushRoom(replyMessges []interface{}) error {
	return c.lineBot.tryLine(
		func() error {
			_, err := c.lineBot.PushMessage(&reqs.PushMessage{
				To:       c.lineBot.lineRoomID,
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
