package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/collector"
	"github.com/pilly-io/api/internal/apis/frontend"
	"github.com/pilly-io/api/internal/apis/middlewares"
	"github.com/pilly-io/api/internal/db"
)

// SetupRouter binds all the routes to their handlers
func SetupRouter(r *gin.Engine, database db.Database) {
	clusters := frontend.ClustersHandler{DB: database}
	metrics := frontend.MetricsHandler{DB: database}
	v1 := r.Group("/api/v1")
	v1.POST("/clusters", clusters.Create)
	v1.GET("/clusters", clusters.List)
	v1.GET("/clusters/:id/owners/metrics", metrics.ListOwners)
	v1.GET("/clusters/:id/namespaces/metrics", metrics.ListNamespaces)

	collectorGroup := r.Group("/api/v1/collector")
	collectorGroup.Use(middlewares.CluserAuthMiddleware(database.Clusters()))
	collectorNodes := collector.NodesHandler{DB: database}
	collectorNamespaces := collector.NamespacesHandler{DB: database}
	collectorOwners := collector.OwnersHandler{DB: database}
	collectorMetrics := collector.MetricsHandler{DB: database}
	collectorGroup.POST("/nodes", collectorNodes.Sync)
	collectorGroup.POST("/namespaces", collectorNamespaces.Sync)
	collectorGroup.POST("/owners", collectorOwners.Sync)
	collectorGroup.POST("/metrics", collectorMetrics.Create)
}
