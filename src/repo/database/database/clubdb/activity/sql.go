package activity

import (
	"fmt"
	dbModel "heroku-line-bot/src/model/database"
	"time"

	"gorm.io/gorm"
)

func (t Activity) Insert(datas ...*dbModel.ClubActivity) error {
	return t.IBaseTable.Insert(datas)
}

func (t Activity) MigrationData(datas ...*dbModel.ClubActivity) error {
	return t.IBaseTable.MigrationData(len(datas), datas)
}

func (t Activity) Delete(arg dbModel.ReqsClubActivity) error {
	return t.IBaseTable.Delete(arg)
}

func (t Activity) Update(trans *gorm.DB, arg dbModel.ReqsClubActivityUpdate) error {
	fields := make(map[string]interface{})
	if p := arg.LogisticID; p != nil {
		fields[string(COLUMN_LogisticID)] = *p
	}
	if p := arg.MemberCount; p != nil {
		fields[string(COLUMN_MemberCount)] = *p
	}
	if p := arg.GuestCount; p != nil {
		fields[string(COLUMN_GuestCount)] = *p
	}
	if p := arg.MemberFee; p != nil {
		fields[string(COLUMN_MemberFee)] = *p
	}
	if p := arg.GuestFee; p != nil {
		fields[string(COLUMN_GuestFee)] = *p
	}

	return t.IBaseTable.Update(arg.ReqsClubActivity, fields)
}

func (t Activity) Select(arg dbModel.ReqsClubActivity, columns ...Column) ([]*dbModel.ClubActivity, error) {
	result := make([]*dbModel.ClubActivity, 0)

	columnStrs := make([]string, 0)
	for _, column := range columns {
		columnStrs = append(columnStrs, string(column))
	}
	if len(columnStrs) == 0 {
		columnStrs = append(columnStrs, "*")
	}

	if err := t.SelectColumns(arg, &result, columnStrs...); err != nil {
		return nil, err
	}

	return result, nil
}

func (t Activity) MinMaxDate(arg dbModel.ReqsClubActivity) (maxDate, minDate time.Time, resultErr error) {
	type MinMaxID struct {
		MaxDate time.Time
		MinDate time.Time
	}

	response := &MinMaxID{}
	if err := t.SelectColumns(
		arg, response,
		fmt.Sprintf("MIN(%s) AS min_date", COLUMN_Date),
		fmt.Sprintf("MAX(%s) AS max_date", COLUMN_Date),
	); err != nil {
		resultErr = err
		return
	}

	maxDate = response.MaxDate
	minDate = response.MinDate
	return
}
