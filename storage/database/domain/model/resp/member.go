package resp

type Role struct {
	Role uint
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
