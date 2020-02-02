package daos

import "github.com/pilly-io/api/internal/models"
import . "github.com/pilly-io/api/internal/config"

type ClusterDAO struct{}

func NewClusterDao() *ClusterDAO {
	return &ClusterDAO{}
}

func (dao *ClusterDAO) GetByName(name string) (*models.Cluster, error) {
	cluster := models.Cluster{}
	err := Settings.DB.Where("name = ?", name).First(&cluster).Error
	return &cluster, err

}

func (dao *ClusterDAO) GetByID(id uint) (*models.Cluster, error) {
	cluster := models.Cluster{}
	err := Settings.DB.Where("id = ?", id).First(&cluster).Error
	return &cluster, err
}
