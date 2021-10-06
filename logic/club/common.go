package club

import (
	"embed"
	"fmt"
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/logic/club/domain"
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/redis"
	"heroku-line-bot/util"

	errLogic "heroku-line-bot/logic/error"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	adminRichMenuImg,
	cadreRichMenuImg,
	guestRichMenuImg []byte

	liffUrl string
)

func Init(cfg *bootstrap.Config, resourceFS embed.FS) errLogic.IError {
	if bs, err := readImg(resourceFS, "adminRichMenu.png"); err != nil {
		return errLogic.NewError(err)
	} else {
		adminRichMenuImg = bs
	}
	if bs, err := readImg(resourceFS, "cadreRichMenu.png"); err != nil {
		return errLogic.NewError(err)
	} else {
		cadreRichMenuImg = bs
	}
	if bs, err := readImg(resourceFS, "guestRichMenu.png"); err != nil {
		return errLogic.NewError(err)
	} else {
		guestRichMenuImg = bs
	}

	liffUrl = cfg.Badminton.LiffUrl
	return nil
}

func readImg(resourceFS embed.FS, fileName string) ([]byte, error) {
	fileName = fmt.Sprintf("resource/img/%s", fileName)
	return resourceFS.ReadFile(fileName)
}

func HandlerTextCmd(text string, lineContext clublinebotDomain.IContext) (resultErrInfo errLogic.IError) {
	cmd := domain.TextCmd(text)
	var cmdHandler domain.ICmdHandler
	paramJson := ""
	isSingelParamText := !util.IsJSON(text)
	if handler, err := getCmdHandler(cmd, lineContext); err != nil {
		resultErrInfo = errLogic.NewError(err)
		return
	} else if handler != nil {
		cmdHandler = handler
		if err := lineContext.DeleteParam(); redis.IsRedisError(err) {
			resultErrInfo = errLogic.NewError(err)
		}
	} else {
		cmd = getCmdFromJson(text)
		if cmd == "" {
			if !isSingelParamText {
				paramJson = text
			}
		} else {
			if handler, err := getCmdHandler(cmd, lineContext); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			} else if handler != nil {
				cmdHandler = handler
			}
			paramJson = text
		}

		dateTimeCmd := getDateTimeCmdFromJson(text)
		switch dateTimeCmd {
		case domain.TIME_POSTBACK_DATE_TIME_CMD, domain.DATE_TIME_POSTBACK_DATE_TIME_CMD, domain.DATE_POSTBACK_DATE_TIME_CMD:
			if js, err := sjson.Delete(text, domain.DATE_TIME_CMD_ATTR); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			} else {
				paramJson = js
			}
			jr := gjson.Get(text, string(dateTimeCmd))
			text = jr.String()
			isSingelParamText = true
		}
	}

	if redisParamJson := lineContext.GetParam(); redisParamJson != "" {
		redisCmd := getCmdFromJson(redisParamJson)
		if isChangeCmd := cmdHandler != nil && redisCmd != cmd; isChangeCmd {
			if err := lineContext.DeleteParam(); redis.IsRedisError(err) {
				resultErrInfo = errLogic.NewError(err)
			}
			cmdHandler = nil
		}

		if cmdHandler == nil || isSingelParamText {
			if handler, err := getCmdHandler(redisCmd, lineContext); err != nil {
				resultErrInfo = errLogic.NewError(err)
				return
			} else {
				cmdHandler = handler
			}
		}

		if isSingelParamText {
			cmdHandler.SetSingleParamMode()
		}

		if errInfo := cmdHandler.ReadParam([]byte(redisParamJson)); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
	}

	if cmdHandler == nil {
		replyMessges := []interface{}{
			linebot.GetTextMessage("聽不懂你在說什麼"),
		}
		if err := lineContext.Reply(replyMessges); err != nil {
			resultErrInfo = errLogic.NewError(err)
			return
		}
		return nil
	}

	if paramJson != "" {
		if errInfo := cmdHandler.ReadParam([]byte(paramJson)); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
	}

	if errInfo := cmdHandler.Do(text); errInfo != nil {
		resultErrInfo = errInfo
		return
	}

	return
}

func getCmdHandler(cmd domain.TextCmd, context clublinebotDomain.IContext) (domain.ICmdHandler, error) {
	var logicHandler domain.ICmdLogic
	switch cmd {
	case domain.NEW_ACTIVITY_TEXT_CMD:
		logicHandler = &NewActivity{}
	case domain.GET_ACTIVITIES_TEXT_CMD:
		logicHandler = &GetActivities{}
	case domain.REGISTER_TEXT_CMD:
		logicHandler = &register{}
	case domain.CONFIRM_REGISTER_TEXT_CMD:
		logicHandler = &confirmRegister{}
	case domain.GET_CONFIRM_REGISTER_TEXT_CMD:
		logicHandler = &GetComfirmRegisters{}
	case domain.SUBMIT_ACTIVITY_TEXT_CMD:
		logicHandler = &submitActivity{}
	case domain.RICH_MENU_TEXT_CMD:
		logicHandler = &richMenu{}
	case domain.NEW_LOGISTIC_TEXT_CMD:
		logicHandler = &NewLogistic{}
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
	if errInfo := logicHandler.Init(result); errInfo != nil {
		return nil, errInfo
	}

	return result, nil
}

func getCmd(cmd domain.TextCmd, pathValueMap map[string]interface{}) (string, errLogic.IError) {
	cmdHandler := &CmdHandler{
		CmdBase: &domain.CmdBase{
			Cmd: cmd,
		},
	}
	return cmdHandler.
		GetCmdInputMode(nil).
		GetKeyValueInputMode(pathValueMap).
		GetSignal()
}

func getCmdFromJson(json string) domain.TextCmd {
	cmdJs := gjson.Get(json, domain.CMD_ATTR)
	return domain.TextCmd(cmdJs.String())
}

func getDateTimeCmdFromJson(json string) domain.DateTimeCmd {
	cmdJs := gjson.Get(json, domain.DATE_TIME_CMD_ATTR)
	return domain.DateTimeCmd(cmdJs.String())
}

func calculateActivity(ballConsume, courtFee util.Float) (activityFee, ballFee util.Float) {
	ballFee = ballConsume.MulFloat(float64(domain.PRICE_PER_BALL))
	return ballFee.Plus(courtFee), ballFee
}

func calculateActivityPay(people int, ballConsume, courtFee, clubSubsidy util.Float) (activityFee util.Float, clubMemberFee, guestFee int) {
	activityFee, _ = calculateActivity(ballConsume, courtFee)
	clubMemberFee, guestFee = calculatePay(people, activityFee, clubSubsidy)
	return
}

func calculatePay(people int, activityFee, clubSubsidy util.Float) (clubMemberFee, guestFee int) {
	shareMoney := activityFee.Minus(clubSubsidy)
	p := util.ToFloat(float64(people * domain.MONEY_UNIT))
	pp := int(shareMoney.Div(p).Ceil().ToInt())
	clubMemberFee = pp * domain.MONEY_UNIT
	guestFee = int(activityFee.Div(p).Ceil().ToInt()) * domain.MONEY_UNIT
	return
}

func getJoinCount(totalCount int, limit *int16) (joinedCount, waitingCount int) {
	joinedCount = totalCount
	peopleLimit := 0
	if limit != nil {
		peopleLimit = int(*limit)
		if joinedCount > peopleLimit {
			waitingCount = joinedCount - peopleLimit
			joinedCount = peopleLimit
		}
	}
	return
}
