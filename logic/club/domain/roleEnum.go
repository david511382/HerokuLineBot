package domain

type ClubRole int16

const (
	ADMIN_CLUB_ROLE  ClubRole = 1
	CADRE_CLUB_ROLE  ClubRole = 2
	MEMBER_CLUB_ROLE ClubRole = 3
	GUEST_CLUB_ROLE  ClubRole = 4
)

var ClubRoleNames = []string{
	"管理者",
	"幹部",
	"社員",
	"球友",
}
