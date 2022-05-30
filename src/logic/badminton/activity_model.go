package badminton

import (
	"heroku-line-bot/src/logic/badminton/domain"
	"heroku-line-bot/src/pkg/util"
)

type Activity struct {
	ID          int
	TeamID      int
	Date        util.DefinedTime[util.DateInt]
	PlaceID     int
	Courts      []*domain.ActivityCourt
	MemberCount int16
	GuestCount  int16
	MemberFee   int16
	GuestFee    int16
	ClubSubsidy int16
	LogisticID  *int
	Description string
	PeopleLimit *int16
}
