package activity

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"testing"
	"time"
)

func TestActivity_Select(t *testing.T) {
	type args struct {
		arg     dbModel.ReqsClubActivity
		columns []Column
	}
	type migrations struct {
		table []*dbModel.ClubActivity
	}
	type wants struct {
		data []*dbModel.ClubActivity
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"array condition",
			args{
				arg: dbModel.ReqsClubActivity{
					Date: dbModel.Date{
						FromDate: util.GetUTCTimeP(2013, 8, 2),
					},
				},
				columns: []Column{
					COLUMN_Date,
				},
			},
			migrations{
				table: []*dbModel.ClubActivity{
					{
						ID:   5,
						Date: util.GetUTCTime(2013, 8, 2),
					},
					{
						ID:   8,
						Date: util.GetUTCTime(2013, 8, 1),
					},
					{
						ID:   2,
						Date: util.GetUTCTime(2013, 8, 3),
					},
				},
			},
			wants{
				data: []*dbModel.ClubActivity{
					{
						Date: *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
					{
						Date: *util.GetTimePLoc(global.Location, 2013, 8, 3),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.MigrationData(tt.migrations.table...); err != nil {
				t.Fatal(err.Error())
			}

			got, err := db.Select(tt.args.arg, tt.args.columns...)
			if err != nil {
				t.Fatal(err)
			}

			if ok, msg := util.Comp(got, tt.wants.data); !ok {
				t.Error(msg)
				return
			}
		})
	}
}

func TestActivity_MinMaxID(t *testing.T) {
	type args struct {
		arg dbModel.ReqsClubActivity
	}
	type migrations struct {
		table []*dbModel.ClubActivity
	}
	type wants struct {
		maxDate time.Time
		minDate time.Time
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"min max date",
			args{
				arg: dbModel.ReqsClubActivity{},
			},
			migrations{
				table: []*dbModel.ClubActivity{
					{
						Date: *util.GetTimePLoc(global.Location, 2013, 8, 5),
					},
					{
						Date: *util.GetTimePLoc(global.Location, 2013, 8, 8),
					},
					{
						Date: *util.GetTimePLoc(global.Location, 2013, 8, 2),
					},
				},
			},
			wants{
				maxDate: *util.GetTimePLoc(global.Location, 2013, 8, 8),
				minDate: *util.GetTimePLoc(global.Location, 2013, 8, 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.MigrationData(tt.migrations.table...); err != nil {
				t.Fatal(err.Error())
			}

			gotMaxDate, gotMinDate, err := db.MinMaxDate(tt.args.arg)
			if err != nil {
				t.Fatal(err)
			}

			if ok, msg := util.Comp(gotMaxDate, tt.wants.maxDate); !ok {
				t.Error(msg)
				return
			}
			if ok, msg := util.Comp(gotMinDate, tt.wants.minDate); !ok {
				t.Error(msg)
				return
			}
		})
	}
}
