package models

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

const MetricCPUUsed = "CPU_USED"
const MetricCPURequested = "CPU_REQUESTED"
const MetricMemoryUsed = "MEMORY_USED"
const MetricMemoryRequested = "MEMORY_REQUESTED"

//IndexedMetrics : index the metrics by owner, period and type as it's easier to process
type IndexedMetrics = map[string]map[time.Time]map[string]Metric

type Resources struct {
	ResourcesTimestamp time.Time              `gorm:"-" json:"timestamp;omitempty"`
	ResourcesUsed      map[string]interface{} `gorm:"-" json:"resources_used;omitempty"`
	ResourcesRequested map[string]interface{} `gorm:"-" json:"resources_requested;omitempty"`
	Price              float64                `gorm:"-" json:"price;omitempty"`
	Score              float64                `gorm:"-" json:"score;omitempty"`
}

//Model : a copy of gorm.Model with json annotations
type Model struct {
	ID        uint       `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type Cluster struct {
	Model
	Resources
	Name       string `gorm:"unique;not null" json:"name"`
	Provider   string `json:"provider"`
	Region     string `json:"region"`
	APIToken   string `json:"api_token"`
	NodesCount int    `json:"nodes_count"`
}

type Node struct {
	Model
	InstanceType      string `json:"instance_type"`
	Region            string `json:"region"`
	Zone              string `json:"zone"`
	Hostname          string `json:"hostname"`
	UID               string `json:"uid"`
	KubernetesVersion string `json:"kubernetes_version"`
	OS                string `json:"os"`
	ClusterID         int    `json:"cluster_id"`
}

type Namespace struct {
	Model
	Name      string         `json:"name"`
	Labels    postgres.Jsonb `json:"labels"`
	ClusterID uint           `json:"cluster_id"`
}

type Metric struct {
	Model
	Name      string    `json:"metric_name"`
	Value     float64   `json:"metric_value"`
	OwnerUID  string    `json:"owner_uid"`
	ClusterID uint      `json:"cluster_id"`
	Period    time.Time `gorm:"-"`
}

type Owner struct {
	Model
	Metrics   []Resources    `gorm:"-" json:"metrics;omitempty"`
	UID       string         `json:"uid"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Namespace string         `json:"namespace"`
	Labels    postgres.Jsonb `json:"labels"`
	ClusterID uint           `json:"cluster_id"`
}
