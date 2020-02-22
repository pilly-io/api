package apis

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestClustersHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClustersHandler Suite")
}

var _ = Describe("Clusters", func() {
	var (
		//handler *ClustersHandler
		engine   *gin.Engine
		database db.Database
	)
	BeforeEach(func() {
		database = tests.SetupDB()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		SetupRouter(engine, database)
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("Create()", func() {
		It("Should create a record and returns a 201", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()

			//1. Create the POST request
			clusterBytes := []byte(`{"name":"cluster1", "provider":"aws"}`)
			req, _ := http.NewRequest("POST", "/api/v1/clusters", bytes.NewBuffer(clusterBytes))
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			json.Unmarshal(res.Body.Bytes(), &payload)
			Expect(res.Code).To(Equal(201))
			Expect(payload["data"]).To(HaveKeyWithValue("name", "cluster1"))
			Expect(payload["data"]).To(HaveKeyWithValue("provider", "aws"))
		})
		It("Should fails because the record already exist", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Clusters().Insert(&cluster1)

			//1. Create the POST request
			clusterBytes := []byte(`{"name":"cluster1", "provider":"aws"}`)
			req, _ := http.NewRequest("POST", "/api/v1/clusters", bytes.NewBuffer(clusterBytes))
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			json.Unmarshal(res.Body.Bytes(), &payload)
			Expect(res.Code).To(Equal(409))
			Expect(payload["errors"]).To(HaveLen(1))
		})
	})
	Describe("List()", func() {
		It("Should get all the clusters and return 200", func() {
			var payload jsonFormat
			res := httptest.NewRecorder()
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Clusters().Insert(&cluster1)
			cluster2 := models.Cluster{Name: "cluster2"}
			database.Clusters().Insert(&cluster2)

			//1. Create the GET request
			req, _ := http.NewRequest("GET", "/api/v1/clusters", nil)
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			json.Unmarshal(res.Body.Bytes(), &payload)
			Expect(res.Code).To(Equal(200))
			Expect(payload["data"]).To(HaveLen(2))
			json := payload["data"].([]interface{})
			Expect(json[0]).To(HaveKeyWithValue("name", "cluster1"))
			Expect(json[1]).To(HaveKeyWithValue("name", "cluster2"))
		})
	})
})
