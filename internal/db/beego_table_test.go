package db_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
	"github.com/pilly-io/api/internal/tests"
)

func TestTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Table Suite")
}

var _ = Describe("BeegoDatabase", func() {
	var (
		database db.Database
		table    db.Table
	)
	BeforeEach(func() {
		database = tests.GetDB()
		table = database.Clusters()
	})

	AfterEach(func() {
		database.Flush()
	})

	Describe("Find()", func() {
		It("returns record matching the query", func() {
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Clusters().Insert(&cluster1)
			cluster2 := models.Cluster{Name: "cluster2"}
			database.Clusters().Insert(&cluster2)

			query := db.Query{
				Conditions: db.QueryConditions{"name": "cluster1"},
			}
			result := models.Cluster{}
			table.Find(query, &result)

			Expect(result.Name).To(Equal("cluster1"))
		})

		It("returns error if record not found", func() {
			result := models.Cluster{}
			query := db.Query{
				Conditions: db.QueryConditions{"name": "cluster2"},
			}
			err := table.Find(query, &result)
			fmt.Println(err)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Exists()", func() {
		It("returns true if record exists", func() {
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Clusters().Insert(&cluster1)

			query := db.Query{
				Conditions: db.QueryConditions{"name": cluster1.Name},
			}
			exists := table.Exists(query)
			Expect(exists).To(BeTrue())
		})

		It("returns false if record not found", func() {
			query := db.Query{
				Conditions: db.QueryConditions{"name": "cluster2"},
			}
			exists := table.Exists(query)
			Expect(exists).To(BeFalse())
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
			database.Clusters().Insert(&cluster1)
			cluster2 = models.Cluster{Name: "cluster2"}
			database.Clusters().Insert(&cluster2)
			cluster3 = models.Cluster{Name: "cluster3"}
			database.Clusters().Insert(&cluster3)
		})

		It("returns records matching the query", func() {
			query := db.Query{
				Conditions: db.QueryConditions{"name": "cluster1"},
			}
			var results []*models.Cluster
			table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster1"))
		})

		It("returns all records if conditions", func() {
			query := db.Query{}
			var results []*models.Cluster
			pagination, _ := table.FindAll(query, &results)

			Expect(len(results)).To(Equal(3))

			Expect(pagination.CurrentPage).To(Equal(0))
			Expect(pagination.MaxPage).To(Equal(0))
			Expect(pagination.TotalCount).To(Equal(3))
			Expect(pagination.Limit).To(Equal(0))
		})

		It("returns cluster3 because of orderby and limit", func() {
			query := db.Query{
				OrderBy: "-name",
				Limit:   1,
			}
			var results []*models.Cluster
			table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster3"))
		})

		It("returns cluster2 because of offset", func() {
			query := db.Query{
				OrderBy: "name",
				Limit:   1,
				Page:    2,
			}
			var results []*models.Cluster
			pagination, _ := table.FindAll(query, &results)

			Expect(len(results)).To(Equal(1))
			Expect(results[0].Name).To(Equal("cluster2"))

			Expect(pagination.CurrentPage).To(Equal(2))
			Expect(pagination.MaxPage).To(Equal(3))
			Expect(pagination.TotalCount).To(Equal(3))
			Expect(pagination.Limit).To(Equal(1))
		})
	})

	Describe("Count()", func() {
		It("returns number of records that match query", func() {
			cluster1 := models.Cluster{Name: "cluster1"}
			database.Clusters().Insert(&cluster1)

			query := db.Query{
				Conditions: db.QueryConditions{"name": cluster1.Name},
			}
			count := table.Count(query)
			Expect(count).To(Equal(1))
		})

		It("returns 0 if no records match", func() {
			query := db.Query{
				Conditions: db.QueryConditions{"name": "cluster2"},
			}
			count := table.Count(query)
			Expect(count).To(Equal(0))
		})
	})

	Describe("Update()", func() {
		It("update record in DB", func() {
			var clusterInDB models.Cluster

			cluster := models.Cluster{Name: "cluster1"}
			table.Insert(&cluster)

			cluster.Name = "New name"
			table.Update(&cluster)

			table.Find(db.Query{
				Conditions: db.QueryConditions{"id": cluster.ID},
			}, &clusterInDB)

			Expect(clusterInDB.Name).To(Equal("New name"))
		})
	})

	Describe("Insert()", func() {
		It("creates a record", func() {
			var count = 0
			cluster := models.Cluster{Name: "cluster1"}
			count = database.Clusters().Count(db.Query{})
			Expect(count).To(Equal(0))

			database.Clusters().Insert(&cluster)

			count = database.Clusters().Count(db.Query{})
			Expect(count).To(Equal(1))
		})
	})
})
