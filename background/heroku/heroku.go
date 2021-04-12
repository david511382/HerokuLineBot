package heroku

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/util"
)

type Heroku struct {
	url string
}

func (b *Heroku) Init(cfg bootstrap.Config) string {
	b.url = cfg.Heroku.Url
	return cfg.Heroku.Spec
}

func (b *Heroku) Run() {
	if b.url != "" {
		util.SendGetRequest(b.url, nil, nil)
	}
}
