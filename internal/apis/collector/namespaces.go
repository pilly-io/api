package collector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// NamespacesHandler definition
type NamespacesHandler struct {
	DB db.Database
}

// Sync synchronize namespaces of cluster
func (handler *NamespacesHandler) Sync(c *gin.Context) {
	var namespaces, existingNamespaces []*models.Namespace
	namespacesTable := handler.DB.Namespaces()
	c.BindJSON(&namespaces)

	cluster := c.MustGet("cluster").(*models.Cluster)

	existingNamespacesQuery := db.Query{
		Conditions: db.QueryConditions{"cluster_id": cluster.ID},
	}

	namespacesTable.FindAll(existingNamespacesQuery, &existingNamespaces)
	existingNamespacesByUID := indexNamespacesByUID(&existingNamespaces)
	namespacesByUID := indexNamespacesByUID(&namespaces)

	// Merge namespaces infos beased on their UID
	for _, namespace := range namespaces {
		namespace.ClusterID = cluster.ID
		existingNamespace, ok := existingNamespacesByUID[namespace.UID]
		if ok {
			existingNamespace.Labels = namespace.Labels
			namespacesTable.Update(&existingNamespace)
		} else {
			err := namespacesTable.Insert(namespace)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorsToJSON(err))
			}
		}
	}

	// Mark namespaces as deleted if not received
	namespaceIDsToDelete := make([]uint, 0)
	for _, existingNamespace := range existingNamespaces {
		if _, ok := namespacesByUID[existingNamespace.UID]; ok == false {
			namespaceIDsToDelete = append(namespaceIDsToDelete, existingNamespace.ID)
		}
	}
	if len(namespaceIDsToDelete) > 0 {
		namespacesTable.Delete(db.Query{
			Conditions: db.QueryConditions{
				"id__in": namespaceIDsToDelete,
			},
		}, true)
	}

	c.JSON(http.StatusCreated, utils.ObjectToJSON(nil))
}

func indexNamespacesByUID(namespaces *[]*models.Namespace) map[string]models.Namespace {
	namepsacesByUID := make(map[string]models.Namespace)
	for _, namespace := range *namespaces {
		namepsacesByUID[namespace.UID] = *namespace
	}
	return namepsacesByUID
}
