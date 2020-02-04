package db

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/models"
)

func TestTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Table Suite")
}

var _ = Describe("GormDatabase", func() {
	var (
		db    Database
		table Table
	)
	BeforeEach(func() {
		db, _ = New("sqlite3", ":memory:")
		db.Migrate()
		table = db.Clusters()
	})

	Describe("Find()", func() {
		It("returns record matching the query", func() {
			cluster1 := models.Cluster{Name: "cluster1"}
			db.Insert(&cluster1)
			cluster2 := models.Cluster{Name: "cluster2"}
			db.Insert(&cluster2)

			query := Query{
				Conditions: QueryConditions{"name": "cluster1"},
			}
			result := models.Cluster{}
			table.Find(query, &result)

			Expect(result.Name).To(Equal("cluster1"))
		})

		It("returns error if record not found", func() {
			result := models.Cluster{}
			query := Query{
				Conditions: QueryConditions{"name": "cluster2"},
			}
			err := table.Find(query, &result)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("FindAll()", func() {
		var (
			cluster1 models.Cluster
			cluster2 models.Cluster
			cluster3 models.Cluster
		)
		BeforeEach(func() {
			cluster1 = models.Cluster{Name: "cluster1"}
			db.Insert(&cluster1)
			cluster2 = models.Cluster{Name: "cluster2"}
			db.Insert(&cluster2)
			cluster3 = models.Cluster{Name: "cluster3"}
			db.Insert(&cluster3)
		})

		It("returns records matching the query", func() {
			query := Query{
				Conditions: QueryConditions{"name": "cluster1"},
			}
			var results []models.Cluster
			table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster1"))
		})

		It("returns all records if conditions", func() {
			query := Query{}
			var results []models.Cluster
			pagination, _ := table.FindAll(query, &results)

			Expect(len(results)).To(Equal(3))

			Expect(pagination.CurrentPage).To(Equal(0))
			Expect(pagination.MaxPage).To(Equal(0))
			Expect(pagination.TotalCount).To(Equal(3))
			Expect(pagination.Limit).To(Equal(0))
		})

		It("returns cluster3 because of orderby and limit", func() {
			query := Query{
				OrderBy: "name",
				Desc:    true,
				Limit:   1,
			}
			var results []models.Cluster
			table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster3"))
		})

		It("returns cluster2 because of offset", func() {
			query := Query{
				OrderBy: "name",
				Limit:   1,
				Page:    2,
			}
			var results []models.Cluster
			pagination, _ := table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster2"))

			Expect(pagination.CurrentPage).To(Equal(2))
			Expect(pagination.MaxPage).To(Equal(3))
			Expect(pagination.TotalCount).To(Equal(3))
			Expect(pagination.Limit).To(Equal(1))
		})

		It("returns error if invalid query", func() {
			query := Query{
				OrderBy: "fake",
			}
			var results []models.Cluster
			_, err := table.FindAll(query, &results)

			Expect(err).To(HaveOccurred())
		})
	})
})
