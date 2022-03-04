package model

import "heroku-line-bot/src/pkg/service/linebot/domain"

type TextMessage struct {
	*Type
	Text string `json:"text,omitempty"`
}

type FlexMessage struct {
	*Type
	AltText  string      `json:"altText,omitempty"`
	Contents interface{} `json:"contents,omitempty"`
}

type FlexMessageBoxComponentOption struct {
	AlignItems      domain.AlignItems     `json:"alignItems,omitempty"`
	JustifyContent  domain.JustifyContent `json:"justifyContent,omitempty"`
	Margin          domain.MessageSize    `json:"margin,omitempty"`
	Spacing         domain.MessageSize    `json:"spacing,omitempty"`
	BackgroundColor string                `json:"backgroundColor,omitempty"`
	CornerRadius    string                `json:"cornerRadius,omitempty"`
	Height          string                `json:"height,omitempty"`
	Width           string                `json:"width,omitempty"`
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
	Footer *Background `json:"footer,omitempty"`
}

type FlexMessagBubbleComponentOption struct {
	Header *FlexMessageBoxComponent        `json:"header,omitempty"`
	Footer *FlexMessageBoxComponent        `json:"footer,omitempty"`
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

type FlexMessageTextComponentOption struct {
	Contents   []*FlexMessageTextComponentSpan `json:"contents,omitempty"`
	Weight     domain.MessageWeight            `json:"weight,omitempty"`
	Color      string                          `json:"color,omitempty"`
	Size       domain.MessageSize              `json:"size,omitempty"`
	Flex       int                             `json:"flex,omitempty"`
	AdjustMode domain.AdjustMode               `json:"adjustMode,omitempty"`
	Align      domain.Align                    `json:"align,omitempty"`
	Margin     domain.MessageSize              `json:"margin,omitempty"`
	Wrap       bool                            `json:"wrap,omitempty"`
}

type FlexMessageTextComponent struct {
	TextMessage
	*FlexMessageTextComponentOption
}

type FlexMessageTextComponentSpan struct {
	TextMessage
	Weight domain.MessageWeight `json:"weight,omitempty"`
	Color  string               `json:"color,omitempty"`
	Size   domain.MessageSize   `json:"size,omitempty"`
}

type PostBackAction struct {
	*Type
	Label string `json:"label,omitempty"`
	Data  string `json:"data,omitempty"`
}

type MessageAction struct {
	*Type
	Label string `json:"label,omitempty"`
	Text  string `json:"text,omitempty"`
}

type UriAction struct {
	Type string `json:"type"`
	Uri  string `json:"uri,omitempty"`
}

type TimeAction struct {
	*PostBackAction
	Mode domain.TimeActionMode `json:"mode,omitempty"`
	Max  string                `json:"max,omitempty"`
	Min  string                `json:"min,omitempty"`
}

type ButtonOption struct {
	Color      string             `json:"color,omitempty"`
	Flex       int                `json:"flex,omitempty"`
	Style      string             `json:"style,omitempty"`
	Height     domain.MessageSize `json:"height,omitempty"`
	AdjustMode domain.AdjustMode  `json:"adjustMode,omitempty"`
}

type Button struct {
	*Type
	*ButtonOption
	Action interface{} `json:"action,omitempty"`
}

type FlexMessageSeparatorComponentOption struct {
	Color  string             `json:"color,omitempty"`
	Margin domain.MessageSize `json:"margin,omitempty"`
}

type FlexMessageSeparatorComponent struct {
	*Type
	*FlexMessageSeparatorComponentOption
}

type FlexMessageFillerComponent struct {
	*Type
	Flex int `json:"flex,omitempty"`
}
