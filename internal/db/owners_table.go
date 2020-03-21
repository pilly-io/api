package db

import (
	"github.com/astaxie/beego/orm"
	"github.com/pilly-io/api/internal/models"
)

// OwnersTable implements Cluster specific operations
type OwnersTable struct {
	*BeegoTable
}

// NewOwnersTable : returns a new Table object
func NewOwnersTable(client orm.Ormer, model models.Owner) *OwnersTable {
	table := NewBeegoTable(client, model)
	return &OwnersTable{table}
}

//ComputeResources : Given metrics compute the owners resources
func (table *OwnersTable) ComputeResources(objects *[]*models.Owner, metrics *models.IndexedMetrics) {
	for i, object := range *objects {
		if tmetric, exist := (*metrics)[object.UID]; exist {
			for timestamp, metric := range tmetric {
				resource := models.Resources{}
				resource.ResourcesTimestamp = timestamp
				resource.ResourcesUsed = map[string]interface{}{
					"cpu":    metric[models.MetricCPUUsed],
					"memory": metric[models.MetricMemoryUsed],
				}
				resource.ResourcesRequested = map[string]interface{}{
					"cpu":    metric[models.MetricCPURequested],
					"memory": metric[models.MetricMemoryRequested],
				}
				(*objects)[i].Resources = append((*objects)[i].Resources, resource)
			}
		}
	}
}
