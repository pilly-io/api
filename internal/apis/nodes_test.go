package apis

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/apis/middlewares"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestNodesHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NodesHandler Suite")
}

var _ = Describe("NodesHandler", func() {
	var (
		//handler *ClustersHandler
		engine   *gin.Engine
		database *db.GormDatabase
		cluster  *models.Cluster
	)
	BeforeEach(func() {
		database, _ = db.New("sqlite3", ":memory:")
		database.Migrate()
		gin.SetMode(gin.TestMode)
		engine = gin.New()
		SetupRouter(engine, database)
		cluster, _ = database.Clusters().Create("test", "aws")
	})

	Describe("Sync()", func() {
		It("should create node if does not exist", func() {
			data := tests.LoadFile("testdata/nodes.json")
			res := httptest.NewRecorder()

			count := 0
			countQuery := db.Query{
				Conditions: db.QueryConditions{
					"cluster_id": cluster.ID,
				},
			}
			database.Nodes().Count(countQuery)
			Expect(count).To(Equal(0))

			req, _ := http.NewRequest("POST", "/api/v1/collector/nodes", bytes.NewBuffer(data))
			req.Header.Set(middlewares.ClusterAuthHeaderName, cluster.APIToken)
			engine.ServeHTTP(res, req)

			database.Nodes().Count(countQuery)
			Expect(count).To(Equal(1))

			// //2. Analyse the result
			// json.Unmarshal(res.Body.Bytes(), &paylod)
			// Expect(res.Code).To(Equal(201))
			// Expect(paylod["data"]).To(HaveKeyWithValue("name", "cluster1"))
			// Expect(paylod["data"]).To(HaveKeyWithValue("provider", "aws"))
		})
	})
})
