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
)

func TestRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Clusters", func() {
	var (
		//handler *ClustersHandler
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

	Describe("Create()", func() {
		It("Should create a record and returns a 201", func() {
			var paylod jsonFormat
			res := httptest.NewRecorder()

			//1. Create the POST request
			clusterBytes := []byte(`{"name":"cluster1", "provider":"aws"}`)
			req, _ := http.NewRequest("POST", "/api/v1/clusters", bytes.NewBuffer(clusterBytes))
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			json.Unmarshal(res.Body.Bytes(), &paylod)
			Expect(res.Code).To(Equal(201))
			Expect(paylod["data"]).To(HaveKeyWithValue("name", "cluster1"))
			Expect(paylod["data"]).To(HaveKeyWithValue("provider", "aws"))
		})
		It("Should fails because the record already exist", func() {
			var paylod jsonFormat
			res := httptest.NewRecorder()
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Insert(&cluster1)

			//1. Create the POST request
			clusterBytes := []byte(`{"name":"cluster1", "provider":"aws"}`)
			req, _ := http.NewRequest("POST", "/api/v1/clusters", bytes.NewBuffer(clusterBytes))
			engine.ServeHTTP(res, req)

			//2. Analyse the result
			json.Unmarshal(res.Body.Bytes(), &paylod)
			Expect(res.Code).To(Equal(409))
			Expect(paylod["errors"]).To(HaveLen(1))
		})
	})
})
