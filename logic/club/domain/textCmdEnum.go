package domain

type TextCmd string

const (
	NEW_ACTIVITY_TEXT_CMD    TextCmd = "新增活動"
	GET_ACTIVITIES_TEXT_CMD  TextCmd = "查看活動"
	REGISTER_TEXT_CMD        TextCmd = "註冊"
	SUBMIT_ACTIVITY_TEXT_CMD TextCmd = "提交活動"
	RICH_MENU_TEXT_CMD       TextCmd = "richmenu"
	TIME_POSTBACK_CMD        TextCmd = "time"
	DATE_TIME_POSTBACK_CMD   TextCmd = "date_time"
	DATE_POSTBACK_CMD        TextCmd = "date"
)
