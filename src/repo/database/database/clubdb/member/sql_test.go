package member

import (
	"heroku-line-bot/src/pkg/global"
	"heroku-line-bot/src/pkg/util"
	"heroku-line-bot/src/repo/database/common"
	"testing"
)

func TestMember_Select(t *testing.T) {
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
			"join date is null",
			args{
				arg: Reqs{
					JoinDateIsNull: util.GetBoolP(true),
				},
				columns: []common.IColumn{
					COLUMN_ID,
				},
			},
			migrations{
				table: []*Model{
					{
						ID:       1,
						JoinDate: util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
					{
						ID: 2,
					},
				},
			},
			wants{
				data: []*Model{
					{
						ID: 2,
					},
				},
			},
		},
		{
			"is delete",
			args{
				arg: Reqs{
					IsDelete: util.GetBoolP(true),
				},
				columns: []common.IColumn{
					COLUMN_ID,
				},
			},
			migrations{
				table: []*Model{
					{
						ID: 1,
					},
					{
						ID:        2,
						DeletedAt: util.GetTimePLoc(global.TimeUtilObj.GetLocation(), 2013, 8, 3),
					},
				},
			},
			wants{
				data: []*Model{
					{
						ID: 2,
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
