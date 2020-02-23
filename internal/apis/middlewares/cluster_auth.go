package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

const ClusterAuthHeaderName = "X-API-TOKEN"

// CluserAuthMiddleware authenticates incoming requests by looking at the
// X-API-TOKEN header and inject the corresponding cluster into the gin context
func CluserAuthMiddleware(table *db.ClustersTable) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiToken := c.Request.Header.Get(ClusterAuthHeaderName)
		cluster := models.Cluster{}
		query := db.Query{
			Conditions: db.QueryConditions{"api_token": apiToken},
		}
		if err := table.Find(query, &cluster); err == nil {
			c.Set("cluster", &cluster)
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
