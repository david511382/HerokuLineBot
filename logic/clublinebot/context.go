package clublinebot

import (
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/redis"
)

type Context struct {
	userID     string
	replyToken string
	lineBot    *ClubLineBot
}

func newContext(
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
	if err := redis.UserUsingStatus.HSet(c.userID, json); redis.IsRedisError(err) {
		return err
	}
	return nil
}

func (c *Context) DeleteParam() error {
	if _, err := redis.UserUsingStatus.HDel(c.userID); redis.IsRedisError(err) {
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

func (c *Context) GetBot() *linebot.LineBot {
	return c.lineBot.LineBot
}
