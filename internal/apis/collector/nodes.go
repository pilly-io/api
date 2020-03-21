package collector

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// NodesHandler definition
type NodesHandler struct {
	DB db.Database
}

// Sync synchronize nodes of cluster
func (handler *NodesHandler) Sync(c *gin.Context) {
	var nodes, existingNodes []*models.Node
	nodesTable := handler.DB.Nodes()
	c.BindJSON(&nodes)

	cluster := c.MustGet("cluster").(*models.Cluster)

	existingNodesQuery := db.Query{
		Conditions: db.QueryConditions{"cluster_id": cluster.ID},
	}

	nodesTable.FindAll(existingNodesQuery, &existingNodes)
	existingNodesByUID := indexNodesByUID(&existingNodes)
	nodesByUID := indexNodesByUID(&nodes)

	// Merge nodes infos beased on their UID
	for _, node := range nodes {
		node.ClusterID = cluster.ID
		existingNode, ok := existingNodesByUID[node.UID]
		if ok {
			existingNode.Labels = node.Labels
			nodesTable.Update(&existingNode)
		} else {
			err := nodesTable.Insert(node)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorsToJSON(err))
			}
		}
	}

	// Mark nodes as deleted if not received
	nodeIDsToDelete := make([]uint, 0)
	for _, existingNode := range existingNodes {
		if _, ok := nodesByUID[existingNode.UID]; ok == false {
			nodeIDsToDelete = append(nodeIDsToDelete, existingNode.ID)
		}
	}
	if len(nodeIDsToDelete) > 0 {
		nodesTable.Delete(db.Query{
			Conditions: db.QueryConditions{
				"id__in": nodeIDsToDelete,
			},
		}, true)
	}

	// Update cluster's nodes count and region if not set
	cluster.NodesCount = len(nodes)
	if cluster.NodesCount > 0 && cluster.Region == "" {
		cluster.Region = nodes[0].Region
	}

	handler.DB.Clusters().Update(cluster)
	c.JSON(http.StatusCreated, utils.ObjectToJSON(nil))
}

func indexNodesByUID(nodes *[]*models.Node) map[string]models.Node {
	nodesByUID := make(map[string]models.Node)
	for _, node := range *nodes {
		nodesByUID[node.UID] = *node
	}
	return nodesByUID
}
