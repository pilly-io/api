package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

var _ = Describe("Owners", func() {
	var (
		engine   *gin.Engine
		database *db.GormDatabase
	)
	BeforeEach(func() {
		database, _ = db.New("sqlite3", ":memory:")
		database.Migrate()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		SetupRouter(engine, database)
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
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(400))
		})
		It("Should return a 400 as start is an invalid timestamp", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

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
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

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
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

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
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

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
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

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
		It("Should return a 200 without the metrics of all the cluster", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters/1/owners/metrics", nil)
			q := req.URL.Query()
			q.Add("start", "1000")
			q.Add("end", "10000000000")
			q.Add("period", "60")
			req.URL.RawQuery = q.Encode()
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			Expect(res.Code).To(Equal(200))
			json.Unmarshal(res.Body.Bytes(), &payload)
			data := payload["data"].([]interface{})
			fmt.Println(data)
			Expect(data).ToNot(BeNil())
		})
		It("Should return a 200 without the metrics of a namespace", func() {
		})
		It("Should return a 200 without the metrics of an owner", func() {
		})
		It("Should return a 200 with different metrics depending of the period", func() {
		})
	})
})
