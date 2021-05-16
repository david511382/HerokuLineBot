package domain

import (
	linebotDomain "heroku-line-bot/service/linebot/domain"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
)

var (
	NormalButtonOption = linebotModel.ButtonOption{
		Color:      "#00dd00",
		Style:      "primary",
		Height:     linebotDomain.SM_FLEX_MESSAGE_SIZE,
		AdjustMode: linebotDomain.SHRINK_TO_FIT_ADJUST_MODE,
	}
	AlertButtonOption = linebotModel.ButtonOption{
		Color:      "#dd00dd",
		Style:      "primary",
		Height:     linebotDomain.SM_FLEX_MESSAGE_SIZE,
		AdjustMode: linebotDomain.SHRINK_TO_FIT_ADJUST_MODE,
	}
	DarkButtonOption = linebotModel.ButtonOption{
		Color:      "#888888",
		Style:      "primary",
		Height:     linebotDomain.SM_FLEX_MESSAGE_SIZE,
		AdjustMode: linebotDomain.SHRINK_TO_FIT_ADJUST_MODE,
	}
)

const (
	RED_COLOR        = "#FF6347"
	WHITE_COLOR      = "#ffffff"
	BLUE_GREEN_COLOR = "#00cc99"
)
