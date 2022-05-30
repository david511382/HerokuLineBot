package util

import (
	"testing"
	"time"
)

func TestTimeInt(t *testing.T) {
	t1 := *GetTimePLoc(time.Local, 2001, 2, 3, 4, 5, 6, 7)

	sec := Second().Of(t1)
	min := Minute().Of(t1)
	{
		secInt := sec.Int()
		minInt := secInt.Minute()
		if ok, msg := Comp(min, minInt.Time(time.Local)); !ok {
			t.Error(msg)
			return
		}
	}

	hour := Hour().Of(t1)
	{
		minInt := min.Int()
		hourInt := minInt.Hour()
		if ok, msg := Comp(hour, hourInt.Time(time.Local)); !ok {
			t.Error(msg)
			return
		}
	}

	date := Date().Of(t1)
	{
		hourInt := hour.Int()
		dateInt := hourInt.Date()
		if ok, msg := Comp(date, dateInt.Time(time.Local)); !ok {
			t.Error(msg)
			return
		}
	}

	month := Month().Of(t1)
	{
		dateInt := date.Int()
		monthInt := dateInt.Month()
		if ok, msg := Comp(month, monthInt.Time(time.Local)); !ok {
			t.Error(msg)
			return
		}
	}

	year := Year().Of(t1)
	{
		monthInt := month.Int()
		yearInt := monthInt.Year()
		if ok, msg := Comp(year, yearInt.Time(time.Local)); !ok {
			t.Error(msg)
			return
		}
	}

	{
		week := Week().Of(t1)
		want := *GetTimePLoc(time.Local, 2001, 1, 28, 4, 5, 6, 7)
		if ok, msg := Comp(week.Time(), want); !ok {
			t.Error(msg)
			return
		}
	}
}

func TestZeroTimeInt(t *testing.T) {
	t1 := time.Time{}

	sec := Second().Of(t1)
	min := Minute().Of(t1)
	{
		secInt := sec.Int()
		minInt := secInt.Minute()
		if ok, msg := Comp(min, minInt.Time(time.UTC)); !ok {
			t.Error(msg)
			return
		}
	}

	hour := Hour().Of(t1)
	{
		minInt := min.Int()
		hourInt := minInt.Hour()
		if ok, msg := Comp(hour, hourInt.Time(time.UTC)); !ok {
			t.Error(msg)
			return
		}
	}

	date := Date().Of(t1)
	{
		hourInt := hour.Int()
		dateInt := hourInt.Date()
		if ok, msg := Comp(date, dateInt.Time(time.UTC)); !ok {
			t.Error(msg)
			return
		}
	}

	month := Month().Of(t1)
	{
		dateInt := date.Int()
		monthInt := dateInt.Month()
		if ok, msg := Comp(month, monthInt.Time(time.UTC)); !ok {
			t.Error(msg)
			return
		}
	}

	year := Year().Of(t1)
	{
		monthInt := month.Int()
		yearInt := monthInt.Year()
		if ok, msg := Comp(year, yearInt.Time(time.UTC)); !ok {
			t.Error(msg)
			return
		}
	}
}

func TestTimeNext(t *testing.T) {
	t1 := *GetTimePLoc(time.Local, 2001, 1, 25, 4, 5, 6, 7)
	week := Week().Of(t1)
	week = week.Next(2)
	want := *GetTimePLoc(time.Local, 2001, 2, 4, 4, 5, 6, 7)
	if ok, msg := Comp(week.Time(), want); !ok {
		t.Error(msg)
		return
	}
}
