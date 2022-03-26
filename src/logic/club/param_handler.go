package club

import (
	"fmt"
	"heroku-line-bot/src/logger"
	"heroku-line-bot/src/logic/club/domain"
	"heroku-line-bot/src/pkg/service/linebot"
	linebotDomain "heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/service/linebot/domain/model"
	errUtil "heroku-line-bot/src/pkg/util/error"
)

type ParamHandler struct {
	requireAttr               string
	attrNameText              string
	valueText                 string
	isNotRequireChecking      bool
	messages                  interface{}
	loadRequireInputTextParam func(requireAttr, text string) errUtil.IError
}

func NewParamHandler(requireAttr string, param domain.IParamTextValue) *ParamHandler {
	result := &ParamHandler{
		requireAttr:               requireAttr,
		loadRequireInputTextParam: param.LoadRequireInputTextParam,
	}
	result.attrNameText, result.valueText, result.isNotRequireChecking = param.GetRequireAttrInfo(result.requireAttr)
	result.messages = param.GetInputTemplate(result.requireAttr)
	return result
}

func (b ParamHandler) IsReading() bool {
	return b.requireAttr != ""
}

func (b *ParamHandler) Read(text string) (resultErrInfo errUtil.IError) {
	if errInfo := b.loadRequireInputTextParam(b.requireAttr, text); errInfo != nil {
		resultErrInfo = errUtil.Append(resultErrInfo, errInfo)
		if errInfo.IsError() {
			return
		}
	}

	return
}

func (b ParamHandler) GetInputTemplate() (resultMessages interface{}) {
	const altText = "請確認或輸入"

	if b.messages != nil {
		resultMessages = b.messages
		return
	}

	text := ""
	if b.valueText != "" {
		if b.attrNameText != "" {
			format := "請確認%s %s，或輸入數值"
			if b.isNotRequireChecking {
				format = "%s %s，請輸入數值"
			}
			text = fmt.Sprintf(format, b.attrNameText, b.valueText)
		} else {
			format := "請確認%s，或輸入數值"
			if b.isNotRequireChecking {
				format = "數值%s，請輸入"
			}
			text = fmt.Sprintf(format, b.valueText)
		}
	} else if b.attrNameText != "" {
		format := "請確認或輸入%s數值"
		if b.isNotRequireChecking {
			format = "請輸入%s數值"
		}
		text = fmt.Sprintf(format, b.attrNameText)
	}

	cancelRequireInputJs, errInfo := NewSignal().
		GetCancelInputMode().
		GetSignal()
	if errInfo != nil {
		logger.Log("line bot club", errInfo)
		return
	}

	contents := make([]interface{}, 0)
	if !b.isNotRequireChecking {
		checkButton := linebot.GetButtonComponent(
			linebot.GetPostBackAction(
				"確認",
				cancelRequireInputJs,
			),
			&domain.NormalButtonOption,
		)
		contents = append(contents, checkButton)
	}
	resultMessages = linebot.GetFlexMessage(
		altText,
		linebot.GetFlexMessageBubbleContent(
			linebot.GetFlexMessageBoxComponent(
				linebotDomain.VERTICAL_MESSAGE_LAYOUT,
				&model.FlexMessageBoxComponentOption{
					JustifyContent: linebotDomain.SPACE_EVENLY_JUSTIFY_CONTENT,
				},
				contents...,
			),
			&model.FlexMessagBubbleComponentOption{
				Header: linebot.GetFlexMessageBoxComponent(
					linebotDomain.VERTICAL_MESSAGE_LAYOUT,
					nil,
					linebot.GetTextMessage(text),
				),
				Styles: &model.FlexMessagBubbleComponentStyle{
					Header: &model.Background{
						BackgroundColor: "#8DFF33",
					},
					Body: &model.Background{
						BackgroundColor: "#FFFFFF",
						SeparatorColor:  "#000000",
						Separator:       true,
					},
				},
			},
		),
	)
	return
}
