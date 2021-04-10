package linebot

import (
	"heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model/reqs"
	"heroku-line-bot/service/linebot/domain/model/resp"
	"heroku-line-bot/util"
	"net/http"
)

type LineBot struct {
	channelAccessToken string
}

func (lb *LineBot) setRequestAuthorization(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+lb.channelAccessToken)
}

func (lb *LineBot) getGetRequest(uri string, param interface{}) (*http.Request, error) {
	request, err := util.GetRequest(uri, param)
	if err != nil {
		return nil, err
	}
	lb.setRequestAuthorization(request)
	return request, nil
}

func (lb *LineBot) getPostRequest(uri string, param interface{}) (*http.Request, error) {
	request, err := util.JsonRequest(uri, util.POST, param)
	if err != nil {
		return nil, err
	}
	lb.setRequestAuthorization(request)
	return request, nil
}

func (lb *LineBot) getDeleteRequest(uri string, param interface{}) (*http.Request, error) {
	request, err := util.GetRequest(uri, param)
	if err != nil {
		return nil, err
	}
	lb.setRequestAuthorization(request)
	return request, nil
}

func (lb *LineBot) GetUserProfile(userID string) (*resp.GetUserProfile, error) {
	url := domain.LINE_URL + "/profile/" + userID
	request, err := lb.getGetRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.GetUserProfile{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) ReplyMessage(param *reqs.ReplyMessage) (*resp.ReplyMessage, error) {
	url := domain.LINE_URL + "/message/reply"
	request, err := lb.getPostRequest(url, param)
	if err != nil {
		return nil, err
	}

	response := &resp.ReplyMessage{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) PushMessage(param *reqs.PushMessage) (*resp.PushMessage, error) {
	url := domain.LINE_URL + "/message/push"
	request, err := lb.getPostRequest(url, param)
	if err != nil {
		return nil, err
	}

	response := &resp.PushMessage{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}
