package model

type Reqs_Service_LokiSend struct {
	Streams []*Reqs_Service_LokiSendStream `json:"streams"`
}

type Reqs_Service_LokiSendStream struct {
	Stream interface{} `json:"stream"`
	Values [][]string  `json:"values"`
}
