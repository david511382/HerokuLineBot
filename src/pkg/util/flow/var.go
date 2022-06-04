package flow

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
)

var (
	FeatureFieldName = "feature"
	StepFieldName    = "step"
)

var (
	ErrorBreak = errUtil.New("Break")
)
