package collector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// OwnersHandler definition
type OwnersHandler struct {
	DB db.Database
}

// Sync synchronize owners of cluster
func (handler *OwnersHandler) Sync(c *gin.Context) {
	var owners, existingOwners []*models.Owner
	ownersTable := handler.DB.Owners()
	c.BindJSON(&owners)

	cluster := c.MustGet("cluster").(*models.Cluster)

	existingOwnersQuery := db.Query{
		Conditions: db.QueryConditions{"cluster_id": cluster.ID},
	}

	ownersTable.FindAll(existingOwnersQuery, &existingOwners)
	existingOwnersByUID := indexOwnersByUID(&existingOwners)

	// Merge owners infos beased on their UID
	for _, owner := range owners {
		owner.ClusterID = cluster.ID
		existingOwner, ok := existingOwnersByUID[owner.UID]
		if ok {
			existingOwner.Labels = owner.Labels
			ownersTable.Update(&existingOwner)
		} else {
			err := ownersTable.Insert(owner)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorsToJSON(err))
			}
		}
	}

	c.JSON(http.StatusCreated, utils.ObjectToJSON(nil))
}

func indexOwnersByUID(owners *[]*models.Owner) map[string]models.Owner {
	ownersByUID := make(map[string]models.Owner)
	for _, owner := range *owners {
		ownersByUID[owner.UID] = *owner
	}
	return ownersByUID
}
