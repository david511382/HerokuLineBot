package resp

type UserInfo struct {
	RoleID   uint8  `json:"role_id"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
}
