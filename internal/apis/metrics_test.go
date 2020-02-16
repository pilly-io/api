package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestMetricsHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MetricsHandler Suite")
}

//MetricsFactory: create a metric given an owner/cluster
func MetricsFactory(database db.Database, clusterID uint, ownerUID string, name string, value float64, createdAt time.Time) *models.Metric {
	metrics := &models.Metric{
		ClusterID: clusterID,
		OwnerUID:  ownerUID,
		Name:      name,
		Value:     value,
	}
	metrics.Model.CreatedAt = createdAt
	return metrics
}

//VeryOldOwner : create an owner and generate its metrics, will create a factory if needed
func VeryOldOwner(database db.Database, clusterID uint) (*models.Owner, []*models.Metric) {
	var metrics []*models.Metric
	uid, _ := uuid.NewRandom()
	past := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	owner := &models.Owner{
		ClusterID: clusterID,
		Name:      "toto",
		Type:      "deployment",
		Namespace: "default",
		UID:       uid.String(),
	}
	owner.Model.CreatedAt = past
	database.Owners().Insert(&owner)
	// TODO convert into a []{}
	//1. Generate CPU metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPUUsed, 1, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPUUsed, 1, past.Add(60)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPUUsed, 2, past.Add(120)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPUUsed, 2, past.Add(180)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPUUsed, 1, past.Add(240)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPURequested, 2, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPURequested, 2, past.Add(60)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPURequested, 2, past.Add(120)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPURequested, 2, past.Add(180)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricCPURequested, 2, past.Add(240)))
	//2. Generate Memory metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryUsed, 1024, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryUsed, 1024, past.Add(60)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryUsed, 2048, past.Add(120)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryUsed, 2048, past.Add(180)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryUsed, 1024, past.Add(240)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryRequested, 2048, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryRequested, 2048, past.Add(60)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryRequested, 2048, past.Add(120)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryRequested, 2048, past.Add(180)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, MetricMemoryRequested, 2048, past.Add(240)))
	return owner, metrics
}

var _ = Describe("Owners", func() {
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

	Describe("ListMetrics() fails", func() {
		It("Should return a 404 as cluster does not exist", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/xxxx/owners/metrics", nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(404))
		})
		It("Should return a 400 as start is not defined", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
		})
		It("Should return a 400 as start is an invalid timestamp", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "start")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
			json.Unmarshal(res.Body.Bytes(), &payload)
			errors := payload["errors"].([]interface{})
			Expect(errors[0]).To(Equal("invalid_start"))
		})
		It("Should return a 400 as end is not defined", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "1")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
		})
		It("Should return a 400 as end is an invalid timestamp", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "1")
			q.Add("end", "end")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
			json.Unmarshal(res.Body.Bytes(), &payload)
			errors := payload["errors"].([]interface{})
			Expect(errors[0]).To(Equal("invalid_end"))
		})
		It("Should return a 400 as period is not defined", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "1")
			q.Add("end", "2")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
		})
		It("Should return a 400 as period is an invalid int", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "1")
			q.Add("end", "2")
			q.Add("period", "3")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
			json.Unmarshal(res.Body.Bytes(), &payload)
			errors := payload["errors"].([]interface{})
			Expect(errors[0]).To(Equal("invalid_period"))
		})
	})
	Describe("ListMetrics() succeeds", func() {
		FIt("Should return a 200 without the metrics of all the cluster", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			now := time.Now().Unix()
			VeryOldOwner(database, cluster.ID)

			//2. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
			q := req.URL.Query()
			q.Add("start", "1000")
			q.Add("end", fmt.Sprintf("%d", now))
			q.Add("period", "60")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			fmt.Println(res)
			//3. Analyse the result
			Expect(res.Code).To(Equal(200))
			json.Unmarshal(res.Body.Bytes(), &payload)
			data := payload["data"].([]interface{})
			fmt.Println(data)
			Expect(data).ToNot(BeEmpty())
		})
		It("Should return a 200 without the metrics of a namespace", func() {
		})
		It("Should return a 200 without the metrics of an owner", func() {
		})
		It("Should return a 200 with different metrics depending of the period", func() {
		})
	})
})
