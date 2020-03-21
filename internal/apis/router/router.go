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
	nodes := collector.NodesHandler{DB: database}
	namespaces := collector.NamespacesHandler{DB: database}
	owners := collector.OwnersHandler{DB: database}
	collectorGroup.POST("/nodes", nodes.Sync)
	collectorGroup.POST("/namespaces", namespaces.Sync)
	collectorGroup.POST("/owners", owners.Sync)
}
