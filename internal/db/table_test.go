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
			result, _ := table.Find(query)

			Expect(result.Name).To(Equal("cluster1"))
		})

		It("returns error if record not found", func() {
			result := models.Cluster{}
			query := Query{
				Conditions: QueryConditions{"name": "cluster2"},
			}
			result, err := table.Find(query)
			Expect(err).To(HaveOccurred())
			Expect(result).To(BeNil())
		})
	})
})
