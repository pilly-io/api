package apis

import (
	"errors"
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

// ValidateRequest : Validate the cluster and start/end
func (handler *MetricsHandler) ValidateRequest(c *gin.Context) bool {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	query := db.Query{
		Conditions: db.QueryConditions{"id": clusterID},
	}
	if err != nil || !handler.DB.Clusters().Exists(query) {
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
	// 1. Check sanity of the request
	if !handler.ValidateRequest(c) {
		return
	}
	start := c.Value("Start").(time.Time)
	end := c.Value("End").(time.Time)
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	period, err := strconv.ParseInt(c.Query("period"), 10, 64)
	if err != nil || period < MinPeriod {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_period")))
		return
	}
	//2. Get the metrics within the interval grouped by period
	interval := db.QueryInterval{Start: start, End: end}
	metrics, _ := handler.DB.Metrics().FindAll(uint(clusterID), uint(period), interval)
	metricsByOwnerAndTimestamp, ownerUIDs := indexMetricsByOwnerAndTimestamp(metrics)
	//3. Get the owners within the interval
	query := db.Query{
		Conditions: db.QueryConditions{"cluster_id": clusterID, "uid__in": ownerUIDs},
		Interval:   &interval,
	}
	handler.DB.Owners().FindAll(query, &owners)
	//4. Compute the owners resources
	handler.DB.Owners().ComputeResources(&owners, metricsByOwnerAndTimestamp)
	//5. This is the end
	c.JSON(http.StatusOK, ObjectToJSON(&owners))
}

//indexMetricsByOwnerAndTimestamp : this is ugly but life is ugly
func indexMetricsByOwnerAndTimestamp(metrics *[]models.Metric) (*models.IndexedMetrics, []string) {
	var ownerUIDs []string
	metricsByOwnerAndTimestamp := make(models.IndexedMetrics)
	for _, metric := range *metrics {
		if _, exist := metricsByOwnerAndTimestamp[metric.OwnerUID]; !exist {
			metricsByOwnerAndTimestamp[metric.OwnerUID] = map[time.Time]map[string]models.Metric{metric.Period: {metric.Name: metric}}
			ownerUIDs = append(ownerUIDs, metric.OwnerUID)
		} else {
			if _, exist := metricsByOwnerAndTimestamp[metric.OwnerUID][metric.Period]; !exist {
				metricsByOwnerAndTimestamp[metric.OwnerUID][metric.Period] = map[string]models.Metric{metric.Name: metric}
			} else {
				metricsByOwnerAndTimestamp[metric.OwnerUID][metric.Period][metric.Name] = metric
			}
		}
	}
	return &metricsByOwnerAndTimestamp, ownerUIDs
}
