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

func TestNamespacesHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NamespacesHandler Suite")
}

var _ = Describe("NamespacesHandler", func() {
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
		It("should create namespaces if does not exist", func() {
			data := tests.LoadFile("testdata/namespaces.json")
			res := httptest.NewRecorder()

			count := 0
			countQuery := db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
			}
			count = database.Namespaces().Count(countQuery)
			Expect(count).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/namespaces", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			count = database.Namespaces().Count(countQuery)
			Expect(count).To(Equal(2))

			Expect(res.Code).To(Equal(201))
		})

		It("should mark namespace as deleted if not sent", func() {
			deletedNamespace := models.Namespace{Name: "deleted", ClusterID: cluster.ID}
			database.Namespaces().Insert(&deletedNamespace)

			Expect(deletedNamespace.DeletedAt).To(BeNil())

			data := tests.LoadFile("testdata/namespaces.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/namespaces", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Namespaces().Find(db.Query{
				Conditions: db.QueryConditions{"id": deletedNamespace.ID},
			}, &deletedNamespace)

			Expect(deletedNamespace.DeletedAt).ToNot(BeNil())
		})

		It("should update current namespaces with new labels", func() {
			currentNamespace := models.Namespace{
				Name:      "monitoring",
				UID:       "0238e796-9409-4152-a370-5d57a79ec6a6",
				ClusterID: cluster.ID,
				Labels: map[string]interface{}{
					"key1": "oldVal",
				},
			}
			database.Namespaces().Insert(&currentNamespace)

			data := tests.LoadFile("testdata/namespaces.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/namespaces", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Namespaces().Find(db.Query{
				Conditions: db.QueryConditions{"id": currentNamespace.ID},
			}, &currentNamespace)

			Expect(currentNamespace.Labels).To(Equal(map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			}))
		})
	})
})
