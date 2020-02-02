package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/daos"
	"github.com/pilly-io/api/internal/services"
	log "github.com/sirupsen/logrus"
)

func GetClusterByName(c *gin.Context) {
	svc := services.NewClusterService(daos.NewClusterDao())
	cluster, err := svc.GetByName(c.Param("name"))
	if err != nil {
		log.Error("Cannot retrieve the cluster")
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}

func GetClusterById(c *gin.Context) {
	svc := services.NewClusterService(daos.NewClusterDao())
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	cluster, err := svc.GetByID(uint(id))
	if err != nil {
		log.Error("Cannot retrieve the cluster")
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, cluster)
	}
}
