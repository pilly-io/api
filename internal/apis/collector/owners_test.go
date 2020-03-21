package collector

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/apis"
	"github.com/pilly-io/api/internal/apis/middlewares"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestOwnersHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OwnersHandler Suite")
}

var _ = Describe("OwnersHandler", func() {
	var (
		engine   *gin.Engine
		database db.Database
		cluster  *models.Cluster
	)

	BeforeEach(func() {
		database = tests.GetDB()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		apis.SetupRouter(engine, database)
		cluster, _ = database.Clusters().Create("test", "aws")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("Sync()", func() {
		It("should create owners if does not exist", func() {
			data := tests.LoadFile("testdata/owners.json")
			res := httptest.NewRecorder()

			count := 0
			countQuery := db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
			}
			count = database.Owners().Count(countQuery)
			Expect(count).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/owners", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			count = database.Owners().Count(countQuery)
			Expect(count).To(Equal(1))

			Expect(res.Code).To(Equal(201))
		})

		It("should update current owners with new labels", func() {
			currentOwner := models.Owner{
				UID:       "816d2d42-4dd0-4e05-97eb-f077983b73dc",
				Name:      "Grafana",
				ClusterID: cluster.ID,
				Labels: map[string]interface{}{
					"key1": "oldVal",
				},
			}
			database.Owners().Insert(&currentOwner)

			data := tests.LoadFile("testdata/owners.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/owners", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Owners().Find(db.Query{
				Conditions: db.QueryConditions{"id": currentOwner.ID},
			}, &currentOwner)

			Expect(currentOwner.Labels).To(Equal(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			}))
		})
	})
})
