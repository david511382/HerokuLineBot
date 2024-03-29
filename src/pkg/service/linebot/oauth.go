package linebot

import (
	"heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/service/linebot/domain/model/reqs"
	"heroku-line-bot/src/pkg/service/linebot/domain/model/resp"
	"heroku-line-bot/src/pkg/util"
	"strconv"
)

type OAuth struct {
	channelID uint64
}

func NewOAuth(channelID uint64) *OAuth {
	return &OAuth{
		channelID: channelID,
	}
}

func (oa *OAuth) VerifyIDToken(param reqs.OAuthVerifyIDToken) (*resp.OAuthVerifyIDToken, error) {
	url := domain.LINE_OAUTH_URL + "/verify"
	param.ClientID = strconv.FormatUint(oa.channelID, 10)
	request, err := util.FormRequest(url, util.POST, param)
	if err != nil {
		return nil, err
	}

	response := &resp.OAuthVerifyIDToken{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}
