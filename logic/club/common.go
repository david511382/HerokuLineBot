package club

import (
	"heroku-line-bot/logic/club/domain"
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/redis"
	"heroku-line-bot/util"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func HandlerTextCmd(text string, lineContext clublinebotDomain.IContext) (resultErr error) {
	cmd := domain.TextCmd(text)
	var cmdHandler domain.ICmdHandler
	paramJson := ""
	isSingelParamText := !util.IsJSON(text)

	if handler, err := getCmdHandler(cmd, lineContext); err != nil {
		return err
	} else if handler != nil {
		cmdHandler = handler
		if err := lineContext.DeleteParam(); redis.IsRedisError(err) {
			resultErr = err
		}
	} else {
		cmd = getCmdFromJson(text)
		switch cmd {
		case domain.TIME_POSTBACK_CMD, domain.DATE_TIME_POSTBACK_CMD, domain.DATE_POSTBACK_CMD:
			if js, err := sjson.Delete(text, domain.CMD_ATTR); err != nil {
				return err
			} else {
				paramJson = js
			}
			jr := gjson.Get(text, string(cmd))
			text = jr.String()
			cmd = ""
			isSingelParamText = true
		default:
			if !isSingelParamText {
				paramJson = text
			}
		}
	}

	if redisParamJson := lineContext.GetParam(); redisParamJson != "" {
		redisCmd := getCmdFromJson(redisParamJson)
		if isChangeCmd := cmdHandler != nil && redisCmd != cmd; isChangeCmd {
			if err := lineContext.DeleteParam(); redis.IsRedisError(err) {
				resultErr = err
			}
			cmdHandler = nil
		}

		if cmdHandler == nil || isSingelParamText {
			if handler, err := getCmdHandler(redisCmd, lineContext); err != nil {
				return err
			} else {
				cmdHandler = handler
			}
		}

		if isSingelParamText {
			cmdHandler.SetSingleParamMode()
		}

		if err := cmdHandler.ReadParam([]byte(redisParamJson)); err != nil {
			return err
		}
	}

	if cmdHandler == nil {
		replyMessges := []interface{}{
			linebot.GetTextMessage("聽不懂你在說什麼"),
		}
		if err := lineContext.Reply(replyMessges); err != nil {
			return err
		}
		return nil
	}

	if paramJson != "" {
		if err := cmdHandler.ReadParam([]byte(paramJson)); err != nil {
			return err
		}
	}

	if err := cmdHandler.Do(text); err != nil {
		return err
	}

	return resultErr
}

func getCmdHandler(cmd domain.TextCmd, context clublinebotDomain.IContext) (domain.ICmdHandler, error) {
	var logicHandler domain.ICmdLogic
	switch cmd {
	case domain.NEW_ACTIVITY_TEXT_CMD:
		logicHandler = &newActivity{}
	case domain.GET_ACTIVITIES_TEXT_CMD:
		logicHandler = &getActivities{}
	case domain.REGISTER_TEXT_CMD:
		logicHandler = &register{}
	default:
		return nil, nil
	}

	result := &CmdHandler{
		CmdBase: &domain.CmdBase{
			Cmd: cmd,
		},
		IContext:  context,
		ICmdLogic: logicHandler,
	}
	if err := logicHandler.Init(
		result,
		func(requireRawParamAttr, requireRawParamAttrText string, isInputImmediately bool) {
			result.RequireRawParamAttr = requireRawParamAttr
			result.RequireRawParamAttrText = requireRawParamAttrText
			result.IsInputImmediately = isInputImmediately
		},
	); err != nil {
		return nil, err
	}

	return result, nil
}

func getCmdFromJson(json string) domain.TextCmd {
	cmdJs := gjson.Get(json, domain.CMD_ATTR)
	return domain.TextCmd(cmdJs.String())
}
