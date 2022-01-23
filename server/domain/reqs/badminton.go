package reqs

import "time"

type GetRentalCourts struct {
	MustFromToDate
	TeamID int `json:"team_id" form:"team_id" binding:"-" uri:"team_id" url:"team_id"`
}

type AddRentalCourt struct {
	MustFromToDate
	PlaceID      int          `json:"place_id" form:"place_id" binding:"required" uri:"place_id" url:"place_id"`
	TeamID       int          `json:"team_id" form:"team_id" binding:"required" uri:"team_id" url:"team_id"`
	EveryWeekday *int         `json:"every_weekday" form:"every_weekday" binding:"-" uri:"every_weekday" url:"every_weekday"`
	ExcludeDates []*time.Time `json:"exclude_dates" form:"exclude_dates" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"exclude_dates"`

	DespositMoney *int       `json:"desposit_money" form:"desposit_money" binding:"-" uri:"desposit_money" url:"desposit_money"`
	DespositDate  *time.Time `json:"desposit_date" form:"desposit_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"desposit_date"`
	BalanceMoney  *int       `json:"balance_money" form:"balance_money" binding:"-" uri:"balance_money" url:"balance_money"`
	BalanceDate   *time.Time `json:"balance_date" form:"balance_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"balance_date"`

	CourtFromTime time.Time `json:"court_from_time" form:"court_from_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"court_from_time"`
	CourtToTime   time.Time `json:"court_to_time" form:"court_to_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"court_to_time"`
	CourtCount    int       `json:"court_count" form:"court_count" binding:"required" uri:"court_count" url:"court_count"`
	PricePerHour  int       `json:"price_per_hour" form:"price_per_hour" binding:"required" uri:"price_per_hour" url:"price_per_hour"`
}

type GetActivitys struct {
	FromToDate
	Page
	PlaceIDs      []int `json:"place_ids" form:"place_ids" binding:"-" uri:"place_ids" url:"place_ids"`
	TeamIDs       []int `json:"team_ids" form:"team_ids" binding:"-" uri:"team_ids" url:"team_ids"`
	EveryWeekdays []int `json:"every_weekdays" form:"every_weekdays" binding:"-" uri:"every_weekdays" url:"every_weekdays"`
}
