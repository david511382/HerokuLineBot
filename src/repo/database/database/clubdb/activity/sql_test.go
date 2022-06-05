package activity

import (
	dbModel "heroku-line-bot/src/model/database"
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/common"
	"testing"
	"time"
)

func TestActivity_Select(t *testing.T) {
	t.Parallel()

	type args struct {
		arg     Reqs
		columns []common.IColumn
	}
	type migrations struct {
		table []*Model
	}
	type wants struct {
		data []*Model
	}
	tests := []struct {
		name       string
		args       args
		migrations migrations
		wants      wants
	}{
		{
			"order",
			args{
				arg: Reqs{
					Date: dbModel.Date{
						FromDate: util.GetUTCTimeP(2013, 8, 1, 16),
					},
				},
				columns: []common.IColumn{
					COLUMN_Date,
					COLUMN_PlaceID.Order(common.DESC),
					COLUMN_ID.Order(common.DESC),
					COLUMN_Date.Order(common.DESC),
				},
			},
			migrations{
				table: []*Model{
					{
						ID:      5,
						Date:    util.GetUTCTime(2013, 8, 2),
						PlaceID: 1,
					},
					{
						ID:   8,
						Date: util.GetUTCTime(2013, 8, 1),
					},
					{
						ID:      2,
						Date:    util.GetUTCTime(2013, 8, 3),
						PlaceID: 1,
					},
				},
			},
			wants{
				data: []*Model{
					{
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
					{
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDb(t)

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
	t.Parallel()

	type args struct {
		arg Reqs
	}
	type migrations struct {
		table []*Model
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
				arg: Reqs{},
			},
			migrations{
				table: []*Model{
					{
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 5),
					},
					{
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
					},
					{
						Date: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
					},
				},
			},
			wants{
				maxDate: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 8),
				minDate: *util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDb(t)

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
