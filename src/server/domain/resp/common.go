package resp

type Base struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Page struct {
	DataCount int `json:"data_count"`
}
