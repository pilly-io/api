package frontend

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
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
		cluster, err := handler.DB.Clusters().Create(cluster.Name, cluster.Provider)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.ErrorsToJSON(err))
		} else {
			c.JSON(http.StatusCreated, utils.ObjectToJSON(&cluster))
		}
	} else {
		c.JSON(http.StatusConflict, utils.ErrorsToJSON(errors.New("already_exist")))
	}

}

// List all the existing clusters
func (handler *ClustersHandler) List(c *gin.Context) {
	var clusters []*models.Cluster
	query := db.Query{}

	_, err := handler.DB.Clusters().FindAll(query, &clusters)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorsToJSON(err))
	} else {
		c.JSON(http.StatusOK, ObjectToJSON(&clusters))
	}
}
