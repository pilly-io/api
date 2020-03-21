package db_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/helpers"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestNamespacesTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NamespacesTable Suite")
}

var _ = Describe("NamespacesTable", func() {
	var (
		database                       db.Database
		cluster                        *models.Cluster
		tutumNamespace, falcoNamespace *models.Namespace
		tutumMetrics, falcoMetrics     []*models.Metric
	)

	BeforeEach(func() {
		database = tests.GetDB()
		cluster, _ = database.Clusters().Create("test", "aws")
		tutumNamespace, tutumMetrics = tests.NamespaceFactory(database, cluster.ID, "default")
		falcoNamespace, falcoMetrics = tests.NamespaceFactory(database, cluster.ID, "infrastructure")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("ComputeResources()", func() {
		It("Should compute the resources for 2 namespaces", func() {
			var namespaces []*models.Namespace
			var metrics []*models.Metric
			namespaces = append(namespaces, tutumNamespace, falcoNamespace)
			metrics = append(tutumMetrics, falcoMetrics...)
			metricsIndexed := helpers.IndexMetrics(&metrics)
			database.Namespaces().ComputeResources(&namespaces, metricsIndexed)
			// Check the Resources field of the first namespace
			Expect(namespaces[0].Resources).To(HaveLen(1))
			Expect(namespaces[0].Resources[0].ResourcesUsed).To(HaveKey("cpu"))
			Expect(namespaces[0].Resources[0].ResourcesUsed).To(HaveKey("memory"))
			Expect(namespaces[0].Resources[0].ResourcesRequested).To(HaveKey("cpu"))
			Expect(namespaces[0].Resources[0].ResourcesRequested).To(HaveKey("memory"))
			// Check the Resources field of the second namespaces
			Expect(namespaces[1].Resources).To(HaveLen(1))
			Expect(namespaces[1].Resources[0].ResourcesUsed).To(HaveKey("cpu"))
			Expect(namespaces[1].Resources[0].ResourcesUsed).To(HaveKey("memory"))
			Expect(namespaces[1].Resources[0].ResourcesRequested).To(HaveKey("cpu"))
			Expect(namespaces[1].Resources[0].ResourcesRequested).To(HaveKey("memory"))
		})
		It("Should compute the resources for only 1 namespace", func() {
			var namespaces []*models.Namespace
			namespaces = append(namespaces, tutumNamespace, falcoNamespace)
			metricsIndexed := helpers.IndexMetrics(&tutumMetrics)
			database.Namespaces().ComputeResources(&namespaces, metricsIndexed)
			// Check the Resources field of the first namespace
			Expect(namespaces[0].Resources[0].ResourcesUsed).To(HaveKey("cpu"))
			Expect(namespaces[0].Resources[0].ResourcesUsed).To(HaveKey("memory"))
			Expect(namespaces[0].Resources[0].ResourcesRequested).To(HaveKey("cpu"))
			Expect(namespaces[0].Resources[0].ResourcesRequested).To(HaveKey("memory"))
			// Check the Resources field of the second namespace
			Expect(namespaces[1].Resources).To(BeNil())
		})
	})
})
