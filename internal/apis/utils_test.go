package apis

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pilly-io/api/internal/models"
)

var _ = Describe("Utils", func() {
	Describe("ObjectToJSON()", func() {
		It("Should return a map with a single object", func() {
			cluster1 := models.Cluster{Name: "cluster1"}
			payload := ObjectToJSON(cluster1)
			Expect(payload).To(HaveKey("data"))
			Expect(payload["data"]).To(Equal(cluster1))
		})
		It("Should return a map with a list of objects", func() {
			clusters := []models.Cluster{
				models.Cluster{Name: "cluster1"},
				models.Cluster{Name: "cluster2"},
			}
			payload := ObjectToJSON(clusters)
			Expect(payload).To(HaveKey("data"))
			Expect(payload["data"]).To(Equal(clusters))
		})
	})
})
