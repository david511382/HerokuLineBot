package club

import (
	"heroku-line-bot/src/logic/club/domain"
	"strings"
)

const DEPARTMENT_NONE = "無"

type Department string

func NewDepartment(處 domain.Department, 部, 組 string) Department {
	result := Department("")
	result.set(處, 部, 組)
	return result
}

func (d Department) Split() (處 domain.Department, 部, 組 string) {
	ds := strings.Split(string(d), "/")
	if len(ds) >= 1 {
		處 = domain.Department(ds[0])
	}
	if len(ds) >= 2 {
		部 = ds[1]
	}
	if len(ds) >= 3 {
		組 = ds[2]
	}
	return
}

func (d Department) IsClubMember() bool {
	處, _, _ := d.Split()
	for _, clubMemberDepartment := range domain.ClubMemberDepartments {
		if 處 == clubMemberDepartment {
			return true
		}
	}
	return false
}

func (d *Department) Set處(data domain.Department) {
	_, 部, 組 := d.Split()
	d.set(data, 部, 組)
}

func (d *Department) Set部(data string) {
	處, _, 組 := d.Split()
	d.set(處, data, 組)
}

func (d *Department) Set組(data string) {
	處, 部, _ := d.Split()
	d.set(處, 部, data)
}

func (d *Department) set(處 domain.Department, 部, 組 string) {
	strs := []string{
		string(處), 部, 組,
	}
	*d = Department(strings.Join(strs, "/"))
}
