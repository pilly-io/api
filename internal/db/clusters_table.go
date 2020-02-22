package db

import (
	"github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"github.com/pilly-io/api/internal/models"
)

// ClustersTable implements Cluster specific operations
type ClustersTable struct {
	*BeegoTable
}

// NewClusterTable : returns a new Table object
func NewClusterTable(client orm.Ormer, model models.Cluster) *ClustersTable {
	table := NewBeegoTable(client, model)
	return &ClustersTable{table}
}

// Create a new cluster and populate missing fields
func (table *ClustersTable) Create(name string, provider string) (*models.Cluster, error) {
	uid, _ := uuid.NewRandom()
	cluster := &models.Cluster{
		Name:     name,
		Provider: provider,
		APIToken: uid.String(),
	}
	err := table.Insert(&cluster)
	return cluster, err
}
