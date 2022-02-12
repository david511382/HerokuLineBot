package telegram

func NewBot(token string) *Bot {
	if token == "" {
		return nil
	}
	return &Bot{
		botToken: token,
	}
}
