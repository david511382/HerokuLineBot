package domain

type ReqsSendMessage struct {
	ChatID int    `url:"chat_id"`
	Text   string `url:"text"`
}

type ReqsGetUpdates struct {
}
