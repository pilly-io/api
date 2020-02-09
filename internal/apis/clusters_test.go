package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
)

func TestRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Suite")
}

var _ = Describe("Clusters", func() {
	var (
		//handler *ClustersHandler
		engine *gin.Engine
	)
	BeforeEach(func() {
		database, _ := db.New("sqlite3", ":memory:")
		database.Migrate()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		SetupRouter(engine, database)
	})

	Describe("Create()", func() {
		It("Should create a record and returns a 201", func() {
			res := httptest.NewRecorder()
			cluster := []byte(`{"name":"cluster1", "provider":"aws"}`)
			req, _ := http.NewRequest("POST", "/api/v1/clusters", bytes.NewBuffer(cluster))
			engine.ServeHTTP(res, req)
			Expect(res.Code).To(Equal(201))
			Expect(res.Body.String()).To(Equal(`{"name":"cluster1", "provider":"aws"}`))
		})
	})
})
