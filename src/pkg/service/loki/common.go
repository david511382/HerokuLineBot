package loki

import (
	"heroku-line-bot/src/model"
	"heroku-line-bot/src/pkg/util"
	"strings"
)

type Loki struct {
	host string
}

func New(host string) *Loki {
	return &Loki{
		host: host,
	}
}

func (l Loki) Send(reqs model.Reqs_Service_LokiSend) error {
	uri := l.host + "/loki/api/v1/push"
	method := util.POST
	// TODO http請求物件化來設定請求
	if _, err := util.SendJsonRequest(uri, method, reqs, nil); err != nil &&
		!strings.Contains(err.Error(), "not 200") {
		return err
	}

	return nil
}
