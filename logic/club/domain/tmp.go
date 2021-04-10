package domain

import linebotModel "heroku-line-bot/service/linebot/domain/model"

type NewActivityLineTemplate struct {
	DateAction,
	PlaceAction,
	ClubSubsidyAction,
	PeopleLimitAction interface{}
	CourtAction *linebotModel.PostBackAction
}
