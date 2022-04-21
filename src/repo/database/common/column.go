package common

import (
	"fmt"
)

type OrderType string

const (
	ASC  OrderType = "ASC"
	DESC OrderType = "DESC"
)

type IColumn interface {
	FullName() string
	Name() string
	TableName(tableName string) IColumn
	Order(orderType OrderType) IColumn
	Alias(alias string) IColumn
	Max() IColumn
	Min() IColumn
	setFormat(format string) IColumn
	Info() (
		name string,
		isOrderColumn bool,
	)
}

type Column struct {
	name      string
	orderType *OrderType
	alias     *string
	tableName *string
	format    *string
}

func NewColumn(name string) Column {
	return Column{
		name: name,
	}
}

func (c Column) FullName() string {
	if c.tableName == nil {
		return c.Name()
	}
	return fmt.Sprintf("%s.%s", *c.tableName, c.name)
}

func (c Column) Name() string {
	return c.name
}

func (c Column) TableName(tableName string) IColumn {
	c.tableName = &tableName
	return c
}

func (c Column) Order(orderType OrderType) IColumn {
	c.orderType = &orderType
	return c
}

func (c Column) Alias(alias string) IColumn {
	c.alias = &alias
	return c
}

func (c Column) Max() IColumn {
	return c.setFormat("MAX(%s)")
}

func (c Column) Min() IColumn {
	return c.setFormat("MIN(%s)")
}

func (c Column) setFormat(format string) IColumn {
	c.format = &format
	return c
}

func (c Column) Info() (
	name string,
	isOrderColumn bool,
) {
	if c.orderType != nil {
		isOrderColumn = true
		name = fmt.Sprintf("%s %s", c.name, string(*c.orderType))
		return
	}

	format := "%s"
	if c.format != nil {
		format = *c.format
	}
	name = fmt.Sprintf(format, c.name)
	if c.alias != nil {
		name += fmt.Sprintf(" AS %s", *c.alias)
	}

	return
}

type ColumnName string

func (c ColumnName) FullName() string {
	return string(c)
}

func (c ColumnName) Name() string {
	return string(c)
}

func (c ColumnName) TableName(tableName string) IColumn {
	return NewColumn(c.Name()).
		TableName(tableName)
}

func (c ColumnName) Order(orderType OrderType) IColumn {
	return NewColumn(c.Name()).
		Order(orderType)
}

func (c ColumnName) Alias(alias string) IColumn {
	return NewColumn(c.Name()).
		Alias(alias)
}

func (c ColumnName) Max() IColumn {
	return NewColumn(c.Name()).
		Max()
}

func (c ColumnName) Min() IColumn {
	return NewColumn(c.Name()).
		Min()
}

func (c ColumnName) setFormat(format string) IColumn {
	return NewColumn(c.Name()).
		setFormat(format)
}

func (c ColumnName) Info() (
	name string,
	isOrderColumn bool,
) {
	name = c.Name()
	isOrderColumn = false
	return
}
