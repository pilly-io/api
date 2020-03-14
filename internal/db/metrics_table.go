package db

import (
	"fmt"

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
func (table *MetricsTable) FindAll(clusterID uint, ownerUIDs []string, period uint, queryInterval QueryInterval) (*[]*models.Metric, error) {
	var results []*models.Metric
	if len(ownerUIDs) == 0 {
		return nil, nil
	}
	ownerUIDsMarks := ""
	for index := range ownerUIDs {
		if index == 0 {
			ownerUIDsMarks += "?"
		} else {
			ownerUIDsMarks += " , ?"
		}
	}

	query := fmt.Sprintf(`
	SELECT AVG(value) as value, name, owner_uid, cluster_id,
	to_timestamp(floor((extract('epoch' from created_at) / ? )) * ?)  AT TIME ZONE 'UTC' as period
	FROM metric
	WHERE cluster_id = ?
	AND owner_uid IN(%s)
	AND created_at <= ?
	AND (deleted_at IS NULL or deleted_at >= ?)
	GROUP BY name, period, owner_uid, cluster_id
	`, ownerUIDsMarks)

	_, err := table.client.Raw(query, period, period, clusterID, ownerUIDs, queryInterval.End, queryInterval.Start).QueryRows(&results)
	return &results, err
}
