package reqs

type GetRentalCourts struct {
	MustFromToDate
	TeamID int `json:"team_id" form:"team_id" binding:"required" uri:"team_id" url:"team_id"`
}
