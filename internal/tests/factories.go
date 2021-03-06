package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/pilly-io/api/internal/db"
	"github.com/pilly-io/api/internal/models"
)

//MetricsFactory: create a metric given an owner/cluster
func MetricsFactory(database db.Database, clusterID uint, ownerUID string, namespaceUID string, name string, value float64, createdAt time.Time) *models.Metric {
	metric := &models.Metric{
		ClusterID:    clusterID,
		OwnerUID:     ownerUID,
		NamespaceUID: namespaceUID,
		Name:         name,
		Value:        value,
	}
	metric.Model.CreatedAt = createdAt
	database.Metrics().Insert(metric)
	return metric
}

//NamespaceFactory : create a namespace
func NamespaceFactory(database db.Database, clusterID uint, name string) (*models.Namespace, []*models.Metric) {
	var metrics []*models.Metric
	uid, _ := uuid.NewRandom()
	past := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	namespace := &models.Namespace{
		ClusterID: clusterID,
		Name:      name,
		UID:       uid.String(),
	}
	database.Namespaces().Insert(namespace)
	// TODO convert into a []{}
	//1. Generate CPU metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPUUsed, 1, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPUUsed, 1, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPUUsed, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPUUsed, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPUUsed, 1, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPURequested, 2, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPURequested, 2, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPURequested, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPURequested, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricCPURequested, 2, past.Add(time.Minute*4)))
	//2. Generate Memory metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryUsed, 1024, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryRequested, 2048, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, namespace.Name, namespace.UID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*4)))
	return namespace, metrics
}

//OwnerFactory : create an owner and generate its metrics
func OwnerFactory(database db.Database, clusterID uint, name string, namespaceUID string) (*models.Owner, []*models.Metric) {
	var metrics []*models.Metric
	uid, _ := uuid.NewRandom()
	past := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	owner := &models.Owner{
		ClusterID:    clusterID,
		Name:         name,
		Type:         "deployment",
		NamespaceUID: namespaceUID,
		UID:          uid.String(),
	}
	owner.Model.CreatedAt = past
	database.Owners().Insert(owner)
	// TODO convert into a []{}
	//1. Generate CPU metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPUUsed, 1, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPUUsed, 1, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPUUsed, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPUUsed, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPUUsed, 1, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPURequested, 2, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPURequested, 2, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPURequested, 2, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPURequested, 2, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricCPURequested, 2, past.Add(time.Minute*4)))
	//2. Generate Memory metrics
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryUsed, 1024, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryUsed, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryUsed, 1024, past.Add(time.Minute*4)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryRequested, 2048, past))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*1)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*2)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*3)))
	metrics = append(metrics, MetricsFactory(database, clusterID, owner.UID, namespaceUID, models.MetricMemoryRequested, 2048, past.Add(time.Minute*4)))
	return owner, metrics
}
