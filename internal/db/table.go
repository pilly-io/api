package db

import (
	"fmt"
	"math"

	"github.com/jinzhu/gorm"
)

type Table interface {
	Find(query Query, result interface{}) error
	FindAll(query Query, results interface{}) (*PaginationInfo, error)
	Exists(query Query) bool
	Insert(value interface{}) error
	//Create(value interface{}) error
}

type GormTable struct {
	*gorm.DB
	kind interface{}
}

// NewTable : returns a new Table object
func NewTable(client *gorm.DB, kind interface{}) Table {
	return &GormTable{client, kind}
}

// Insert creates a new record in the right table
func (table *GormTable) Insert(value interface{}) error {
	return table.Create(value).Error
}

// Find first object that matches the conditions
func (table *GormTable) Find(query Query, result interface{}) error {
	return table.Where(query.Conditions).First(result).Error
}

// Exists check if at least one record exist for this query
func (table *GormTable) Exists(query Query) bool {
	count := 0
	table.Model(table.kind).Where(query.Conditions).Limit(1).Count(&count)
	return count > 0
}

// FindAll returns all object matching parameters
func (table *GormTable) FindAll(query Query, results interface{}) (*PaginationInfo, error) {
	count := 0

	builder := table.Where(query.Conditions)

	if query.Interval != nil {
		fmt.Println(query.Interval.Start)
		fmt.Println(query.Interval.End)
		builder = builder.Unscoped()
		builder = builder.Where("created_at <= ?", query.Interval.End)
		builder = builder.Where("deleted_at IS NULL OR deleted_at >= ?", query.Interval.Start)
	}

	if query.Limit > 0 {
		builder = builder.Limit(query.Limit)
	}

	if query.OrderBy != "" {
		direction := "ASC"
		if query.Desc {
			direction = "DESC"
		}
		orderBy := query.OrderBy + " " + direction
		builder = builder.Order(orderBy)
	}

	if query.Page > 1 {
		builder = builder.Offset((query.Page - 1) * query.Limit)
	}

	err := builder.Find(results).Error
	if err != nil {
		return nil, err
	}

	countBuilder := builder.Limit(nil).Order(nil).Offset(nil)
	err = countBuilder.Model(results).Count(&count).Error
	if err != nil {
		return nil, err
	}

	maxPage := 0
	if query.Limit > 0 {
		maxPage = int(math.Ceil(float64(count) / float64(query.Limit)))
	}
	return &PaginationInfo{
		TotalCount:  count,
		Limit:       query.Limit,
		MaxPage:     maxPage,
		CurrentPage: query.Page,
	}, err
}
