package domain

import linebotModel "heroku-line-bot/service/linebot/domain/model"

var (
	NormalButtonOption = linebotModel.ButtonOption{
		Color: "#00dd00",
	}
	AlertButtonOption = linebotModel.ButtonOption{
		Color: "#dd00dd",
	}
)
