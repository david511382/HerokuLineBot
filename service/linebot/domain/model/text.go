package model

import "heroku-line-bot/service/linebot/domain"

type TextMessage struct {
	*Type
	Text   string               `json:"text,omitempty"`
	Weight domain.MessageWeight `json:"weight,omitempty"`
	Size   domain.MessageSize   `json:"size,omitempty"`
}

type FlexMessage struct {
	*Type
	AltText  string      `json:"altText,omitempty"`
	Contents interface{} `json:"contents,omitempty"`
}

type FlexMessageBoxComponentOption struct {
	AlignItems     domain.AlignItems     `json:"alignItems,omitempty"`
	JustifyContent domain.JustifyContent `json:"justifyContent,omitempty"`
}

type FlexMessageBoxComponent struct {
	*Type
	*FlexMessageBoxComponentOption
	Flex     string               `json:"flex,omitempty"`
	Layout   domain.MessageLayout `json:"layout,omitempty"`
	Contents []interface{}        `json:"contents,omitempty"`
}

type Background struct {
	BackgroundColor string `json:"backgroundColor,omitempty"`
	SeparatorColor  string `json:"separatorColor,omitempty"`
	Separator       bool   `json:"separator,omitempty"`
}

type FlexMessagBubbleComponentStyle struct {
	Header *Background `json:"header,omitempty"`
	Body   *Background `json:"body,omitempty"`
}

type FlexMessagBubbleComponentOption struct {
	Header *FlexMessageBoxComponent        `json:"header,omitempty"`
	Styles *FlexMessagBubbleComponentStyle `json:"styles,omitempty"`
}

type FlexMessagBubbleComponent struct {
	*Type
	*FlexMessagBubbleComponentOption
	Body *FlexMessageBoxComponent `json:"body,omitempty"`
}

type FlexMessagCarouselComponent struct {
	*Type
	Contents []*FlexMessagBubbleComponent `json:"contents,omitempty"`
}

type FlexMessageTextComponent struct {
	*Type
	Contents   []*FlexMessageTextComponentSpan `json:"contents,omitempty"`
	Flex       int                             `json:"flex,omitempty"`
	AdjustMode domain.AdjustMode               `json:"adjustMode,omitempty"`
	Align      string                          `json:"align,omitempty"`
}

type FlexMessageTextComponentSpan struct {
	*Type
	Text   string               `json:"text,omitempty"`
	Size   domain.MessageSize   `json:"size,omitempty"`
	Weight domain.MessageWeight `json:"weight,omitempty"`
}

type PostBackAction struct {
	*Type
	Label string `json:"label,omitempty"`
	Data  string `json:"data,omitempty"`
}

type TimeAction struct {
	*PostBackAction
	Mode domain.TimeActionMode `json:"mode,omitempty"`
	Max  string                `json:"max,omitempty"`
	Min  string                `json:"min,omitempty"`
}

type Button struct {
	*Type
	Action     interface{}        `json:"action,omitempty"`
	Style      string             `json:"style,omitempty"`
	Color      string             `json:"color,omitempty"`
	Height     domain.MessageSize `json:"height,omitempty"`
	Flex       int                `json:"flex,omitempty"`
	AdjustMode domain.AdjustMode  `json:"adjustMode,omitempty"`
}

type FlexMessageSeparatorComponent struct {
	*Type
	Color string `json:"color,omitempty"`
}
