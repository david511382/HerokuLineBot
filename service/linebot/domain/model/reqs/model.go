package reqs

type ReplyMessage struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

type PushMessage struct {
	To       string        `json:"to"`
	Messages []interface{} `json:"messages"`
}
