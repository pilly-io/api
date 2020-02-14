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
