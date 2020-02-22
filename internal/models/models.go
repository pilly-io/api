package models

import (
	"time"
)

const MetricCPUUsed = "CPU_USED"
const MetricCPURequested = "CPU_REQUESTED"
const MetricMemoryUsed = "MEMORY_USED"
const MetricMemoryRequested = "MEMORY_REQUESTED"

//IndexedMetrics : index the metrics by owner, period and type as it's easier to process
type IndexedMetrics = map[string]map[time.Time]map[string]Metric

type Resources struct {
	ResourcesTimestamp time.Time              `orm:"-" json:"timestamp;omitempty"`
	ResourcesUsed      map[string]interface{} `orm:"-" json:"resources_used;omitempty"`
	ResourcesRequested map[string]interface{} `orm:"-" json:"resources_requested;omitempty"`
	Price              float64                `orm:"-" json:"price;omitempty"`
	Score              float64                `orm:"-" json:"score;omitempty"`
}

//Model : a copy of orm.Model with json annotations
type Model struct {
	ID        uint       `orm:"pk;column(id);auto" json:"id"`
	CreatedAt time.Time  `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time  `orm:"auto_now_add;type(datetime)" json:"updated_at"`
	DeletedAt *time.Time `orm:"null" json:"deleted_at"`
}

type Cluster struct {
	Model
	Resources
	Name       string `orm:"unique" json:"name"`
	Provider   string `orm:"null" json:"provider"`
	Region     string `orm:"null" json:"region"`
	APIToken   string `orm:"null" json:"api_token"`
	NodesCount int    `orm:"null" json:"nodes_count"`
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
	ClusterID         uint   `json:"cluster_id"`
	//Labels            postgres.Jsonb `json:"labels"`
}

type Namespace struct {
	Model
	Name string `json:"name"`
	//Labels    postgres.Jsonb `json:"labels"`
	ClusterID uint `json:"cluster_id"`
}

type Metric struct {
	Model
	Name      string    `json:"metric_name"`
	Value     float64   `json:"metric_value"`
	OwnerUID  string    `json:"owner_uid"`
	ClusterID uint      `json:"cluster_id"`
	Period    time.Time `orm:"-"`
}

type Owner struct {
	Model
	Metrics   []Resources `orm:"-" json:"metrics;omitempty"`
	UID       string      `json:"uid"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Namespace string      `json:"namespace"`
	//Labels    postgres.Jsonb `json:"labels"`
	ClusterID uint `json:"cluster_id"`
}
