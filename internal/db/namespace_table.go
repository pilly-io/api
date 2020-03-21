package db

import (
	"github.com/astaxie/beego/orm"
	"github.com/pilly-io/api/internal/models"
)

// OwnersTable implements Cluster specific operations
type NamespacesTable struct {
	*BeegoTable
}

// NewOwnersTable : returns a new Table object
func NewNamespacesTable(client orm.Ormer, model models.Namespace) *NamespacesTable {
	table := NewBeegoTable(client, model)
	return &NamespacesTable{table}
}

//ComputeResources : Given metrics compute the owners resources
func (table *NamespacesTable) ComputeResources(objects *[]*models.Namespace, metrics *models.IndexedMetrics) {
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
