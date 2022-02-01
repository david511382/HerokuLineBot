package domain

import "fmt"

var (
	USER_NOT_REGISTERED error = fmt.Errorf("用戶尚未註冊")
	NO_AUTH_ERROR       error = fmt.Errorf("沒有權限")
)
