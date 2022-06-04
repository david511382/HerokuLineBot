package flow

import (
	"time"
)

type StepRecord struct {
	Name         string
	DurationTime time.Duration
}

type StepRecords []*StepRecord

func (sr StepRecords) Last() *StepRecord {
	len := len(sr)
	if len == 0 {
		return nil
	}

	return sr[len-1]
}

func (sr StepRecords) TotalDuration() time.Duration {
	dur := time.Duration(0)
	for _, v := range sr {
		dur += v.DurationTime
	}
	return dur
}
