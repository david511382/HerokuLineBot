package domain

import "heroku-line-bot/src/logic/club/domain"

type Model struct {
	ID   int             `json:"id,omitempty"`
	Name string          `json:"name,omitempty"`
	Role domain.ClubRole `json:"role,omitempty"`
}
