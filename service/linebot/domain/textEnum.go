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
	XL_FLEX_MESSAGE_SIZE MessageSize = "xl"
	MD_FLEX_MESSAGE_SIZE MessageSize = "md"
	SM_FLEX_MESSAGE_SIZE MessageSize = "sm"
)

type AdjustMode string

const (
	SHRINK_TO_FIT_ADJUST_MODE AdjustMode = "shrink-to-fit"
)

type TimeActionMode string

const (
	DATE_TIME_ACTION_MODE TimeActionMode = "date"
)
