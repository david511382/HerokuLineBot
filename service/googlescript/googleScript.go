package googlescript

import (
	"heroku-line-bot/service/googlescript/domain/model/reqs"
	"heroku-line-bot/util"
)

type GoogleScript struct {
	url string
}

func (gs *GoogleScript) LinePage(json string) error {
	queryStringParam := &reqs.Page{
		Page: "line",
	}
	url := util.QueryString(gs.url, queryStringParam)
	body := []byte(json)
	request, err := util.RawJsonRequest(url, util.POST, body)
	if err != nil {
		return err
	}

	if _, err := util.SendRequest(request, nil); err != nil {
		return err
	}

	return nil
}
