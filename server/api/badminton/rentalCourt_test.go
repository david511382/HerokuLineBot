package badminton

import (
	"heroku-line-bot/global"
	"heroku-line-bot/util"
	"testing"
	"time"
)

func Test_addRentalCourtGetRentalDates(t *testing.T) {
	type args struct {
		fromDate     util.DateTime
		toDate       util.DateTime
		everyWeekday *int
		excludeDates []*time.Time
	}
	tests := []struct {
		name            string
		args            args
		wantRentalDates []util.DateTime
	}{
		{
			"hour exclude date",
			args{
				fromDate:     *util.NewDateTimeP(global.Location, 2013, 8, 1),
				toDate:       *util.NewDateTimeP(global.Location, 2013, 8, 3),
				everyWeekday: nil,
				excludeDates: []*time.Time{
					util.GetTimePLoc(global.Location, 2013, 8, 1, 23),
					util.GetTimePLoc(global.Location, 2013, 8, 3),
				},
			},
			[]util.DateTime{
				*util.NewDateTimeP(global.Location, 2013, 8, 2),
			},
		},
		{
			"everyweekdate exclude date",
			args{
				fromDate:     *util.NewDateTimeP(global.Location, 2013, 8, 2),
				toDate:       *util.NewDateTimeP(global.Location, 2013, 8, 9),
				everyWeekday: util.GetIntP(5),
				excludeDates: []*time.Time{
					util.GetTimePLoc(global.Location, 2013, 8, 2),
				},
			},
			[]util.DateTime{
				*util.NewDateTimeP(global.Location, 2013, 8, 9),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRentalDates := addRentalCourtGetRentalDates(tt.args.fromDate, tt.args.toDate, tt.args.everyWeekday, tt.args.excludeDates)
			if ok, msg := util.Comp(gotRentalDates, tt.wantRentalDates); !ok {
				t.Fatal(msg)
			}
		})
	}
}
