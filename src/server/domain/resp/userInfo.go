package resp

type UserInfo struct {
	RoleID   int16  `json:"role_id"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}
