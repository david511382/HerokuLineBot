package activity

import (
	courtBadmintonLogicDomain "heroku-line-bot/src/logic/badminton/court/domain"
	"heroku-line-bot/src/util"
)

type Activity struct {
	ID          int
	TeamID      int
	Date        util.DateTime
	PlaceID     int
	Courts      []*courtBadmintonLogicDomain.ActivityCourt
	MemberCount int16
	GuestCount  int16
	MemberFee   int16
	GuestFee    int16
	ClubSubsidy int16
	LogisticID  *int
	Description string
	PeopleLimit *int16
}
