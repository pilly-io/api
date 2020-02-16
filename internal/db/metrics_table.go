package db

import (
	"fmt"

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

// FindAll returns all the metrics using an average on a specific perio
func (table *MetricsTable) FindAll(period uint, ownerIDs []string) (*[]models.Metric, error) {
	var results []models.Metric
	/*builder := table.Where(query.Conditions)
	builder = table.Unscoped()
	builder = builder.Where("created_at <= ?", query.Interval.End)
	builder = builder.Where("deleted_at IS NULL OR deleted_at >= ?", query.Interval.Start)
	builder = builder.Select("AVG(value)").Select("to_timestamp(floor((extract('epoch' from created_at) / 180 )) * 180)  AT TIME ZONE 'UTC' as period")
	return builder.Find(results).Error*/
	fmt.Println(period)
	table.Raw(`
	SELECT AVG(value) as value, name, owner_uid, cluster_id,
	to_timestamp(floor((extract('epoch' from created_at) / ? )) * ?)  AT TIME ZONE 'UTC' as period
	FROM metrics
	WHERE owner_uid IN (?)
	GROUP BY name, period, owner_uid, cluster_id
	`, period, period, ownerIDs).Scan(&results)
	return &results, nil
}
