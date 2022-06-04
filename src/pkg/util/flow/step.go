package flow

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
	"time"
)

type IStep interface {
	Name() string
	Run() (records StepRecords, resultErrInfo errUtil.IError)
}

type Step struct {
	StepName string
	Fun      func() (resultErrInfo errUtil.IError)
}

func (s Step) Name() string {
	return s.StepName
}

func (s Step) Run() (records StepRecords, resultErrInfo errUtil.IError) {
	startTime := time.Now()
	resultErrInfo = s.Fun()
	durationTime := time.Since(startTime)
	record := &StepRecord{
		Name:         s.Name(),
		DurationTime: durationTime,
	}
	records = append(records, record)
	return
}

type Steps []IStep
