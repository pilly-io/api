package frontend_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/apis/router"
	"github.com/pilly-io/api/internal/apis/utils"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestMetricsHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MetricsHandler Suite")
}

var _ = Describe("Metrics", func() {
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
		defaultNS, _ := tests.NamespaceFactory(database, cluster.ID, "default")
		infraNS, _ := tests.NamespaceFactory(database, cluster.ID, "infrastructure")
		tests.OwnerFactory(database, cluster.ID, "tutum", defaultNS.UID)
		tests.OwnerFactory(database, cluster.ID, "falco", infraNS.UID)
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("ListOwners() fails", func() {
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
			var payload utils.JsonFormat
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
			var payload utils.JsonFormat
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
			var payload utils.JsonFormat
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
			var payload utils.JsonFormat
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
			var payload utils.JsonFormat
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
			var payload utils.JsonFormat
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
})
