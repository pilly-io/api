package collector

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/apis/middlewares"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestNodesHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NodesHandler Suite")
}

var _ = Describe("NodesHandler", func() {
	var (
		engine   *gin.Engine
		database db.Database
		cluster  *models.Cluster
	)
	BeforeEach(func() {
		database = tests.SetupDB()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		SetupRouter(engine, database)
		cluster, _ = database.Clusters().Create("test", "aws")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("Sync()", func() {
		It("should create node if does not exist", func() {
			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			count := 0
			countQuery := db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
			}
			count = database.Nodes().Count(countQuery)
			Expect(count).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			count = database.Nodes().Count(countQuery)
			Expect(count).To(Equal(1))

			Expect(res.Code).To(Equal(201))
		})

		It("should update nodes count of cluster", func() {
			var clusterFromDB models.Cluster
			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			clusterQuery := db.Query{
				Conditions: db.QueryConditions{
					"id": cluster.ID,
				},
			}

			database.Clusters().Find(clusterQuery, &clusterFromDB)
			Expect(clusterFromDB.NodesCount).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Clusters().Find(clusterQuery, &clusterFromDB)
			Expect(clusterFromDB.NodesCount).To(Equal(1))
		})

		It("should update region of cluster if not set", func() {
			var clusterFromDB models.Cluster
			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			clusterQuery := db.Query{
				Conditions: db.QueryConditions{
					"id": cluster.ID,
				},
			}

			database.Clusters().Find(clusterQuery, &clusterFromDB)
			Expect(clusterFromDB.Region).To(Equal(""))

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Clusters().Find(clusterQuery, &clusterFromDB)
			Expect(clusterFromDB.Region).To(Equal("euwest1"))
		})

		It("should mark nodes as deleted if not sent", func() {
			deletedNode := models.Node{UID: "cbd46a8e-faa1-4f2a-a826-f45169d5ba14", ClusterID: cluster.ID}
			database.Nodes().Insert(&deletedNode)

			Expect(deletedNode.DeletedAt).To(BeNil())

			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Nodes().Find(db.Query{
				Conditions: db.QueryConditions{"id": deletedNode.ID},
			}, &deletedNode)

			Expect(deletedNode.DeletedAt).ToNot(BeNil())
		})

		It("should update current nodes with new labels", func() {
			currentNode := models.Node{
				UID:       "816d2d42-4dd0-4e05-97eb-f077983b73dc",
				ClusterID: cluster.ID,
				Labels: map[string]interface{}{
					"key1": "oldVal",
				},
			}
			database.Nodes().Insert(&currentNode)

			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Nodes().Find(db.Query{
				Conditions: db.QueryConditions{"id": currentNode.ID},
			}, &currentNode)

			Expect(currentNode.Labels).To(Equal(map[string]interface{}{
				"key1": "val1",
			}))
		})
	})
})
