package conn

import "github.com/jinzhu/gorm"

type IConnect interface {
	Connect() (*gorm.DB, error)
}
