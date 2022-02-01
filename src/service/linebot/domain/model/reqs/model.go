package reqs

type OAuthVerifyIDToken struct {
	// ID token
	IDToken string `json:"id_token"`
	// Expected channel ID. Unique identifier for your channel issued by LINE. Found in the LINE Developers Console.
	ClientID string `json:"client_id"`
	// Expected nonce value. Use the nonce value provided in the authorization request. Omit if the nonce value was not specified in the authorization request.
	Nonce string `json:"nonce"`
	// Expected user ID. Learn how to get the user ID from Get user profile.
	UserID string `json:"user_id"`
}

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
