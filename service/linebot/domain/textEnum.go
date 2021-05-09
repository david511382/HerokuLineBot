package domain

type MessageLayout string

const (
	VERTICAL_MESSAGE_LAYOUT   MessageLayout = "vertical"
	HORIZONTAL_MESSAGE_LAYOUT MessageLayout = "horizontal"
)

type FlexMessageContentType string

const (
	BUBBLE_FLEX_MESSAGE_CONTENT_TYPE FlexMessageContentType = "bubble"
)

type MessageWeight string

const (
	BOLD_FLEX_MESSAGE_WEIGHT    MessageWeight = "bold"
	REGULAR_FLEX_MESSAGE_WEIGHT MessageWeight = "regular"
)

type MessageSize string

const (
	LG_FLEX_MESSAGE_SIZE  MessageSize = "lg"
	XXL_FLEX_MESSAGE_SIZE MessageSize = "xxl"
	XL_FLEX_MESSAGE_SIZE  MessageSize = "xl"
	MD_FLEX_MESSAGE_SIZE  MessageSize = "md"
	XS_FLEX_MESSAGE_SIZE  MessageSize = "xs"
	SM_FLEX_MESSAGE_SIZE  MessageSize = "sm"
)

type AdjustMode string

const (
	SHRINK_TO_FIT_ADJUST_MODE AdjustMode = "shrink-to-fit"
)

type Align string

const (
	START_Align  Align = "start"
	END_Align    Align = "end"
	CENTER_Align Align = "center"
)

type TimeActionMode string

const (
	DATE_TIME_ACTION_MODE TimeActionMode = "date"
)
