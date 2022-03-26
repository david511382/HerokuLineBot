package club

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/src/logic/club/domain"
	clublinebotDomain "heroku-line-bot/src/logic/clublinebot/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	"heroku-line-bot/src/pkg/util"
	commonRedis "heroku-line-bot/src/repo/redis/common"
	"io/ioutil"
	"path/filepath"
	"strings"

	errUtil "heroku-line-bot/src/pkg/util/error"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	adminRichMenuImg,
	cadreRichMenuImg,
	guestRichMenuImg []byte

	liffUrl    string
	clubTeamID int
)

func Init(cfg *bootstrap.Config) errUtil.IError {
	root, err := bootstrap.GetRootDirPath()
	if err != nil {
		return errUtil.NewError(err)
	}

	if errInfo := readImg(root); errInfo != nil {
		return errInfo
	}

	clubTeamID = cfg.Badminton.ClubTeamID
	liffUrl = cfg.Badminton.LiffUrl
	return nil
}

func readImg(rootPath string) errUtil.IError {
	var err error
	rootPath = filepath.Join(rootPath, "resource", "img")

	{
		fileName := filepath.Join(rootPath, "adminRichMenu.png")
		adminRichMenuImg, err = ioutil.ReadFile(fileName)
		if err != nil {
			return errUtil.NewError(err)
		}
	}

	{
		fileName := filepath.Join(rootPath, "cadreRichMenu.png")
		cadreRichMenuImg, err = ioutil.ReadFile(fileName)
		if err != nil {
			return errUtil.NewError(err)
		}
	}

	{
		fileName := filepath.Join(rootPath, "guestRichMenu.png")
		guestRichMenuImg, err = ioutil.ReadFile(fileName)
		if err != nil {
			return errUtil.NewError(err)
		}
	}

	return nil
}

func HandlerTextCmd(text string, lineContext clublinebotDomain.IContext) (resultErrInfo errUtil.IError) {
	cmd := domain.TextCmd(text)
	var cmdHandler domain.ICmdHandler
	paramJson := ""
	textJsonResult := gjson.Parse(text)
	isInputTextJson := textJsonResult.Type == gjson.JSON
	if handler, err := getCmdHandler(cmd, lineContext); err != nil {
		resultErrInfo = errUtil.NewError(err)
		return
	} else if isChangeHandler := handler != nil; isChangeHandler {
		cmdHandler = handler
		text = ""
		if err := lineContext.DeleteParam(); commonRedis.IsRedisError(err) {
			resultErrInfo = errUtil.NewError(err)
		}
		if errInfo := cmdHandler.CacheParams(); errInfo != nil {
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
	} else {
		cmd = getCmdFromJson(textJsonResult)
		if cmd == "" {
			if isInputTextJson {
				if pathConverter := getCmdBasePathConverterFromJson(textJsonResult); pathConverter != nil {
					paramJson = "{}"
					for path, jr := range textJsonResult.Map() {
						path := pathConverter(path)
						s, err := sjson.Set(paramJson, path, jr.Value())
						if err != nil {
							errInfo := errUtil.NewError(err)
							resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
							return
						}
						paramJson = s
					}
				} else {
					paramJson = text
				}
			}
		} else {
			if handler, err := getCmdHandler(cmd, lineContext); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			} else if handler != nil {
				cmdHandler = handler
			}
			paramJson = text
		}
	}

	if redisParamJson := lineContext.GetParam(); redisParamJson != "" {
		textJsonResult := gjson.Parse(redisParamJson)
		redisCmd := getCmdFromJson(textJsonResult)
		isChangeCmd := cmdHandler != nil && redisCmd != cmd
		if isChangeCmd {
			// if err := lineContext.DeleteParam(); commonRedis.IsRedisError(err) {
			// 	resultErrInfo = errUtil.NewError(err)
			// }
			// cmdHandler = nil
		}

		if cmdHandler == nil || !isInputTextJson {
			if handler, err := getCmdHandler(redisCmd, lineContext); err != nil {
				resultErrInfo = errUtil.NewError(err)
				return
			} else {
				cmdHandler = handler
			}
		}

		if errInfo := cmdHandler.ReadParam(textJsonResult); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
	}

	if cmdHandler == nil {
		replyMessges := []interface{}{
			linebot.GetTextMessage("聽不懂你在說什麼"),
		}
		if err := lineContext.Reply(replyMessges); err != nil {
			resultErrInfo = errUtil.NewError(err)
			return
		}
		return nil
	}

	if paramJson != "" {
		jr := gjson.Parse(paramJson)
		if errInfo := cmdHandler.ReadParam(jr); errInfo != nil {
			resultErrInfo = errInfo
			return
		}
	}

	if isInputTextJson {
		text = ""
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
		logicHandler = &GetConfirmRegisters{}
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
		IContext:   context,
		ICmdLogic:  logicHandler,
		InputParam: *NewInputParam(logicHandler),
	}
	if errInfo := logicHandler.Init(result); errInfo != nil {
		return nil, errInfo
	}

	return result, nil
}

func getCmdFromJson(jr gjson.Result) domain.TextCmd {
	cmdJs := jr.Get(domain.ATTR_CMD)
	return domain.TextCmd(cmdJs.String())
}

func getDateTimeCmdFromJson(jr gjson.Result) domain.DateTimeCmd {
	cmdJs := jr.Get(domain.ATTR_DATE_TIME_CMD)
	return domain.DateTimeCmd(cmdJs.String())
}

func getCmdBasePathConverterFromJson(jr gjson.Result) (cmdBasePathConverter func(attr string) string) {
	basePath := jr.Get(domain.ATTR_CMD_BASE_PATH).Str
	if basePath == "" {
		return nil
	}
	return func(attr string) string {
		if attr == domain.ATTR_CMD_BASE_PATH {
			return attr
		}
		return strings.Join(
			[]string{
				basePath,
				attr,
			},
			".",
		)
	}
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
	p := util.NewFloat(float64(people * domain.MONEY_UNIT))
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
