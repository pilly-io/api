package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pilly-io/api/internal/models"
)

// MetricsTable implements Cluster specific operations
type MetricsTable struct {
	*GormTable
}

// NewMetricsTable : returns a new Table object
func NewMetricsTable(client *gorm.DB, model models.Metric) *MetricsTable {
	table := GormTable{client, model}
	return &MetricsTable{&table}
}

// FindAll returns all object matching parameters
func (table *MetricsTable) FindAll(query Query, results interface{}) error {
	builder := table.Where(query.Conditions)

	if query.Interval != nil {
		builder = builder.Unscoped()
		builder = builder.Where("created_at <= ?", query.Interval.End)
		builder = builder.Where("deleted_at IS NULL OR deleted_at >= ?", query.Interval.Start)
	}
	return builder.Find(results).Error
}
