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

func GetFlexMessage(altText string, content *model.FlexMessagBubbleComponent) *model.FlexMessage {
	return &model.FlexMessage{
		Type: &model.Type{
			Type: "flex",
		},
		AltText:  altText,
		Contents: content,
	}
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

func GetButtonComponent(flex int, action interface{}) *model.Button {
	return &model.Button{
		Type: &model.Type{
			Type: "button",
		},
		Action:     action,
		Style:      "primary",
		Color:      "#00dd00",
		Height:     domain.SM_FLEX_MESSAGE_SIZE,
		Flex:       flex,
		AdjustMode: domain.SHRINK_TO_FIT_ADJUST_MODE,
	}
}

func GetKeyValueEditComponent(name, value string, action interface{}, sizeP, valueSizeP *domain.MessageSize) *model.FlexMessageBoxComponent {
	size := domain.XL_FLEX_MESSAGE_SIZE
	if sizeP != nil {
		size = *sizeP
	}
	valueSize := size
	if valueSizeP != nil {
		valueSize = *valueSizeP
	}

	contents := []interface{}{
		GetFlexMessageTextComponent(
			5,
			GetFlexMessageTextComponentSpan(name, size, domain.BOLD_FLEX_MESSAGE_WEIGHT),
			GetFlexMessageTextComponentSpan(" : ", size, domain.BOLD_FLEX_MESSAGE_WEIGHT),
			GetFlexMessageTextComponentSpan(value, valueSize, domain.REGULAR_FLEX_MESSAGE_WEIGHT),
		),
	}
	if action != nil {
		contents = append(contents, GetButtonComponent(2, action))
	}
	return GetFlexMessageBoxComponent(
		domain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		contents...,
	)
}

func GetComfirmComponent(leftAction, rightAction interface{}) *model.FlexMessageBoxComponent {
	return GetFlexMessageBoxComponent(
		domain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		GetButtonComponent(0, leftAction),
		GetButtonComponent(0, rightAction),
	)
}

func GetFlexMessageTextComponent(flex int, contents ...*model.FlexMessageTextComponentSpan) *model.FlexMessageTextComponent {
	return &model.FlexMessageTextComponent{
		Type: &model.Type{
			Type: "text",
		},
		Contents:   contents,
		Flex:       flex,
		AdjustMode: domain.SHRINK_TO_FIT_ADJUST_MODE,
		Align:      "start",
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
