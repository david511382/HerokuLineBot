package resp

import "time"

type Role struct {
	Role uint
}

type LineID struct {
	LineID *string
}

type NameLineID struct {
	Name   string
	LineID *string
}

type IDNameLineID struct {
	NameLineID
	ID int
}

type IDName struct {
	ID   int
	Name string
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

type IDNameDepartmentJoinDate struct {
	IDName
	Department string
	JoinDate   *time.Time
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
