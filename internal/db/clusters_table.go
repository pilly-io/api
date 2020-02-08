package db

import (
	"github.com/google/uuid"
	"github.com/pilly-io/api/internal/models"
)

// ClustersTable implements Cluster specific operations
type ClustersTable struct {
	*GormTable
}

// Create a new cluster and populate missing fields
func (table *ClustersTable) Create(name string) (*models.Cluster, error) {
	uid, _ := uuid.NewRandom()
	cluster := &models.Cluster{
		Name:     name,
		APIToken: uid.String(),
	}
	err := table.Insert(&cluster)
	return cluster, err
}
