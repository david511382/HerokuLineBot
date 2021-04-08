package club

import (
	"heroku-line-bot/logic/club/domain"
	"heroku-line-bot/service/linebot"
	linebotDomain "heroku-line-bot/service/linebot/domain"
	"heroku-line-bot/service/linebot/domain/model"
	linebotModel "heroku-line-bot/service/linebot/domain/model"
)

func GetDoubleKeyValueComponent(key1, value1, key2, value2 string, option *linebotModel.FlexMessageBoxComponentOption, keyValueEditComponentOption *domain.KeyValueEditComponentOption) *linebotModel.FlexMessageBoxComponent {
	components := []interface{}{}
	components = append(components, GetKeyValueEditComponent(
		key1,
		value1,
		keyValueEditComponentOption,
	))
	components = append(components, GetKeyValueEditComponent(
		key2,
		value2,
		keyValueEditComponentOption,
	))
	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		option,
		components...,
	)
}

func GetKeyValueEditComponent(name, value string, option *domain.KeyValueEditComponentOption) *model.FlexMessageBoxComponent {
	size := linebotDomain.XL_FLEX_MESSAGE_SIZE
	if option != nil && option.SizeP != nil {
		size = *option.SizeP
	}
	valueSize := size
	if option != nil && option.ValueSizeP != nil {
		valueSize = *option.ValueSizeP
	}

	contents := []interface{}{}

	if option != nil && option.Indent != nil {
		for i := 0; i < *option.Indent; i++ {
			contents = append(contents, linebot.GetFillerComponent())
		}
	}

	contents = append(contents,
		linebot.GetFlexMessageTextComponent(
			5,
			"",
			linebot.GetFlexMessageTextComponentSpan(name, size, linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT),
			linebot.GetFlexMessageTextComponentSpan(" : ", size, linebotDomain.BOLD_FLEX_MESSAGE_WEIGHT),
			linebot.GetFlexMessageTextComponentSpan(value, valueSize, linebotDomain.REGULAR_FLEX_MESSAGE_WEIGHT),
		),
	)

	if option != nil && option.Action != nil {
		contents = append(contents, linebot.GetButtonComponent(2, option.Action, &domain.NormalButtonOption))
	}

	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		contents...,
	)
}

func GetComfirmComponent(leftAction, rightAction interface{}) *model.FlexMessageBoxComponent {
	return linebot.GetFlexMessageBoxComponent(
		linebotDomain.HORIZONTAL_MESSAGE_LAYOUT,
		nil,
		linebot.GetButtonComponent(0, leftAction, &domain.NormalButtonOption),
		linebot.GetButtonComponent(0, rightAction, &domain.NormalButtonOption),
	)
}
