package db

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/helpers"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestOwnersTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OwnersTable Suite")
}

var _ = Describe("OwnersTable", func() {
	var (
		database                   db.Database
		cluster                    *models.Cluster
		tutumOwner, falcoOwner     *models.Owner
		tutumMetrics, falcoMetrics []*models.Metric
	)

	BeforeEach(func() {
		database = tests.SetupDB()
		cluster, _ = database.Clusters().Create("test", "aws")
		tests.NamespaceFactory(database, cluster.ID, "default")
		tests.NamespaceFactory(database, cluster.ID, "infrastructure")
		tutumOwner, tutumMetrics = tests.OwnerFactory(database, cluster.ID, "tutum", "default")
		falcoOwner, falcoMetrics = tests.OwnerFactory(database, cluster.ID, "falco", "infrastructure")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("ComputeResources()", func() {
		It("Should just cast the owners without filtering", func() {
			var owners []*models.Owner
			var metrics []*models.Metric
			owners = append(owners, tutumOwner, falcoOwner)
			metrics = append(tutumMetrics, falcoMetrics...)
			metricsIndexed := helpers.IndexMetrics(&metrics)
			database.Owners().ComputeResources(&owners, metricsIndexed)
			//2. Analyse the result
		})
	})
})
