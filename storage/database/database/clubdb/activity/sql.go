package activity

import (
	dbModel "heroku-line-bot/model/database"

	"gorm.io/gorm"
)

func (t Activity) Insert(datas ...*dbModel.ClubActivity) error {
	return t.BaseTable.Insert(datas)
}

func (t Activity) MigrationData(datas ...*dbModel.ClubActivity) error {
	return t.BaseTable.MigrationData(len(datas), datas)
}

func (t Activity) Delete(arg dbModel.ReqsClubActivity) error {
	return t.BaseTable.Delete(arg)
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

	return t.BaseTable.Update(arg.ReqsClubActivity, fields)
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
