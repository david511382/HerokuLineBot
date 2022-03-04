package linebot

import (
	"heroku-line-bot/src/pkg/service/linebot/domain"
	"heroku-line-bot/src/pkg/service/linebot/domain/model"
)

func GetTextMessage(text string) *model.TextMessage {
	return &model.TextMessage{
		Type: &model.Type{
			Type: "text",
		},
		Text: text,
	}
}

func GetFlexMessage(altText string, content interface{}) *model.FlexMessage {
	return &model.FlexMessage{
		Type: &model.Type{
			Type: "flex",
		},
		AltText:  altText,
		Contents: content,
	}
}

func GetFlexMessageCarouselContent(contents ...*model.FlexMessagBubbleComponent) *model.FlexMessagCarouselComponent {
	result := &model.FlexMessagCarouselComponent{
		Type: &model.Type{
			Type: "carousel",
		},
		Contents: contents,
	}
	return result
}

func GetFlexMessageBubbleContent(body *model.FlexMessageBoxComponent, option *model.FlexMessagBubbleComponentOption) *model.FlexMessagBubbleComponent {
	result := &model.FlexMessagBubbleComponent{
		Type: &model.Type{
			Type: "bubble",
		},
		Body: body,
	}
	if option != nil {
		result.FlexMessagBubbleComponentOption = option
	}
	return result
}

func GetFlexMessageBoxComponent(layout domain.MessageLayout, option *model.FlexMessageBoxComponentOption, contents ...interface{}) *model.FlexMessageBoxComponent {
	result := &model.FlexMessageBoxComponent{
		Type: &model.Type{
			Type: "box",
		},
		Layout:                        layout,
		Contents:                      contents,
		FlexMessageBoxComponentOption: option,
	}

	return result
}

func GetPostBackAction(text, data string) *model.PostBackAction {
	return &model.PostBackAction{
		Type: &model.Type{
			Type: "postback",
		},
		Label: text,
		Data:  data,
	}
}

func GetMessageAction(text string) *model.MessageAction {
	return &model.MessageAction{
		Type: &model.Type{
			Type: "message",
		},
		Label: text,
		Text:  text,
	}
}

func GetUriAction(uri string) *model.UriAction {
	return &model.UriAction{
		Type: "uri",
		Uri:  uri,
	}
}

func GetTimeAction(text, data, max, min string, mode domain.TimeActionMode) *model.TimeAction {
	return &model.TimeAction{
		PostBackAction: &model.PostBackAction{
			Type: &model.Type{
				Type: "datetimepicker",
			},
			Label: text,
			Data:  data,
		},
		Mode: mode,
		Max:  max,
		Min:  min,
	}
}

func GetButtonComponent(action interface{}, option *model.ButtonOption) *model.Button {
	result := &model.Button{
		Type: &model.Type{
			Type: "button",
		},
		ButtonOption: option,
		Action:       action,
	}

	return result
}

func GetClassButtonComponent(action interface{}) *model.Button {
	option := &model.ButtonOption{
		Style:      "primary",
		Height:     domain.SM_FLEX_MESSAGE_SIZE,
		AdjustMode: domain.SHRINK_TO_FIT_ADJUST_MODE,
	}
	result := &model.Button{
		Type: &model.Type{
			Type: "button",
		},
		ButtonOption: option,
		Action:       action,
	}

	return result
}

func GetFlexMessageTextComponent(text string, option *model.FlexMessageTextComponentOption) *model.FlexMessageTextComponent {
	return &model.FlexMessageTextComponent{
		TextMessage:                    *GetTextMessage(text),
		FlexMessageTextComponentOption: option,
		//	AdjustMode:                     domain.SHRINK_TO_FIT_ADJUST_MODE,
		//	Align:                          "start",
	}
}

func GetFlexMessageTextComponentSpan(text string, size domain.MessageSize, weight domain.MessageWeight) *model.FlexMessageTextComponentSpan {
	return &model.FlexMessageTextComponentSpan{
		TextMessage: model.TextMessage{
			Type: &model.Type{
				Type: "span",
			},
			Text: text,
		},
		Size:   size,
		Weight: weight,
	}
}

func GetSeparatorComponent(option *model.FlexMessageSeparatorComponentOption) *model.FlexMessageSeparatorComponent {
	return &model.FlexMessageSeparatorComponent{
		Type: &model.Type{
			Type: "separator",
		},
		FlexMessageSeparatorComponentOption: option,
	}
}

func GetFillerComponent() *model.FlexMessageFillerComponent {
	return &model.FlexMessageFillerComponent{
		Type: &model.Type{
			Type: "filler",
		},
	}
}
