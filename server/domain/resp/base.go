package resp

type Base struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
