package helpers

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestConvert(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Convert Suite")
}

var _ = Describe("Convert", func() {
	var (
		database db.Database
		cluster  *models.Cluster
	)

	BeforeEach(func() {
		database = tests.GetDB()
		cluster, _ = database.Clusters().Create("test", "aws")
		tests.NamespaceFactory(database, cluster.ID, "default")
		tests.NamespaceFactory(database, cluster.ID, "infrastructure")
		tests.OwnerFactory(database, cluster.ID, "tutum", "default")
		tests.OwnerFactory(database, cluster.ID, "falco", "infrastructure")
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("IndexMetrics()", func() {
		It("Map per OwnerUID, per period and per metric", func() {
			now := time.Now().UTC()
			past := now.Add(time.Minute * -10)
			var metrics []*models.Metric
			metric1 := tests.MetricsFactory(database, cluster.ID, "owner1", models.MetricCPUUsed, 10, now)
			metric1.Period = now
			metric2 := tests.MetricsFactory(database, cluster.ID, "owner1", models.MetricMemoryUsed, 100, now)
			metric2.Period = now
			metric3 := tests.MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 20, now)
			metric3.Period = now
			metric4 := tests.MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 200, past)
			metric4.Period = past
			metric5 := tests.MetricsFactory(database, cluster.ID, "owner2", models.MetricMemoryUsed, 300, past)
			metric5.Period = past
			metrics = append(metrics, metric1, metric2, metric3, metric4, metric5)
			indexed := IndexMetrics(&metrics)
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
