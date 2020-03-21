package frontend

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/helpers"
	"github.com/pilly-io/api/internal/models"
)

// MetricsHandler definition
type MetricsHandler struct {
	DB db.Database
}

const MinPeriod = 60 // in seconds

// ValidateRequest : Validate the cluster, interval, period, namespace, owners
func (handler *MetricsHandler) ValidateRequest(c *gin.Context) bool {
	clusterID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	var ownerUIDs []string
	var namespace models.Namespace
	query := db.Query{
		Conditions: db.QueryConditions{"id": clusterID},
	}
	if err != nil || !handler.DB.Clusters().Exists(query) {
		c.JSON(http.StatusNotFound, utils.ErrorsToJSON(errors.New("cluster_does_not_exist")))
		return false
	}
	start, err := utils.ConvertTimestampToTime(c.Query("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorsToJSON(errors.New("invalid_start")))
		return false
	}
	end, err := utils.ConvertTimestampToTime(c.Query("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorsToJSON(errors.New("invalid_end")))
		return false
	}
	period, err := strconv.ParseInt(c.DefaultQuery("period", string(MinPeriod)), 10, 64)
	if err != nil || period < MinPeriod {
		c.JSON(http.StatusBadRequest, utils.ErrorsToJSON(errors.New("invalid_period")))
		return false
	}
	namespaceName := c.Query("namespace")
	if namespaceName != "" {
		query = db.Query{
			Conditions: db.QueryConditions{"cluster_id": clusterID, "name": namespaceName},
		}
		err := handler.DB.Namespaces().Find(query, &namespace)
		if err != nil {
			c.JSON(http.StatusNotFound, utils.ErrorsToJSON(errors.New("namespace_does_not_exist")))
			return false
		}
		if c.Query("owners") != "" {
			for _, owner := range strings.Split(c.Query("owners"), ",") {
				details := strings.Split(owner, "/")
				if len(details) == 2 {
					kind := utils.GetFullKindName(details[0])
					query = db.Query{
						Conditions: db.QueryConditions{"cluster_id": clusterID, "namespace_uid": namespace.UID, "type": kind, "name": details[1]},
					}
					owner := models.Owner{}
					err := handler.DB.Owners().Find(query, &owner)
					if err == nil {
						ownerUIDs = append(ownerUIDs, owner.UID)
					}
				}
			}
			if len(ownerUIDs) == 0 {
				c.JSON(http.StatusNotFound, utils.ErrorsToJSON(errors.New("owners_do_not_exist")))
				return false
			}
		}
	}
	c.Set("Start", *start)
	c.Set("End", *end)
	c.Set("Period", period)
	c.Set("NamespaceUID", namespace.UID)
	c.Set("OwnerUIDs", ownerUIDs)
	return true
}

// ListOwners List the metrics of the cluster owners within an interval
func (handler *MetricsHandler) ListOwners(c *gin.Context) {
	var owners []*models.Owner
	// 1. Check sanity of the request
	if !handler.ValidateRequest(c) {
		return
	}
	start := c.Value("Start").(time.Time)
	end := c.Value("End").(time.Time)
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	namespaceUID := c.Value("NamespaceUID").(string)
	period := c.Value("Period").(int64)
	ownerUIDs := c.Value("OwnerUIDs").([]string)

	//2. Get all the owners or specific ones
	interval := db.QueryInterval{Start: start, End: end}
	conditions := db.QueryConditions{"cluster_id": clusterID}
	if namespaceUID != "" {
		conditions["namespaceUID"] = namespaceUID
		if len(ownerUIDs) > 0 {
			conditions["uid__in"] = ownerUIDs
		}
	}
	query := db.Query{
		Conditions: conditions,
		Interval:   &interval,
	}
	handler.DB.Owners().FindAll(query, &owners)
	GetOwnerUIDs := func(owners []*models.Owner) []string {
		var uids []string
		for _, owner := range owners {
			uids = append(uids, owner.UID)
		}
		return uids
	}
	ownerUIDs = GetOwnerUIDs(owners)
	//3. Get the metrics within the interval grouped by period
	metrics, _ := handler.DB.Metrics().FindAll(uint(clusterID), ownerUIDs, uint(period), interval)
	metricsIndexed := helpers.IndexMetrics(metrics, "owner")
	//4. Compute the owners resources
	handler.DB.Owners().ComputeResources(&owners, metricsIndexed)
	//5. This is the end
	c.JSON(http.StatusOK, utils.ObjectToJSON(&owners))
}

// ListNamespaces List the metrics of the cluster namespaces within an interval
func (handler *MetricsHandler) ListNamespaces(c *gin.Context) {
	var namespaces []*models.Namespace
	// 1. Check sanity of the request
	if !handler.ValidateRequest(c) {
		return
	}
	start := c.Value("Start").(time.Time)
	end := c.Value("End").(time.Time)
	clusterID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	period := c.Value("Period").(int64)

	interval := db.QueryInterval{Start: start, End: end}
	conditions := db.QueryConditions{"cluster_id": clusterID}

	query := db.Query{
		Conditions: conditions,
		Interval:   &interval,
	}
	handler.DB.Namespaces().FindAll(query, &namespaces)
	//3. Get the metrics within the interval grouped by period
	metrics, _ := handler.DB.Metrics().FindAll(uint(clusterID), nil, uint(period), interval)
	metricsIndexed := helpers.IndexMetrics(metrics, "namespace")
	//4. Compute the owners resources
	handler.DB.Namespaces().ComputeResources(&namespaces, metricsIndexed)
	//5. This is the end
	c.JSON(http.StatusOK, utils.ObjectToJSON(&namespaces))
}
