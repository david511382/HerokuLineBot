package domain

import "fmt"

var (
	USER_NOT_REGISTERED error = fmt.Errorf("用戶尚未註冊")
)
