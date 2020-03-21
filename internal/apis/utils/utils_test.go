package utils

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
	Describe("ConvertTimestampToTime()", func() {
		It("Should return a valid time", func() {
			t1, err := ConvertTimestampToTime("0")
			Expect(err).To(BeNil())
			Expect(t1.UTC().String()).To(Equal("1970-01-01 00:00:00 +0000 UTC"))
		})
		It("Should fails", func() {
			t1, err := ConvertTimestampToTime("error")
			Expect(err).ToNot(BeNil())
			Expect(t1).To(BeNil())
		})
	})
})
