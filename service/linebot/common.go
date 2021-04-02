package linebot

func New(channelAccessToken string) *LineBot {
	return &LineBot{
		channelAccessToken: channelAccessToken,
	}
}
