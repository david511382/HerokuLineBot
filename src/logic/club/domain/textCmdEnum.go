package domain

type TextCmd string

const (
	NEW_ACTIVITY_TEXT_CMD         TextCmd = "新增活動"
	GET_ACTIVITIES_TEXT_CMD       TextCmd = "查看活動"
	REGISTE_COMPANY_TEXT_CMD      TextCmd = "登記公司"
	CONFIRM_REGISTER_TEXT_CMD     TextCmd = "確認入社"
	GET_CONFIRM_REGISTER_TEXT_CMD TextCmd = "入社列表"
	SUBMIT_ACTIVITY_TEXT_CMD      TextCmd = "提交活動"
	RICH_MENU_TEXT_CMD            TextCmd = "richmenu"
	NEW_LOGISTIC_TEXT_CMD         TextCmd = "新增品項紀錄"
	UPDATE_MEMBER_INFO_TEXT_CMD   TextCmd = "修改個人資訊"
)

type DateTimeCmd string

const (
	TIME_POSTBACK_DATE_TIME_CMD      DateTimeCmd = "time"
	DATE_TIME_POSTBACK_DATE_TIME_CMD DateTimeCmd = "date_time"
	DATE_POSTBACK_DATE_TIME_CMD      DateTimeCmd = "date"
)
