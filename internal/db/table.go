package db

import "github.com/jinzhu/gorm"

type Table interface {
	Find(Query) error
	FindAll(Query) error
}

type GormTable struct {
	*gorm.DB
}

// NewTable : returns a new Table object
func NewTable(client *gorm.DB) Table {
	return &GormTable{client}
}

// Find first object that matches the conditions
func (table *GormTable) Find(query Query) error {
	return table.Where(query.Conditions).First(query.Result).Error
}

// FindAll returns all object matching parameters
func (table *GormTable) FindAll(query Query) error {
	return table.Where(query.Conditions).Find(query.Result).Error
}
