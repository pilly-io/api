package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

// NodesHandler definition
type NodesHandler struct {
	DB db.Database
}

// Sync synchronize nodes of cluster
func (handler *NodesHandler) Sync(c *gin.Context) {
	var nodes, existingNodes []models.Node
	nodesTable := handler.DB.Nodes()
	c.BindJSON(&nodes)

	fmt.Println(nodes)
	cluster := c.MustGet("cluster").(*models.Cluster)

	existingNodesQuery := db.Query{
		Conditions: db.QueryConditions{"cluster_id": cluster.ID},
	}

	nodesTable.FindAll(existingNodesQuery, &existingNodes)
	existingNodesByUID := indexNodesByUID(&existingNodes)

	for _, node := range nodes {
		existingNode, ok := existingNodesByUID[node.UID]
		if ok {
			existingNode.Labels = node.Labels
			
		} else {
			nodesTable.Insert(node)
		}
	}
	cluster.NodesCount = len(nodes)
	if cluster.NodesCount > 0 && cluster.Region != "" {
		cluster.Region = nodes[0].Region
	}
}

func indexNodesByUID(nodes *[]models.Node) map[string]models.Node {
	nodesByUID := make(map[string]models.Node)
	for _, node := range *nodes {
		nodesByUID[node.UID] = node
	}
	return nodesByUID
}
