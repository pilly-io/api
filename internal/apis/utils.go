package apis

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
)

type jsonFormat = map[string]interface{}

const MinPeriod = 60 // in seconds
const MetricCPUUsed = "fill_with_proper_cpu_value"
const MetricMemoryUsed = "fill_with_proper_memory_value"

// SetupRouter binds all the routes to their handlers
func SetupRouter(r *gin.Engine, database *db.GormDatabase) {
	clusters := ClustersHandler{DB: database}
	owners := OwnersHandler{DB: database}
	v1 := r.Group("/api/v1")
	v1.POST("/clusters", clusters.Create)
	v1.GET("/clusters", clusters.List)
	v1.GET("/clusters/:id/owners/metrics", owners.ListMetrics)
	v1.GET("/clusters/:id/owners/usage", owners.Usage)

	collector := r.Group("/api/v1/collector")
	nodes := NodesHandler{DB: database}
	collector.POST("/nodes", nodes.Sync)
}

// ConvertStringToTime : convert a string to a time
func ConvertStringToTime(str string) (time.Time, error) {
	toInt, err := strconv.ParseInt(str, 10, 64)
	toTime := time.Unix(int64(toInt), 0)
	return toTime, err
}

// ErrorsToJSON returns JSON format for multiple errors
func ErrorsToJSON(errors ...error) map[string]interface{} {
	var errorsString []string
	for _, e := range errors {
		errorsString = append(errorsString, e.Error())
	}
	return jsonFormat{
		"errors": errorsString,
	}
}

//ObjectToJSON returns JSON format for the data
func ObjectToJSON(object interface{}) map[string]interface{} {
	return jsonFormat{
		"data": object,
	}
}

/*func PaginatedObjectToJSON(collection db.PaginationInfo) map[string]interface{} {
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
