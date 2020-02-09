package apis

import "github.com/pilly-io/api/internal/db"
import "github.com/gin-gonic/gin"

type jsonFormat = map[string]interface{}

//TODO: maybe move it to another file ? routes.go ?

// SetupRouter binds all the routes to their handlers
func SetupRouter(r *gin.Engine, database *db.GormDatabase) {
	clusters := ClustersHandler{DB: database}
	v1 := r.Group("/api/v1")
	v1.POST("/clusters", clusters.Create)
}

// ErrorsToJSON returns JSON format for multiple errors
func ErrorsToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"errors": "errors",
	}
}

func ObjectToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"data": object,
	}
}

/*func ObjectsToJSON(collection db.PaginationInfo) map[string]interface{} {
	return jsonFormat{
		"data": collection.Objects,
		"pagination": {
			"current":     collection.CurrentPage,
			"limit":       collection.Limit,
			"max_page":    collection.MaxPage,
			"total_count": collection.TotalCount,
		},
	}
}*/
