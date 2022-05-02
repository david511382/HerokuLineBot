package domain

import (
	"time"
)

type JwtClaims struct {
	RoleID   uint8
	ID       uint
	Username string
	ExpTime  time.Time
}
