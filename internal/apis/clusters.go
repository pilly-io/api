package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	log "github.com/sirupsen/logrus"
)

// Clusters endpoints
type ClustersHandler struct {
	db db.Database
}

// FindAll get all clusters
func (handler *ClustersHandler) FindAll(c *gin.Context) {
	clusters := []*models.Cluster{}
	query := db.Query{
		Result: &clusters,
	}
	err := handler.db.Clusters().FindAll(query)

	if err != nil {
		log.Errorf("Cannot retrieve the clusters: %s", err)
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, clusters)
	}
}
