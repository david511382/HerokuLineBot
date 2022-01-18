package database

type ReqsClubJoinActivityDetail struct {
	*ReqsClubActivity
}

type RespClubJoinActivityDetail struct {
	ActivityID                 int
	RentalCourtDetailStartTime string
	RentalCourtDetailEndTime   string
	RentalCourtDetailCount     int
}
