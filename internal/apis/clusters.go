package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// ClustersHandler definition
type ClustersHandler struct {
	DB db.Database
}

// Create a cluster if not exists
func (handler *ClustersHandler) Create(c *gin.Context) {
	cluster := models.Cluster{}
	c.BindJSON(&cluster)
	query := db.Query{
		Conditions: db.QueryConditions{"name": cluster.Name},
	}

	if !handler.DB.Clusters().Exists(query) {
		_, err := handler.DB.Clusters().Create(cluster.Name, cluster.Provider)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, ErrorsToJSON("error"))
		} else {
			c.JSON(http.StatusCreated, ObjectToJSON(&cluster))
		}
	} else {
		c.JSON(http.StatusConflict, ErrorsToJSON("error"))
	}

}
