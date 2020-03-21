package collector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// MetricsHandler definition
type MetricsHandler struct {
	DB db.Database
}

// Create insert metrics for current cluster
func (handler *MetricsHandler) Create(c *gin.Context) {
	var metrics []*models.Metric
	c.BindJSON(&metrics)

	cluster := c.MustGet("cluster").(*models.Cluster)

	for _, metric := range metrics {
		metric.ClusterID = cluster.ID
	}
	handler.DB.Metrics().BulkInsert(metrics)

	c.JSON(http.StatusCreated, utils.ObjectToJSON(nil))
}
