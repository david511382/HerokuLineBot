package telegram

func NewBot(token string) *Bot {
	return &Bot{
		botToken: token,
	}
}
