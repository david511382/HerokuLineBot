package resp

type Role struct {
	Role uint
}

type NameLineID struct {
	Name   string
	LineID *string
}

type IDNameRole struct {
	Role
	ID   int
	Name string
}

type IDDepartment struct {
	Department string
	ID         int
}

type IDNameRoleDepartment struct {
	IDNameRole
	Department string
}

type NameRoleDepartmentLineIDCompanyID struct {
	NameLineID
	Role
	Department string
	CompanyID  *string
}
