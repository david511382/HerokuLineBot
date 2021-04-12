package background

import (
	"heroku-line-bot/background/heroku"
	"heroku-line-bot/bootstrap"

	cron "gopkg.in/robfig/cron.v2"
)

var (
	cr = cron.New()
)

func Init(cfg *bootstrap.Config) error {
	herokuBg := &heroku.Heroku{}
	spec := herokuBg.Init(*cfg)
	_, err := cr.AddJob(spec, herokuBg)
	if err != nil {
		return err
	}

	cr.Start()

	return nil
}
