package reqs

type ReplyMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

type PushMessage struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}

type SetRichMenuTos struct {
	RichMenuID string   `json:"richMenuId"`
	UserID     []string `json:"userId"`
}

type CreateRichMenuSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type CreateRichMenuAreasBounds struct {
	X int `json:"x"`
	Y int `json:"y"`
	CreateRichMenuSize
}

type CreateRichMenuAreas struct {
	Action interface{}               `json:"action"`
	Bounds CreateRichMenuAreasBounds `json:"bounds"`
}

type CreateRichMenu struct {
	Size        CreateRichMenuSize     `json:"size"`
	Selected    bool                   `json:"selected"`
	Name        string                 `json:"name"`
	ChatBarText string                 `json:"chatBarText"`
	Areas       []*CreateRichMenuAreas `json:"areas"`
}
