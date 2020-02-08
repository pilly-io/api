package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// Clusters endpoints
type ClustersHandler struct {
	db db.Database
}

// FindAll get all clusters
// func (handler *ClustersHandler) FindAll(c *gin.Context) {
// 	clusters := []*models.Cluster{}
// 	query := db.Query{
// 		Result: &clusters,
// 	}
// 	err := handler.db.Clusters().FindAll(query)

// 	if err != nil {
// 		log.Errorf("Cannot retrieve the clusters: %s", err)
// 		c.AbortWithStatus(http.StatusNotFound)
// 	} else {
// 		c.JSON(http.StatusOK, clusters)
// 	}
// }

func (handler *ClustersHandler) Create(c *gin.Context) {
	cluster = models.Cluster{}
	c.BindJSON(&cluster)
	query := db.Query{
		Conditions: db.QueryConditions{"name": cluster.Name},
	}

	if !handler.db.Clusters().Exists(query) {
		cluster, err = handler.db.Cluster().Create(cluster.Name)
		if err != nil {
			// c.JSON(http.StatusCreated, ObjectToJSON(&cluster))
		} else {
			// c.JSON(http.StatusUnprocessableEntity, ErrorsToJSON([err]))
		}
	} else {
		// c.JSON(http.StatusConflict, ErrorsToJSON([errors.New("Cluster already exist")]))
	}

}
