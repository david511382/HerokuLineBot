package domain

type ApiMethod string

const (
	SEND_MESSAGE_API_METHOD ApiMethod = "sendMessage"
	GET_UPDATES_API_METHOD  ApiMethod = "getUpdates"
)
