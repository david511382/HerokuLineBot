package common

type ErrorBaseTable struct {
	err error
}

func NewErrorBaseTable(err error) *ErrorBaseTable {
	result := &ErrorBaseTable{
		err: err,
	}
	return result
}

// response: pointer of slice / struct
func (t ErrorBaseTable) SelectColumns(arg interface{}, response interface{}, columns ...string) error {
	return t.err
}

func (t ErrorBaseTable) Count(arg interface{}) (int64, error) {
	return 0, t.err
}

func (t ErrorBaseTable) Insert(datas interface{}) error {
	return t.err
}

func (t ErrorBaseTable) MigrationTable() error {
	return t.err
}

func (t ErrorBaseTable) MigrationData(length int, datas interface{}) error {
	return t.err
}

func (t ErrorBaseTable) Delete(arg interface{}) error {
	return t.err
}

func (t ErrorBaseTable) Update(arg interface{}, fields map[string]interface{}) error {
	return t.err
}

func (t ErrorBaseTable) IsExist() bool {
	return false
}

func (t ErrorBaseTable) CreateTable() error {
	return t.err
}
