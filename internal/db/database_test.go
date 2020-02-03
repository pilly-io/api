package db

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/models"
)

func TestRunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Database Suite")
}

var _ = Describe("GormDatabase", func() {
	var (
		db *GormDatabase
	)
	BeforeEach(func() {
		db, _ = New("sqlite3", ":memory:")
		db.Migrate()
	})

	Describe("Insert()", func() {
		It("creates a record", func() {
			var count = 0
			cluster := models.Cluster{Name: "cluster1"}
			db.Model(&models.Cluster{}).Count(&count)
			Expect(count).To(Equal(0))

			db.Insert(&cluster)

			db.Model(&models.Cluster{}).Count(&count)
			Expect(count).To(Equal(1))
		})
	})
})
