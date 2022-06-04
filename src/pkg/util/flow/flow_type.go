package flow

import (
	errUtil "heroku-line-bot/src/pkg/util/error"
)

type flowType struct {
	feature string
	steps   []IStep
}

func Flow(功能說明 string, 聲明步驟 ...IStep) flowType {
	return flowType{
		feature: 功能說明,
		steps:   聲明步驟,
	}
}

func (df flowType) Name() string {
	return df.feature
}

func (df flowType) Run() (resultRecords StepRecords, resultErrInfo errUtil.IError) {
	for _, f := range df.steps {
		records, errInfo := f.Run()
		for _, record := range records {
			record.Name = df.stepName(record.Name)
		}
		resultRecords = append(resultRecords, records...)

		if errInfo != nil {
			if errUtil.Equal(errInfo, ErrorBreak) {
				return
			}

			errInfo.Attr(FeatureFieldName, df.feature)
			errInfo.Attr(StepFieldName, f.Name())
			resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
			if resultErrInfo.IsError() {
				return
			}
		}
	}

	return
}

func (df flowType) stepName(stepName string) string {
	return df.Name() + " > " + stepName
}
