package heroku

import (
	"heroku-line-bot/bootstrap"
	errLogic "heroku-line-bot/logic/error"
	"heroku-line-bot/util"
	"time"
)

type BackGround struct {
	url string
}

func (b *BackGround) Init(cfg bootstrap.Backgrounds) (name string, backgroundCfg bootstrap.Background, resultErrInfo errLogic.IError) {
	b.url = cfg.Heroku.Url
	return "Heroku", cfg.Heroku.Background, nil
}

func (b *BackGround) Run(runTime time.Time) (resultErrInfo errLogic.IError) {
	if b.url != "" {
		util.SendGetRequest(b.url, nil, nil)
	}
	return nil
}
