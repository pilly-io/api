package services

import (
	"github.com/pilly-io/api/internal/models"
)

type clusterDAO interface {
	GetByName(name string) (*models.Cluster, error)
	GetByID(id uint) (*models.Cluster, error)
}

type ClusterService struct {
	dao clusterDAO
}

func NewClusterService(dao clusterDAO) *ClusterService {
	return &ClusterService{dao}
}

func (svc *ClusterService) GetByName(name string) (*models.Cluster, error) {
	return svc.dao.GetByName(name)
}

func (svc *ClusterService) GetByID(id uint) (*models.Cluster, error) {
	return svc.dao.GetByID(id)
}
