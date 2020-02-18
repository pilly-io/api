package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pilly-io/api/internal/models"
)

// OwnersTable implements Cluster specific operations
type OwnersTable struct {
	*GormTable
}

// NewOwnersTable : returns a new Table object
func NewOwnersTable(client *gorm.DB, model models.Owner) *OwnersTable {
	table := GormTable{client, model}
	return &OwnersTable{&table}
}

//ComputeResources : Given metrics compute the owners resources
func (table *OwnersTable) ComputeResources(owners *[]models.Owner, metrics *models.IndexedMetrics) {
	for i, owner := range *owners {
		if tmetric, exist := (*metrics)[owner.UID]; exist {
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
				(*owners)[i].Metrics = append((*owners)[i].Metrics, resource)
			}
		}
	}
}
