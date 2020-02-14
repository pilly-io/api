package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// MetricsHandler definition
type MetricsHandler struct {
	DB db.Database
}

const MinPeriod = 60 // in seconds
const MetricCPUUsed = "fill_with_proper_cpu_value"
const MetricMemoryUsed = "fill_with_proper_memory_value"

// ValidateRequest : Validate the cluster and start/end
func (handler *MetricsHandler) ValidateRequest(c *gin.Context) bool {
	query := db.Query{
		Conditions: db.QueryConditions{"id": c.Param("id")},
	}
	if !handler.DB.Clusters().Exists(query) {
		c.JSON(http.StatusNotFound, ErrorsToJSON(errors.New("cluster_does_not_exist")))
		return false
	}
	start, err := ConvertTimestampToTime(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_start")))
		return false
	}
	end, err := ConvertTimestampToTime(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_end")))
		return false
	}
	c.Set("Start", *start)
	c.Set("End", *end)
	return true
}

// List the metrics of the cluster within an interval
func (handler *MetricsHandler) List(c *gin.Context) {
	var owners []models.Owner
	//var metrics []models.Metric
	// 1. Check sanity of the request
	if !handler.ValidateRequest(c) {
		return
	}
	itStart, _ := c.Get("Start")
	itEnd, _ := c.Get("End")
	start, _ := itStart.(time.Time)
	end, _ := itEnd.(time.Time)
	clusterID := c.Param("id")
	period, err := strconv.ParseInt(c.Query("period"), 10, 64)
	if err != nil || period < MinPeriod {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_period")))
		return
	}
	//2. Get the owners within the interval
	query := db.Query{
		Conditions: db.QueryConditions{"cluster_id": clusterID},
		Interval:   &db.QueryInterval{Start: start, End: end},
	}
	handler.DB.Owners().FindAll(query, &owners)
	//3. Get the metrics within the interval grouped by period
	fmt.Println(period)
	c.JSON(http.StatusOK, ObjectToJSON(&owners))
}
