package domain

type MessageEventMessageType string

var (
	TEXT_MESSAGE_EVENT_MESSAGE_TYPE     MessageEventMessageType = "text"
	IMAGE_MESSAGE_EVENT_MESSAGE_TYPE    MessageEventMessageType = "image"
	VIDEO_MESSAGE_EVENT_MESSAGE_TYPE    MessageEventMessageType = "video"
	AUDIO_MESSAGE_EVENT_MESSAGE_TYPE    MessageEventMessageType = "audio"
	LOCATION_MESSAGE_EVENT_MESSAGE_TYPE MessageEventMessageType = "location"
	STICKER_MESSAGE_EVENT_MESSAGE_TYPE  MessageEventMessageType = "sticker"
)
