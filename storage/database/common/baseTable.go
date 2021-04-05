package common

import (
	"fmt"
	"heroku-line-bot/storage/database/domain"
	"strings"

	"github.com/jinzhu/gorm"
)

type ITable interface {
	WhereArg(connection *gorm.DB, arg interface{}) *gorm.DB
	GetTable() interface{}
}

type BaseTable struct {
	BaseDatabase
	table ITable
}

func NewBaseTable(table ITable, writeDb, readDb *gorm.DB) *BaseTable {
	return &BaseTable{
		table: table,
		BaseDatabase: BaseDatabase{
			Write: writeDb,
			Read:  readDb,
		},
	}
}

func (t BaseTable) DbModel() *gorm.DB {
	table := t.table.GetTable()
	return t.Read.Model(table)
}

func (t BaseTable) Insert(trans *gorm.DB, datas ...interface{}) error {
	dp := t.Write
	if trans != nil {
		dp = trans
	}
	if cap(datas) == 0 || len(datas) == 0 {
		return domain.DB_NO_AFFECTED_Error
	} else if len(datas) > 1 {
		idatas := make([]interface{}, 0)
		for _, d := range datas {
			idatas = append(idatas, d)
		}

		count, err := batchInsert(dp, idatas)
		if err != nil {
			return err
		} else if count == 0 {
			return domain.DB_NO_AFFECTED_Error
		}
	} else {
		if err := dp.Create(datas[0]).Error; err != nil {
			return err
		}
	}
	return nil
}

func batchInsert(db *gorm.DB, data []interface{}) (int64, error) {
	// If there is no data, nothing to do.
	if len(data) == 0 {
		return 0, nil
	}

	table := data[0]
	mainScope := db.NewScope(table)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(data))
	for _, obj := range data {
		scope := db.NewScope(obj)
		fields := scope.Fields()
		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			var vars interface{}
			if (fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt") && fields[i].IsBlank {
				vars = gorm.NowFunc()
			} else {
				vars = fields[i].Field.Interface()
			}
			placeholders = append(placeholders, scope.AddToVars(vars))
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		// add real variables for the replacement of placeholders' '?' letter later.
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	//Execute and Log
	if err := mainScope.Exec().DB().Error; err != nil {
		return 0, err
	}

	return mainScope.DB().RowsAffected, nil
}

func (t BaseTable) Delete(trans *gorm.DB, arg interface{}) error {
	dp := t.Write
	if trans != nil {
		dp = trans
	}

	dp = t.table.WhereArg(dp, arg)

	table := t.table.GetTable()
	if err := dp.Delete(table).Error; err != nil {
		return err
	}

	return nil
}

func (t BaseTable) IsExist() bool {
	dp := t.Read
	table := t.table.GetTable()
	return dp.HasTable(table)
}

func (t BaseTable) CreateTable() error {
	dp := t.Read
	table := t.table.GetTable()
	if err := dp.CreateTable(table).Error; err != nil {
		return err
	}

	return nil
}
