package domain

type ClubRole int16

const (
	ADMIN_CLUB_ROLE  ClubRole = 1
	CADRE_CLUB_ROLE  ClubRole = 2
	MEMBER_CLUB_ROLE ClubRole = 3
	GUEST_CLUB_ROLE  ClubRole = 4
)

func (r ClubRole) Name() string {
	switch r {
	case ADMIN_CLUB_ROLE:
		return "管理者"
	case CADRE_CLUB_ROLE:
		return "幹部"
	case MEMBER_CLUB_ROLE:
		return "社員"
	case GUEST_CLUB_ROLE:
		return "球友"
	default:
		return "未定義"
	}
}
