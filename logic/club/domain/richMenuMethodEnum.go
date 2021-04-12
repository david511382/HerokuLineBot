package domain

type RichMenuMethod string

const (
	LIST_RICH_MENU_METHOD        RichMenuMethod = "List"
	DELETE_RICH_MENU_METHOD      RichMenuMethod = "Delete"
	NEW_RICH_MENU_METHOD         RichMenuMethod = "New"
	SET_DEFAULT_RICH_MENU_METHOD RichMenuMethod = "Set Default"
	SET_TO_RICH_MENU_METHOD      RichMenuMethod = "Set To"
)
