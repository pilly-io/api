package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
)

// OwnersHandler definition
type OwnersHandler struct {
	DB db.Database
}

// ListMetrics of owners within a period
func (handler *OwnersHandler) ListMetrics(c *gin.Context) {
	// 1. Check sanity of the request
	query := db.Query{
		Conditions: db.QueryConditions{"id": c.Param("id")},
	}
	if !handler.DB.Clusters().Exists(query) {
		c.JSON(http.StatusNotFound, ErrorsToJSON(errors.New("cluster_does_not_exist")))
		return
	}
	start, err := strconv.ParseInt(c.Query("start"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_start")))
		return
	}
	end, err := strconv.ParseInt(c.Query("end"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_end")))
		return
	}
	period, err := strconv.ParseInt(c.Query("period"), 10, 64)
	if err != nil || period < MinPeriod {
		c.JSON(http.StatusBadRequest, ErrorsToJSON(errors.New("invalid_period")))
		return
	}
	//2. Get the metrics
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(period)
}

// Usage of owners within a period
func (handler *OwnersHandler) Usage(c *gin.Context) {
}
