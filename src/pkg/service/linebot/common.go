package linebot

func New(channelAccessToken string) *LineBot {
	return &LineBot{
		channelAccessToken: channelAccessToken,
	}
}

func NewOAuth(channelID uint64) *OAuth {
	return &OAuth{
		channelID: channelID,
	}
}
