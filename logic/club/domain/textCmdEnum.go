package domain

type TextCmd string

const (
	NEW_ACTIVITY_TEXT_CMD  TextCmd = "新增活動"
	TIME_POSTBACK_CMD      TextCmd = "time"
	DATE_TIME_POSTBACK_CMD TextCmd = "date_time"
	DATE_POSTBACK_CMD      TextCmd = "date"
)
