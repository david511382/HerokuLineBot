package linebot

import (
	"heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model"
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
		Layout:   layout,
		Contents: contents,
	}

	if option != nil {
		result.FlexMessageBoxComponentOption = option
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

func GetButtonComponent(flex int, action interface{}, option *model.ButtonOption) *model.Button {
	result := &model.Button{
		Type: &model.Type{
			Type: "button",
		},
		ButtonOption: option,
		Action:       action,
		Style:        "primary",
		Height:       domain.SM_FLEX_MESSAGE_SIZE,
		Flex:         flex,
		AdjustMode:   domain.SHRINK_TO_FIT_ADJUST_MODE,
	}

	return result
}

func GetFlexMessageTextComponent(flex int, text string, contents ...*model.FlexMessageTextComponentSpan) *model.FlexMessageTextComponent {
	return &model.FlexMessageTextComponent{
		TextMessage: *GetTextMessage(text),
		Contents:    contents,
		Flex:        flex,
		AdjustMode:  domain.SHRINK_TO_FIT_ADJUST_MODE,
		Align:       "start",
	}
}

func GetFlexMessageTextComponentSpan(text string, size domain.MessageSize, weight domain.MessageWeight) *model.FlexMessageTextComponentSpan {
	return &model.FlexMessageTextComponentSpan{
		Type: &model.Type{
			Type: "span",
		},
		Text:   text,
		Size:   size,
		Weight: weight,
	}
}

func GetSeparatorComponent(colorP *string) *model.FlexMessageSeparatorComponent {
	color := "#000000"
	if colorP != nil {
		color = *colorP
	}
	return &model.FlexMessageSeparatorComponent{
		Type: &model.Type{
			Type: "separator",
		},
		Color: color,
	}
}

func GetFillerComponent() *model.FlexMessageFillerComponent {
	return &model.FlexMessageFillerComponent{
		Type: &model.Type{
			Type: "filler",
		},
	}
}
