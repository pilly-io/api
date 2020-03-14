package helpers

import (
	"time"

	"github.com/pilly-io/api/internal/models"
)

//IndexMetrics : this is ugly but life is ugly
// Given a list of metrics, convert it to a IndexedMetrics type
func IndexMetrics(metrics *[]*models.Metric) *models.IndexedMetrics {
	metricsIndexed := make(models.IndexedMetrics)
	for _, metric := range *metrics {
		if _, exist := metricsIndexed[metric.OwnerUID]; !exist {
			metricsIndexed[metric.OwnerUID] = map[time.Time]map[string]models.Metric{metric.Period: {metric.Name: *metric}}
		} else {
			if _, exist := metricsIndexed[metric.OwnerUID][metric.Period]; !exist {
				metricsIndexed[metric.OwnerUID][metric.Period] = map[string]models.Metric{metric.Name: *metric}
			} else {
				metricsIndexed[metric.OwnerUID][metric.Period][metric.Name] = *metric
			}
		}
	}
	return &metricsIndexed
}
