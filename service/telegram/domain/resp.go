package domain

type RespSendMessage struct {
}

type RespGetUpdatesResultMessageChat struct {
	ID int `json:"id"`
}
type RespGetUpdatesResultMessage struct {
	Chat RespGetUpdatesResultMessageChat `json:"chat"`
}
type RespGetUpdatesResult struct {
	Message RespGetUpdatesResultMessage `json:"message"`
}
type RespGetUpdates struct {
	Result []RespGetUpdatesResult `json:"result"`
}
