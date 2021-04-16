package domain

type TextCmd string

const (
	NEW_ACTIVITY_TEXT_CMD     TextCmd = "新增活動"
	GET_ACTIVITIES_TEXT_CMD   TextCmd = "查看活動"
	REGISTER_TEXT_CMD         TextCmd = "註冊"
	CONFIRM_REGISTER_TEXT_CMD TextCmd = "確認入社"
	SUBMIT_ACTIVITY_TEXT_CMD  TextCmd = "提交活動"
	RICH_MENU_TEXT_CMD        TextCmd = "richmenu"
)

type DateTimeCmd string

const (
	TIME_POSTBACK_DATE_TIME_CMD      DateTimeCmd = "time"
	DATE_TIME_POSTBACK_DATE_TIME_CMD DateTimeCmd = "date_time"
	DATE_POSTBACK_DATE_TIME_CMD      DateTimeCmd = "date"
)
