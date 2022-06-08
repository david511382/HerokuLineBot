package database

import (
	"time"
)

type Date struct {
	Date       *time.Time
	Dates      []*time.Time
	FromDate   *time.Time
	AfterDate  *time.Time
	ToDate     *time.Time
	BeforeDate *time.Time
}
