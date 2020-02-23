package apis

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
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

// ValidateRequest : Validate the cluster, interval, period
func (handler *MetricsHandler) ValidateRequest(c *gin.Context) bool {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var ownerUIDs []string
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
	period, err := strconv.ParseInt(c.DefaultQuery("period", string(MinPeriod)), 10, 64)
	if err != nil || period < MinPeriod {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_period")))
		return false
	}
	namespace := c.Query("namespace")
	if namespace != "" {
		query = db.Query{
			Conditions: db.QueryConditions{"id": clusterID, "name": namespace},
		}
		if !handler.DB.Namespaces().Exists(query) {
			c.JSON(http.StatusNotFound, ErrorsToJSON(errors.New("namespace_does_not_exist")))
			return false
		}

		if c.Query("owners") != "" {
			for _, owner := range strings.Split(c.Query("owners"), ",") {
				details := strings.Split(owner, "/")
				if len(details) == 2 {
					//kind := GetFullKindName(details[0])
				}
			}
		}
	}
	c.Set("Start", *start)
	c.Set("End", *end)
	c.Set("Period", period)
	c.Set("OwnerUIDs", ownerUIDs)
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
	period := c.Value("Period").(int64)

	//2. Get all the owners or specific ones
	interval := db.QueryInterval{Start: start, End: end}
	query := db.Query{
		Conditions: db.QueryConditions{"cluster_id": clusterID},
		Interval:   &interval,
	}
	handler.DB.Owners().FindAll(query, &owners)
	//3. Get the metrics within the interval grouped by period
	metrics, _ := handler.DB.Metrics().FindAll(uint(clusterID), uint(period), interval)
	metricsIndexed := indexMetrics(metrics)
	//4. Compute the owners resources
	handler.DB.Owners().ComputeResources(&owners, metricsIndexed)
	//5. This is the end
	c.JSON(http.StatusOK, ObjectToJSON(&owners))
}

//indexMetrics : this is ugly but life is ugly
func indexMetrics(metrics *[]models.Metric) *models.IndexedMetrics {
	metricsIndexed := make(models.IndexedMetrics)
	for _, metric := range *metrics {
		if _, exist := metricsIndexed[metric.OwnerUID]; !exist {
			metricsIndexed[metric.OwnerUID] = map[time.Time]map[string]models.Metric{metric.Period: {metric.Name: metric}}
		} else {
			if _, exist := metricsIndexed[metric.OwnerUID][metric.Period]; !exist {
				metricsIndexed[metric.OwnerUID][metric.Period] = map[string]models.Metric{metric.Name: metric}
			} else {
				metricsIndexed[metric.OwnerUID][metric.Period][metric.Name] = metric
			}
		}
	}
	return &metricsIndexed
}
