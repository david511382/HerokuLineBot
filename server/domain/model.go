package domain

import (
	"time"
)

type JwtClaims struct {
	RoleID   int16
	ID       int
	Username string
	ExpTime  time.Time
}
