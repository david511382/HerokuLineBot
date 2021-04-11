package resp

type GetUserProfile struct {
	DisplayName string `json:"displayName"`
}

type ReplyMessage struct{}

type PushMessage struct{}

type ListRichMenuRichMenu struct {
	Name       string `json:"name"`
	RichMenuID string `json:"richMenuId"`
}

type ListRichMenu struct {
	RichMenus []*ListRichMenuRichMenu `json:"richmenus"`
}

type DeleteRichMenu struct{}

type SetDefaultRichMenu struct{}

type SetRichMenuTo struct{}

type SetRichMenuTos struct{}

type GetDefaultRichMenu struct {
	RichMenuID string `json:"richMenuId"`
}

type UploadRichMenuImage struct{}

type CreateRichMenu struct {
	RichMenuID string `json:"richMenuId"`
}
