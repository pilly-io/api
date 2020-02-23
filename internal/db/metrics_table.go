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
func (table *MetricsTable) FindAll(clusterID uint, namespace string, ownerUIDs []string, period uint, queryInterval QueryInterval) (*[]models.Metric, error) {
	var results []models.Metric
	_, err := table.client.Raw(`
	SELECT AVG(metric.value) as value, metric.name, metric.owner_uid, metric.cluster_id,
	to_timestamp(floor((extract('epoch' from metric.created_at) / ? )) * ?)  AT TIME ZONE 'UTC' as period
	FROM metric
	LEFT JOIN owner (ON owner.uid = metric.owner_uid)
	WHERE metric.cluster_id = ?
	AND metric.created_at <= ?
	AND (metric.deleted_at IS NULL or metric.deleted_at >= ?)
	AND ( ? = '' OR owner.namespace = ?)
	AND ( COALESCE( ? ) IS NULL OR metric.owner_uid IN ( ? ))
	GROUP BY metric.name, metric.period, metric.owner_uid, metric.cluster_id
	`, period, period, clusterID, queryInterval.End, queryInterval.Start, namespace, namespace, ownerUIDs, ownerUIDs).QueryRows(&results)
	return &results, err
}
