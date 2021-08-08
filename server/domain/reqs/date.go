package reqs

import "time"

type FromTo struct {
	FromTime *time.Time `json:"from_time" form:"from_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"from_time" default:"2013-08-02T20:13:14+08:00"`
	ToTime   *time.Time `json:"to_time" form:"to_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"to_time" default:"2013-08-02T20:13:14+08:00"`
}

type FromToDate struct {
	FromDate *time.Time `json:"from_date" form:"from_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"from_date" default:"2013-08-02T00:00:00+08:00"`
	ToDate   *time.Time `json:"to_date" form:"to_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"to_date" default:"2013-08-02T00:00:00+08:00"`
}

type MustDate struct {
	Date time.Time `json:"date" form:"date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"date" default:"2013-08-02T00:00:00+08:00"`
}

type MustFromTo struct {
	FromTime time.Time `json:"from_time" form:"from_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"from_time" default:"2013-08-02T20:13:14+08:00"`
	ToTime   time.Time `json:"to_time" form:"to_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"to_time" default:"2013-08-02T20:13:14+08:00"`
}

type MustFromToDate struct {
	FromDate time.Time `json:"from_date" form:"from_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"from_date" default:"2013-08-02T00:00:00+08:00"`
	ToDate   time.Time `json:"to_date" form:"to_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"to_date" default:"2013-08-02T00:00:00+08:00"`
}

type MustFromBefore struct {
	FromTime   time.Time `json:"from_time" form:"from_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"from_time" default:"2013-08-02T20:13:14+08:00"`
	BeforeTime time.Time `json:"before_time" form:"before_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"before_time" default:"2013-08-02T20:13:14+08:00"`
}

type MustFromBeforeDate struct {
	FromDate   time.Time `json:"from_date" form:"from_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"from_date" default:"2013-08-02T00:00:00+08:00"`
	BeforeDate time.Time `json:"before_date" form:"before_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"required" url:"before_date" default:"2013-08-02T00:00:00+08:00"`
}

type FromBefore struct {
	FromTime   *time.Time `json:"from_time" form:"from_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"from_time" default:"2013-08-02T20:13:14+08:00"`
	BeforeTime *time.Time `json:"before_time" form:"before_time" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"before_time" default:"2013-08-02T20:13:14+08:00"`
}

type FromBeforeDate struct {
	FromDate   *time.Time `json:"from_date" form:"from_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"from_date" default:"2013-08-02T00:00:00+08:00"`
	BeforeDate *time.Time `json:"before_date" form:"before_date" time_format:"2006-01-02T15:04:05Z07:00" binding:"-" url:"before_date" default:"2013-08-02T00:00:00+08:00"`
}
