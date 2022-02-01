package conn

import "gorm.io/gorm"

type IConnect interface {
	GetDialector() gorm.Dialector
}
