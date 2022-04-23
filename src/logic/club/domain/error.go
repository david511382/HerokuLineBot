package domain

import "fmt"

var (
	NO_AUTH_ERROR error = fmt.Errorf("沒有權限")
)
