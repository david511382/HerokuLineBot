package domain

import clubLogicDomain "heroku-line-bot/logic/club/domain"

type Model struct {
	ID   int                      `json:"id,omitempty"`
	Name string                   `json:"name,omitempty"`
	Role clubLogicDomain.ClubRole `json:"role,omitempty"`
}
