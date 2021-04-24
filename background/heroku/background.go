package heroku

import (
	"heroku-line-bot/bootstrap"
	"heroku-line-bot/util"
	"time"
)

type BackGround struct {
	url string
}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErr error) {
	b.url = cfg.Heroku.Url

	return "heroku", cfg.Heroku.Background, nil
}

func (b *BackGround) Run(runTime time.Time) error {
	if b.url != "" {
		util.SendGetRequest(b.url, nil, nil)
	}

	return nil
}
