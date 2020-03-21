package models

import (
	"encoding/json"
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
	APIToken   string `orm:"null;column(api_token)" json:"api_token"`
	NodesCount int    `orm:"null" json:"nodes_count"`
}

type Node struct {
	Model
	InstanceType      string                 `json:"instance_type"`
	Region            string                 `json:"region"`
	Zone              string                 `json:"zone"`
	Hostname          string                 `json:"hostname"`
	UID               string                 `orm:"column(uid)" json:"uid"`
	KubernetesVersion string                 `json:"kubernetes_version"`
	OS                string                 `orm:"column(os)" json:"os"`
	ClusterID         uint                   `orm:"column(cluster_id)" json:"cluster_id"`
	Labels            map[string]interface{} `orm:"-" json:"labels"`
	LabelsAsString    string                 `orm:"type(jsonb);column(labels);null" json:"-"`
}

type Namespace struct {
	Model
	Name           string                 `json:"name"`
	UID            string                 `orm:"column(uid)" json:"uid"`
	Resources      []Resources            `orm:"-" json:"metrics;omitempty"`
	Labels         map[string]interface{} `orm:"-" json:"labels"`
	LabelsAsString string                 `orm:"type(jsonb);column(labels);null" json:"-"`
	ClusterID      uint                   `orm:"column(cluster_id)" json:"cluster_id"`
}

type Metric struct {
	Model
	Name         string    `json:"metric_name"`
	NamespaceUID string    `orm:"column(namespace_uid)" json:"namespace_uid"`
	Value        float64   `json:"metric_value"`
	OwnerUID     string    `orm:"column(owner_uid)" json:"owner_uid"`
	ClusterID    uint      `orm:"column(cluster_id)" json:"cluster_id"`
	Period       time.Time `orm:"-"`
}

type Owner struct {
	Model
	Resources      []Resources            `orm:"-" json:"metrics;omitempty"`
	UID            string                 `orm:"column(uid)" json:"uid"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	NamespaceUID   string                 `orm:"column(namespace_uid)" json:"namespace_uid"`
	Labels         map[string]interface{} `orm:"-" json:"labels"`
	LabelsAsString string                 `orm:"type(jsonb);column(labels);null" json:"-"`
	ClusterID      uint                   `orm:"column(cluster_id)" json:"cluster_id"`
}

// PersistedModel interface used by Tables to calls callback methods on models
type PersistedModel interface {
	AfterLoad()
	BeforeSave()
}

// AfterLoad is called after loading object from DB
func (object *Model) AfterLoad() {

}

// BeforeSave called before the object is saved (updated or created) in DB
func (object *Model) BeforeSave() {

}

// AfterLoad is called after loading object from DB
func (cluster *Cluster) AfterLoad() {
}

// BeforeSave called before the object is saved (updated or created) in DB
func (cluster *Cluster) BeforeSave() {
}

// AfterLoad is called after loading object from DB
func (node *Node) AfterLoad() {
	var labels map[string]interface{}

	json.Unmarshal([]byte(node.LabelsAsString), &labels)
	node.Labels = labels
}

// BeforeSave called before the object is saved (updated or created) in DB
func (node *Node) BeforeSave() {
	labelsStr, _ := json.Marshal(node.Labels)
	node.LabelsAsString = string(labelsStr)
}

// AfterLoad is called after loading object from DB
func (ns *Namespace) AfterLoad() {
	var labels map[string]interface{}

	json.Unmarshal([]byte(ns.LabelsAsString), &labels)
	ns.Labels = labels
}

// BeforeSave called before the object is saved (updated or created) in DB
func (ns *Namespace) BeforeSave() {
	labelsStr, _ := json.Marshal(ns.Labels)
	ns.LabelsAsString = string(labelsStr)
}

// AfterLoad is called after loading object from DB
func (owner *Owner) AfterLoad() {
	var labels map[string]interface{}

	json.Unmarshal([]byte(owner.LabelsAsString), &labels)
	owner.Labels = labels
}

// BeforeSave called before the object is saved (updated or created) in DB
func (owner *Owner) BeforeSave() {
	labelsStr, _ := json.Marshal(owner.Labels)
	owner.LabelsAsString = string(labelsStr)
}

// GetUID Depending the metric reference returns the owner or namespace UID
func (metric *Metric) GetUID(refType string) string {
	switch refType {
	case "owner":
		return metric.OwnerUID
	case "namespace":
		return metric.NamespaceUID
	default:
		return ""
	}
}
