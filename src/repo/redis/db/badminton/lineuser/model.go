package lineuser

type LineUser struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Role int16  `json:"role,omitempty"`
}
