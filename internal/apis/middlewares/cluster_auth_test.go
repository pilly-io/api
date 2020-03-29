package middlewares

import (
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

func TestClustersAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClustersAuth Suite")
}

var _ = Describe("ClustersAuth", func() {
	var (
		middleware      gin.HandlerFunc
		router          *gin.Engine
		called          bool
		cluster         *models.Cluster
		receivedCluster *models.Cluster
		database        db.Database
	)
	BeforeEach(func() {
		called = false
		database = tests.GetDB()
		cluster = &models.Cluster{
			Name:     "cluster1",
			APIToken: "1234567",
		}
		database.Clusters().Insert(cluster)

		middleware = CluserAuthMiddleware(database.Clusters())
		router = gin.New()
		router.Use(middleware)
		router.GET("/", func(c *gin.Context) {
			called = true
			contextValue, _ := c.Get("cluster")
			receivedCluster = contextValue.(*models.Cluster)
		})
	})

	It("should set cluster in context and call handler", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(ClusterAuthHeaderName, cluster.APIToken)
		router.ServeHTTP(w, req)

		Expect(called).To(BeTrue())
		Expect(receivedCluster.Name).To(Equal(cluster.Name))
	})

	It("should not call handler if cluster not found", func() {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set(ClusterAuthHeaderName, "XXXXXX")
		router.ServeHTTP(w, req)

		Expect(called).To(BeFalse())
		Expect(w.Result().StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
