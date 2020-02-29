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
	metric := &models.Metric{
		ClusterID: clusterID,
		OwnerUID:  ownerUID,
		Name:      name,
		Value:     value,
	}
	metric.Model.CreatedAt = createdAt
	database.Metrics().Insert(metric)
	return metric
}

//NamespaceFactory : create a namespace
func NamespaceFactory(database db.Database, clusterID uint, name string) {
	namespace := &models.Namespace{
		ClusterID: clusterID,
		Name:      name,
	}
	database.Namespaces().Insert(namespace)
}

//OwnerFactory : create an owner and generate its metrics
func OwnerFactory(database db.Database, clusterID uint, name string, namespace string) (*models.Owner, []*models.Metric) {
	var metrics []*models.Metric
	uid, _ := uuid.NewRandom()
	past := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	owner := &models.Owner{
		ClusterID: clusterID,
		Name:      name,
		Type:      "deployment",
		Namespace: namespace,
		UID:       uid.String(),
	}
	owner.Model.CreatedAt = past
	database.Owners().Insert(owner)
	// TODO convert into a []{}
	//1. Generate CPU metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPUUsed, 1, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPUUsed, 1, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPUUsed, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPUUsed, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPUUsed, 1, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPURequested, 2, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPURequested, 2, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPURequested, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPURequested, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricCPURequested, 2, past.Add(time.Minute*4)))
	//2. Generate Memory metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryUsed, 1024, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryRequested, 2048, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*4)))
	return owner, metrics
}

var _ = Describe("Metrics", func() {
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
		NamespaceFactory(database, cluster.ID, "default")
		NamespaceFactory(database, cluster.ID, "infrastructure")
		OwnerFactory(database, cluster.ID, "tutum", "default")
		OwnerFactory(database, cluster.ID, "falco", "infrastructure")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("ListMetrics() fails", func() {
		It("Should return a 404 as cluster does not exist", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/xxx/owners/metrics")
			req, _ := http.NewRequest("GET", url, nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(404))
		})
		It("Should return a 400 as start is not defined", func() {
			res := httptest.NewRecorder()

			//1. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
		})
		It("Should return a 400 as start is an invalid timestamp", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
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
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
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
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
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
		It("Should return a 400 as period is invalid", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
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
		It("Should return a 200 without the metrics of all the cluster", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			now := time.Now().Unix()

			//2. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
			q := req.URL.Query()
			q.Add("start", "1000")
			q.Add("end", fmt.Sprintf("%d", now))
			q.Add("period", "180")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//3. Analyse the result
			Expect(res.Code).To(Equal(200))
			json.Unmarshal(res.Body.Bytes(), &payload)
			data := payload["data"].([]interface{})
			Expect(data).To(HaveLen(2))
		})
		It("Should return a 200 without the metrics of a namespace", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			now := time.Now().Unix()

			//2. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
			q := req.URL.Query()
			q.Add("start", "1000")
			q.Add("end", fmt.Sprintf("%d", now))
			q.Add("period", "180")
			q.Add("namespace", "infrastructure")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//3. Analyse the result
			Expect(res.Code).To(Equal(200))
			json.Unmarshal(res.Body.Bytes(), &payload)
			data := payload["data"].([]interface{})
			Expect(data).To(HaveLen(1))
		})
		It("Should return a 200 without the metrics of an owner", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			now := time.Now().Unix()

			//2. Create the GET request
			url := fmt.Sprintf("/api/v1/clusters/%d/owners/metrics", cluster.ID)
			req, _ := http.NewRequest("GET", url, nil)
			q := req.URL.Query()
			q.Add("start", "1000")
			q.Add("end", fmt.Sprintf("%d", now))
			q.Add("period", "180")
			q.Add("namespace", "infrastructure")
			q.Add("owners", "deploy/falco")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//3. Analyse the result
			Expect(res.Code).To(Equal(200))
			json.Unmarshal(res.Body.Bytes(), &payload)
			data := payload["data"].([]interface{})
			Expect(data).To(HaveLen(1))
		})
		It("Should return a 200 with different metrics depending of the period", func() {
		})
	})
	Describe("indexMetrics()", func() {
		It("Map per OwnerUID, per period and per metric", func() {
			now := time.Now().UTC()
			past := now.Add(time.Minute * -10)
			var metrics []models.Metric
			metric1 := MetricsFactory(database, cluster.ID, "owner1", models.MetricCPUUsed, 10, now)
			metric1.Period = now
			metric2 := MetricsFactory(database, cluster.ID, "owner1", models.MetricMemoryUsed, 100, now)
			metric2.Period = now
			metric3 := MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 20, now)
			metric3.Period = now
			metric4 := MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 200, past)
			metric4.Period = past
			metric5 := MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 300, past)
			metric5.Period = past
			metrics = append(metrics, *metric1, *metric2, *metric3, *metric4, *metric5)
			indexed := indexMetrics(&metrics)
			// 1. First check the mapping per OwnerUID
			Expect(*indexed).To(HaveLen(2))
			Expect(*indexed).To(HaveKey("owner1"))
			Expect(*indexed).To(HaveKey("owner2"))
			// 2. Then check the mapping per Period
			owner1 := (*indexed)["owner1"]
			Expect(owner1).To(HaveLen(1))
			Expect(owner1).To(HaveKey(now))
			owner2 := (*indexed)["owner2"]
			Expect(owner2).To(HaveLen(2))
			Expect(owner2).To(HaveKey(now))
			Expect(owner2).To(HaveKey(past))
			// 3. Then check the mapping per Metric
			owner1Now := (*indexed)["owner1"][now]
			Expect(owner1Now).To(HaveLen(2))
			Expect(owner1Now).To(HaveKey(models.MetricCPUUsed))
			Expect(owner1Now).To(HaveKey(models.MetricMemoryUsed))
			owner2Now := (*indexed)["owner2"][now]
			Expect(owner2Now).To(HaveLen(1))
			Expect(owner2Now).To(HaveKey(models.MetricMemoryUsed))
			owner2Past := (*indexed)["owner2"][past]
			Expect(owner2Past).To(HaveLen(1))
			Expect(owner2Past).To(HaveKey(models.MetricMemoryUsed))
		})
	})
})
