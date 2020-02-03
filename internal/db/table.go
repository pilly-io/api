package db

import (
	"reflect"

	"github.com/jinzhu/gorm"
)

type Table interface {
	Find(query Query) (interface{}, error)
	FindAll(query Query) (*PaginatedCollection, error)
}

type GormTable struct {
	*gorm.DB
	kind interface{}
}

// NewTable : returns a new Table object
func NewTable(client *gorm.DB, kind interface{}) Table {
	return &GormTable{client, kind}
}

// Find first object that matches the conditions
func (table *GormTable) Find(query Query) (interface{}, error) {
	kind := reflect.TypeOf(table.kind)
	result := reflect.New(kind).Interface()
	err := table.Where(query.Conditions).First(&result).Error
	if err != nil {
		result = nil
	}
	return result, err
}

// FindAll returns all object matching parameters
func (table *GormTable) FindAll(query Query) (*PaginatedCollection, error) {
	count := 0
	kind := reflect.TypeOf(table.kind)

	results := reflect.New(reflect.SliceOf(kind)).Elem()
	//results := reflect. MakeSlice(kind, 0, 0)
	err := table.Where(query.Conditions).Find(&results).Error
	if err != nil {
		return nil, err
	}

	err = table.Model(&table.kind).Where(query.Conditions).Count(&count).Error
	if err != nil {
		return nil, err
	}
	return &PaginatedCollection{
		Objects:    results.Interface().([]interface{}),
		TotalCount: count,
	}, err
}
