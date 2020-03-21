package helpers

import (
	"time"

	"github.com/pilly-io/api/internal/models"
)

//IndexMetrics : this is ugly but life is ugly
// Given a list of metrics, convert it to an IndexedMetrics type
func IndexMetrics(metrics *[]*models.Metric, refType string) *models.IndexedMetrics {
	metricsIndexed := make(models.IndexedMetrics)
	if metrics == nil {
		return &metricsIndexed
	}
	for _, metric := range *metrics {
		if _, exist := metricsIndexed[metric.GetUID(refType)]; !exist {
			metricsIndexed[metric.GetUID(refType)] = map[time.Time]map[string]models.Metric{metric.Period: {metric.Name: *metric}}
		} else {
			if _, exist := metricsIndexed[metric.GetUID(refType)][metric.Period]; !exist {
				metricsIndexed[metric.GetUID(refType)][metric.Period] = map[string]models.Metric{metric.Name: *metric}
			} else {
				metricsIndexed[metric.GetUID(refType)][metric.Period][metric.Name] = *metric
			}
		}
	}
	return &metricsIndexed
}
