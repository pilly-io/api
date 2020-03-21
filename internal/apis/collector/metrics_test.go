package collector_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/apis/middlewares"
	"github.com/pilly-io/api/internal/apis/router"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestMetricsHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MetricsHandler Suite")
}

var _ = Describe("MetricsHandler", func() {
	var (
		engine   *gin.Engine
		database db.Database
		cluster  *models.Cluster
	)

	BeforeEach(func() {
		database = tests.GetDB()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		router.SetupRouter(engine, database)
		cluster, _ = database.Clusters().Create("test", "aws")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("Create()", func() {
		It("should create 2 metrics", func() {
			data := tests.LoadFile("testdata/metrics.json")
			res := httptest.NewRecorder()

			count := 0
			countQuery := db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
			}
			count = database.Namespaces().Count(countQuery)
			Expect(count).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/metrics", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			count = database.Metrics().Count(countQuery)
			Expect(count).To(Equal(2))

			Expect(res.Code).To(Equal(201))
		})

		It("should create metrics with proper data", func() {
			data := tests.LoadFile("testdata/metrics.json")
			res := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/api/v1/collector/metrics", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			var metrics []*models.Metric
			database.Metrics().FindAll(db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
				OrderBy: "created_at",
			}, &metrics)
			Expect(len(metrics)).To(Equal(2))

			metric1, metric2 := metrics[0], metrics[1]

			Expect(metric1.ClusterID).To(Equal(cluster.ID))
			Expect(metric1.OwnerUID).To(Equal("c5dea447-7914-40c2-9382-227c44731ee0"))
			Expect(metric1.NamespaceUID).To(Equal("723d8295-50aa-49e7-a325-d72541539174"))
			Expect(metric1.Name).To(Equal("cpu_used"))
			Expect(metric1.Value).To(Equal(0.12))
			Expect(metric1.CreatedAt.UTC()).To(Equal(time.Date(2020, 03, 21, 12, 43, 0, 0, time.UTC)))

			Expect(metric2.ClusterID).To(Equal(cluster.ID))
			Expect(metric2.OwnerUID).To(Equal("3fa8c97e-6f97-49c4-a48c-e971de55224d"))
			Expect(metric2.NamespaceUID).To(Equal("8f8d64bc-ec22-41e2-a3d7-b4b83776c2a1"))
			Expect(metric2.Name).To(Equal("memory_used"))
			Expect(metric2.Value).To(Equal(130.4))
			Expect(metric2.CreatedAt.UTC()).To(Equal(time.Date(2020, 03, 21, 13, 43, 0, 0, time.UTC)))
		})
	})
})
