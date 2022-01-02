package common

import (
	"gorm.io/gorm"
)

type ITransaction interface {
	Commit() error
	Rollback() error
}

type Transaction struct {
	*gorm.DB
}

func NewTransaction(conn *gorm.DB) *Transaction {
	return &Transaction{
		DB: conn,
	}
}

func (t *Transaction) Commit() error {
	return t.DB.Commit().Error
}

func (t *Transaction) Rollback() error {
	return t.DB.Rollback().Error
}
