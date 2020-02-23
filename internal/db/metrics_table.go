package db

import (
	"github.com/astaxie/beego/orm"
	"github.com/pilly-io/api/internal/models"
)

// MetricsTable implements Cluster specific operations
type MetricsTable struct {
	*BeegoTable
}

// NewMetricsTable : returns a new Table object
func NewMetricsTable(client orm.Ormer, model models.Metric) *MetricsTable {
	table := NewBeegoTable(client, model)
	return &MetricsTable{table}
}

// FindAll returns all the metrics using an average on a specific period
func (table *MetricsTable) FindAll(clusterID uint, period uint, queryInterval QueryInterval) (*[]models.Metric, error) {
	var results []models.Metric
	_, err := table.client.Raw(`
	SELECT AVG(value) as value, name, owner_uid, cluster_id,
	to_timestamp(floor((extract('epoch' from created_at) / ? )) * ?)  AT TIME ZONE 'UTC' as period
	FROM metrics
	WHERE cluster_id = ?
	AND created_at <= ?
	AND (deleted_at IS NULL or deleted_at >= ?)
	GROUP BY name, period, owner_uid, cluster_id
	`, period, period, clusterID, queryInterval.End, queryInterval.Start).QueryRows(&results)
	return &results, err
}
