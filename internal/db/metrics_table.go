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
func (table *MetricsTable) FindAll(clusterID uint, refType string, refUIDs []string, period uint, queryInterval QueryInterval) (*[]*models.Metric, error) {
	refColumns := map[string]string{
		"owner":     "owner_uid",
		"namespace": "namespace_uid",
	}
	var results []*models.Metric
	if len(refUIDs) == 0 {
		return nil, nil
	}
	UIDsMarks := ""
	for index := range refUIDs {
		if index == 0 {
			UIDsMarks += "?"
		} else {
			UIDsMarks += " , ?"
		}
	}

	query := fmt.Sprintf(`
	SELECT AVG(value) as value, name, %s, cluster_id,
	to_timestamp(floor((extract('epoch' from created_at) / ? )) * ?)  AT TIME ZONE 'UTC' as period
	FROM metric
	WHERE cluster_id = ?
	AND %s IN(%s)
	AND created_at <= ?
	AND (deleted_at IS NULL or deleted_at >= ?)
	GROUP BY name, period, %s, cluster_id
	`, refColumns[refType], refColumns[refType], UIDsMarks, refColumns[refType])

	_, err := table.client.Raw(query, period, period, clusterID, refUIDs, queryInterval.End, queryInterval.Start).QueryRows(&results)
	return &results, err
}
