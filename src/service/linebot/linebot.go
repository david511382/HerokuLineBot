package linebot

import (
	"heroku-line-bot/src/service/linebot/domain"
	"heroku-line-bot/src/service/linebot/domain/model/reqs"
	"heroku-line-bot/src/service/linebot/domain/model/resp"
	"heroku-line-bot/src/util"
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

func (lb *LineBot) getRawPostRequest(uri string, body []byte) (*http.Request, error) {
	request, err := util.RawJsonRequest(uri, util.POST, body)
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

// Limit 5 Messages
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

func (lb *LineBot) ListRichMenu() (*resp.ListRichMenu, error) {
	url := domain.LINE_URL + "/richmenu/list"
	request, err := lb.getPostRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.ListRichMenu{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) DeleteRichMenu(richMenuID string) (*resp.DeleteRichMenu, error) {
	url := domain.LINE_URL + "/richmenu/" + richMenuID
	request, err := lb.getDeleteRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.DeleteRichMenu{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) SetDefaultRichMenu(richMenuID string) (*resp.SetDefaultRichMenu, error) {
	url := domain.LINE_URL + "/user/all/richmenu/" + richMenuID
	request, err := lb.getPostRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.SetDefaultRichMenu{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) SetRichMenuTo(richMenuID, userID string) (*resp.SetRichMenuTo, error) {
	url := domain.LINE_URL + "/user/" + userID + "/richmenu/" + richMenuID
	request, err := lb.getPostRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.SetRichMenuTo{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) SetRichMenuTos(richMenuID string, userIDs ...string) (*resp.SetRichMenuTos, error) {
	url := domain.LINE_URL + "/user/all/richmenu/" + richMenuID
	requestArg := &reqs.SetRichMenuTos{
		RichMenuID: richMenuID,
		UserID:     userIDs,
	}
	request, err := lb.getPostRequest(url, requestArg)
	if err != nil {
		return nil, err
	}

	response := &resp.SetRichMenuTos{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) GetDefaultRichMenu() (*resp.GetDefaultRichMenu, error) {
	url := domain.LINE_URL + "/user/all/richmenu"
	request, err := lb.getGetRequest(url, nil)
	if err != nil {
		return nil, err
	}

	response := &resp.GetDefaultRichMenu{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) UploadRichMenuImage(richMenuID string, image []byte) (*resp.UploadRichMenuImage, error) {
	url := domain.LINE_DATA_URL + "/richmenu/" + richMenuID + "/content"
	request, err := lb.getRawPostRequest(url, image)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "image/png")

	response := &resp.UploadRichMenuImage{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (lb *LineBot) CreateRichMenu(param *reqs.CreateRichMenu) (*resp.CreateRichMenu, error) {
	url := domain.LINE_URL + "/richmenu"
	request, err := lb.getPostRequest(url, param)
	if err != nil {
		return nil, err
	}

	response := &resp.CreateRichMenu{}
	if _, err := util.SendRequest(request, response); err != nil {
		return nil, err
	}
	return response, nil
}
