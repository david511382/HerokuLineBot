package club

import (
	"embed"
	"fmt"
	"heroku-line-bot/logic/club/domain"
	clublinebotDomain "heroku-line-bot/logic/clublinebot/domain"
	commonLogic "heroku-line-bot/logic/common"
	"heroku-line-bot/service/linebot"
	"heroku-line-bot/storage/redis"
	"heroku-line-bot/util"
	"math"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	adminRichMenuImg,
	cadreRichMenuImg,
	guestRichMenuImg []byte
)

func Init(f embed.FS) error {
	if bs, err := readImg(f, "adminRichMenu.png"); err != nil {
		return err
	} else {
		adminRichMenuImg = bs
	}
	if bs, err := readImg(f, "cadreRichMenu.png"); err != nil {
		return err
	} else {
		cadreRichMenuImg = bs
	}
	if bs, err := readImg(f, "guestRichMenu.png"); err != nil {
		return err
	} else {
		guestRichMenuImg = bs
	}
	return nil
}

func readImg(f embed.FS, fileName string) ([]byte, error) {
	fileName = fmt.Sprintf("resource/img/%s", fileName)
	return f.ReadFile(fileName)
}

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
		if cmd == "" {
			if !isSingelParamText {
				paramJson = text
			}
		} else {
			if handler, err := getCmdHandler(cmd, lineContext); err != nil {
				return err
			} else if handler != nil {
				cmdHandler = handler
			}
			paramJson = text
		}

		dateTimeCmd := getDateTimeCmdFromJson(text)
		switch dateTimeCmd {
		case domain.TIME_POSTBACK_DATE_TIME_CMD, domain.DATE_TIME_POSTBACK_DATE_TIME_CMD, domain.DATE_POSTBACK_DATE_TIME_CMD:
			if js, err := sjson.Delete(text, domain.DATE_TIME_CMD_ATTR); err != nil {
				return err
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
	if err := logicHandler.Init(result); err != nil {
		return nil, err
	}

	return result, nil
}

func getCmd(cmd domain.TextCmd, pathValueMap map[string]interface{}) (string, error) {
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

func calculateActivity(ballConsume, courtFee float64) (activityFee, ballFee float64) {
	ballFee = ballConsume * domain.PRICE_PER_BALL
	return commonLogic.FloatPlus(ballFee, courtFee), ballFee
}

func calculateActivityPay(people int, ballConsume, courtFee, clubSubsidy float64) (activityFee float64, clubMemberFee, guestFee int) {
	activityFee, _ = calculateActivity(ballConsume, courtFee)
	clubMemberFee, guestFee = calculatePay(people, activityFee, clubSubsidy)
	return
}

func calculatePay(people int, activityFee, clubSubsidy float64) (clubMemberFee, guestFee int) {
	shareMoney := commonLogic.FloatMinus(activityFee, clubSubsidy)

	p := people * domain.MONEY_UNIT
	clubMemberFee = int(math.Ceil(shareMoney/float64(p)) * domain.MONEY_UNIT)
	guestFee = int(math.Ceil(activityFee/float64(p)) * domain.MONEY_UNIT)
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
